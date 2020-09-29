from datetime import datetime, timedelta
from functools import wraps
from html.parser import HTMLParser
from io import StringIO
import os
import sys
import time
from urllib.parse import quote_plus
from urllib.request import urlopen

import bleach
import feedparser
from flask import Flask, render_template, request
from peewee import Model
from peewee import CharField, DateTimeField, ForeignKeyField, TextField
from playhouse.sqlite_ext import FTSModel, SqliteExtDatabase
from playhouse.sqlite_ext import SearchField


DATABASE = os.getenv('BLOGGULUS_DATABASE') or 'bloggulus.sqlite3'
SECRET_KEY = os.getenv('BLOGGULUS_SECRET_KEY') or 'bloggulus_development_secret_key'

app = Flask(__name__, root_path='.')
app.config.from_object(__name__)

# https://www.sqlite.org/pragma.html
pragmas = {
    'journal_mode': 'wal',  # write-ahead log mode
    'cache_size': -64 * 1024,  # 64MB cache
    'foreign_keys': 1,  # enforce foreign-key constraints
    'ignore_check_constraints': 0,  # enforce CHECK constraints
#    'synchronous': 0,  # let OS handle file syncing (!!! use with caution !!!)
}
database = SqliteExtDatabase(DATABASE, autoconnect=False, pragmas=pragmas)


@app.before_request
def before_request():
    database.connect()

@app.teardown_request
def teardown_request(exc):
    if not database.is_closed():
        database.close()

@app.template_filter('pretty_date')
def pretty_date(date):
    return date.strftime('%B %d, %Y')


class BaseModel(Model):
    class Meta:
        database = database

class Feed(BaseModel):
    url = CharField(unique=True)
    site_url = CharField()
    title = CharField()

    def __str__(self):
        return self.title

class Post(BaseModel):
    feed = ForeignKeyField(Feed, backref='posts')
    url = CharField(unique=True)
    title = CharField()
    updated = DateTimeField()

    def __str__(self):
        return self.title

class PostIndex(FTSModel):
    feed = SearchField()
    title = SearchField()
    content = SearchField()

    class Meta:
        database = database


def add_feed(url, site_url):
    with database:
        d = feedparser.parse(url)
        feed = d['feed']
        title = feed['title']

        print(url, title)

        # exit early if feed already exists
        f = Feed.get_or_none(url=url)
        if f is not None:
            print('  exists')
            return

        Feed.create(url=url, site_url=site_url, title=title)
        print('  create')

def sync_feeds():
    with database:
        feeds = list(Feed.select())

    for feed in feeds:
        with database:
            d = feedparser.parse(feed.url)
            posts = d['entries']
            for post in posts:
                url = post['link']
                title = post['title']
                print(feed.title, '-', title)

                updated = post.get('updated_parsed')
                if updated is None:
                    # if no updated date, just set to 30 days ago
                    updated = datetime.utcnow() - timedelta(days=30)
                else:
                    updated = datetime.fromtimestamp(time.mktime(updated))

                # continue early if post already exists
                p = Post.get_or_none(feed=feed, url=url)
                if p is not None:
                    print('  exists')
                    continue

                # grab contend from feed, fallback to manual fetch
                content = post.get('content')
                if content is None:
                    print(' no content, fetching manually...')
                    try:
                        with urlopen(url) as resp:
                            content = resp.read().decode()
                    except:
                        print('  manual fetch failed, too... RIP')
                        continue
                else:
                    content = content[0]['value']

                # strip any HTML from the content
                content = bleach.clean(content, strip=True, attributes={}, styles=[], tags=[])

                p = Post.create(feed=feed, url=url, title=title, updated=updated)
                PostIndex.create(docid=p.id, feed=feed.title, title=p.title, content=content)
                print('  create')

# http://charlesleifer.com/blog/saturday-morning-hacks-adding-full-text-search-to-the-flask-note-taking-app/
def search_posts(text):
    # make sure the search is clean and reasonably formatted
    words = [word.strip() for word in text.split() if word]
    if not words:
        return Post.select().where(Post.id == 0)
    else:
        search = ' '.join(words)

    return (Post
            .select()
            .join(PostIndex, on=(Post.id == PostIndex.docid))
            .where(PostIndex.match(search))
            .order_by(PostIndex.bm25()))  # could use rank, bm25, bm25f, or lucene


@app.route('/')
def index():
    search_text = request.args.get('q')
    search = search_text or ''

    search_param = quote_plus(search)

    page = request.args.get('p', 1)
    try:
        page = int(page)
    except:
        pass

    if search_text:
        posts = search_posts(search_text)
    else:
        posts = Post.select().order_by(Post.updated.desc())

    pages = posts.count() // 20 + 1
    posts = posts.paginate(page, 20)

    return render_template('index.html',
        posts=posts,
        search=search,
        search_param=search_param,
        page=page,
        pages=pages)

@app.route('/about')
def about():
    feeds = Feed.select().order_by(Feed.title)
    feeds = sorted(feeds, key=lambda feed: feed.title.lower())
    return render_template('about.html', feeds=feeds)

@app.route('/docs')
def docs():
    return render_template('docs.html')


# ensure the database and its tables exist
with database:
    database.create_tables([Feed, Post, PostIndex])

def main():
    # CLI usage and help
    usage = 'usage: {} {{gunicorn,addfeed,syncfeeds}} [extra_args]'.format(sys.argv[0])
    if '-h' in sys.argv or '--help' in sys.argv:
        raise SystemExit(usage)

    # ensure a valid quantity of args
    if len(sys.argv) < 2:
        raise SystemExit(usage)

    # choose an action based on the given command
    if sys.argv[1] == 'gunicorn':
        from gunicorn.app import wsgiapp
        sys.argv[1] = 'app:app'  # swap 'gunicorn' in argv for the WSGI app name
        wsgiapp.run()
    elif sys.argv[1] == 'addfeed':
        if len(sys.argv) < 3:
            raise SystemExit('addfeed: <url> <site_url>')
        url = sys.argv[2]
        site_url = sys.argv[3]
        add_feed(url, site_url)
    elif sys.argv[1] == 'syncfeeds':
        sync_feeds()
    else:
        raise SystemExit(usage)

if __name__ == '__main__':
    main()
