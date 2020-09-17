from datetime import datetime
from functools import wraps
import os
import sys
import time

import feedparser
from flask import Flask, render_template
from peewee import Model, SqliteDatabase
from peewee import CharField, DateTimeField, ForeignKeyField
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
database = SqliteDatabase(DATABASE, autoconnect=False, pragmas=pragmas)


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

class Post(BaseModel):
    feed = ForeignKeyField(Feed, backref='posts')
    url = CharField(unique=True)
    title = CharField()
    updated = DateTimeField()


def add_feed(url):
    with database:
        d = feedparser.parse(url)
        feed = d['feed']

        title = feed['title']
        updated = feed['updated_parsed']
        updated = datetime.fromtimestamp(time.mktime(updated))
        updated = pytz.utc.localize(updated)

        print(url, title, updated)
        Feed.get_or_create(url=url, defaults={'title': title, 'updated': updated})

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

                print(url, title, updated)
                Post.get_or_create(feed=feed, url=url, defaults={'title': title, 'updated': updated})


@app.route('/')
def index():
    posts = Post.select().order_by(Post.updated.desc())
    return render_template('index.html', posts=posts)


# ensure the database and its tables exist
with database:
    database.create_tables([Feed, Post])

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
