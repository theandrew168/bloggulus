from collections import namedtuple

import app


Feed = namedtuple('Feed', 'url site_url')

feeds = [
    Feed('http://lucumr.pocoo.org/feed.atom', 'https://lucumr.pocoo.org'),
    Feed('https://eli.thegreenplace.net/feeds/all.atom.xml', 'https://eli.thegreenplace.net'),
    Feed('https://eliasdaler.github.io/feed.xml', 'https://eliasdaler.github.io'),
    Feed('https://nullprogram.com/feed/', 'https://nullprogram.com'),
    Feed('https://erikbern.com/atom.xml', 'https://erikbern.com'),
    Feed('http://stevehanov.ca/blog/?atom', 'http://stevehanov.ca'),
    Feed('https://stuartsierra.com/feed', 'https://stuartsierra.com'),
    Feed('http://headerphile.com/feed/', 'http://headerphile.com'),
    Feed('https://yogthos.net/feed.xml', 'https://yogthos.net'),
    Feed('http://fabiensanglard.net/rss.xml', 'http://fabiensanglard.net'),
    Feed('https://ruslanspivak.com/feeds/all.atom.xml', 'https://ruslanspivak.com'),
    Feed('http://antirez.com/rss', 'http://antirez.com'),
    Feed('http://charlesleifer.com/blog/rss/', 'http://charlesleifer.com'),
    Feed('https://dave.cheney.net/feed', 'https://dave.cheney.net'),
    Feed('https://www.alexedwards.net/static/feed.rss', 'https://www.alexedwards.net'),
    Feed('https://www.joshmcguigan.com/rss.xml', 'https://www.joshmcguigan.com'),
    Feed('https://rachelbythebay.com/w/atom.xml', 'https://rachelbythebay.com'),
    Feed('https://floooh.github.io/feed.xml', 'https://floooh.github.io'),
    Feed('https://blog.notryan.com/rss.xml', 'https://blog.notryan.com'),
    Feed('https://hookrace.net/blog/feed/', 'https://hookrace.net'),
    Feed('http://peter.michaux.ca/feed/atom.xml', 'http://peter.michaux.ca'),
    Feed('http://jakob.space/feed.xml', 'http://jakob.space'),
    Feed('https://blog.regehr.org/feed', 'https://blog.regehr.org'),
    Feed('https://maryrosecook.com/blog/feed.xml', 'https://maryrosecook.com'),
    Feed('https://www.more-magic.net/feed.atom', 'https://www.more-magic.net'),
    Feed('https://jacobian.org/posts/index.xml', 'https://jacobian.org'),
    Feed('http://www.aaronsw.com/2002/feeds/pgessays.rss', 'http://www.paulgraham.com'),
    Feed('https://caseymuratori.com/blog_atom.rss', 'https://caseymuratori.com/'),
    Feed('https://shallowbrooksoftware.com/posts/index.xml', 'https://shallowbrooksoftware.com'),
]

for feed in feeds:
    print(feed)
    app.add_feed(feed.url, feed.site_url)

app.sync_feeds()
