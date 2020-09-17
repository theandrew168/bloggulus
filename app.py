from datetime import datetime
from functools import wraps
from html.parser import HTMLParser
from io import StringIO
import os
import sys
import time

import bleach
import feedparser
from flask import Flask, render_template
from peewee import Model, SqliteDatabase
from peewee import CharField, DateTimeField, ForeignKeyField, TextField
from playhouse.sqlite_ext import FTSModel, SqliteExtDatabase
from playhouse.sqlite_ext import SearchField
import pytz


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
    'synchronous': 0,  # let OS handle file syncing (!!! use with caution !!!)
}
database = SqliteExtDatabase(DATABASE, autoconnect=False, pragmas=pragmas)


@app.before_request
def before_request():
    database.connect()

@app.teardown_request
def teardown_request(exc):
    if not database.is_closed():
        database.close()


class BaseModel(Model):
    class Meta:
        database = database

class Feed(BaseModel):
    url = CharField(unique=True)
    title = CharField()
    updated = DateTimeField()

    def __str__(self):
        return self.title

class Post(BaseModel):
    feed = ForeignKeyField(Feed, backref='posts')
    url = CharField(unique=True)
    title = CharField()
    updated = DateTimeField()
    content = TextField()

    def __str__(self):
        return self.title

class PostContent(FTSModel):
    content = SearchField()

    class Meta:
        database = database
        options = {'content': Post.content}


# https://stackoverflow.com/questions/753052/strip-html-from-strings-in-python
class MLStripper(HTMLParser):
    def __init__(self):
        super().__init__()
        self.reset()
        self.strict = False
        self.convert_charrefs = True
        self.text = StringIO()
    def handle_data(self, d):
        self.text.write(d)
    def get_data(self):
        return self.text.getvalue()

def strip_tags(html):
    s = MLStripper()
    s.feed(html)
    return s.get_data()

def add_feed(url):
    with database:
        d = feedparser.parse(url)
        feed = d['feed']

        title = feed['title']
        updated = feed['updated_parsed']
        updated = datetime.fromtimestamp(time.mktime(updated))
        updated = pytz.utc.localize(updated)

        print(url, title, updated)
        f = Feed.get_or_none(url=url)
        if f is None:
            f = Feed.create(url=url, title=title, updated=updated)
            print('  created')
        else:
            f.update(title=title, updated=updated).execute()
            print('  exists')

def sync_feeds():
    with database:
        for feed in Feed.select():
            d = feedparser.parse(feed.url)
            posts = d['entries']
            for post in posts:
                url = post['link']
                title = post['title']
                updated = post['updated_parsed']
                updated = datetime.fromtimestamp(time.mktime(updated))
                updated = pytz.utc.localize(updated)
                content = post['content']
                content = content[0]['value']
                content = strip_tags(content)

                print(url, title, updated)
                p = Post.get_or_none(feed=feed, url=url)
                if p is None:
                    p = Post.create(feed=feed, url=url, title=title, updated=updated, content=content)
                    print('  created')
                else:
                    p.update(title=title, updated=updated, content=content)
                    print('  exists')

        PostContent.rebuild()
        PostContent.optimize()

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
            .join(PostContent, on=(Post.id == PostContent.docid))
            .where(PostContent.match(search))
            .order_by(PostContent.rank()))  # could use bm25, bm25f, or lucene


@app.route('/')
def index():
    posts = Post.select().order_by(Post.updated.desc())[:20]
    return render_template('index.html', posts=posts)


# ensure the database and its tables exist
with database:
    database.create_tables([Feed, Post, PostContent])

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
        url = sys.argv[2]
        add_feed(url)
    elif sys.argv[1] == 'syncfeeds':
        sync_feeds()
    else:
        raise SystemExit(usage)

if __name__ == '__main__':
    main()
