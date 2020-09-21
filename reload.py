import app

feeds = [
#    'https://dmerej.info/blog/index.xml',  # lots of manual fetches
    'http://lucumr.pocoo.org/feed.atom',
    'https://eli.thegreenplace.net/feeds/all.atom.xml',
    'https://eliasdaler.github.io/feed.xml',
    'https://nullprogram.com/feed/',
    'https://erikbern.com/atom.xml',
#    'http://apoorvaj.io/feed.xml',  # 404
#    'https://tratt.net/laurie/blog/entries.rss',  # noisy
    'http://stevehanov.ca/blog/?atom',
    'https://stuartsierra.com/feed',
    'http://headerphile.com/feed/',
    'https://yogthos.net/feed.xml',
    'http://fabiensanglard.net/rss.xml',
    'https://ruslanspivak.com/feeds/all.atom.xml',
    'http://antirez.com/rss',
#    'https://eev.ee/feeds/blog.atom.xml',  # just kinda noisy
#    'https://eklitzke.org/index.rss',  # 404
    'http://charlesleifer.com/blog/rss/',
    'https://dave.cheney.net/feed',
    'https://www.alexedwards.net/static/feed.rss',
    'https://www.joshmcguigan.com/rss.xml',
    'https://rachelbythebay.com/w/atom.xml',
    'https://floooh.github.io/feed.xml',
    'https://blog.notryan.com/rss.xml',
#    'https://blog.veitheller.de/feed.rss',  # noisy
    'https://hookrace.net/blog/feed/',
    'http://peter.michaux.ca/feed/atom.xml',
    'http://jakob.space/feed.xml',
#    'https://fasterthanli.me/index.xml',  # borked
    'https://blog.regehr.org/feed',
#    'http://www.wilfred.me.uk/rss.xml',  # manual fetches AND no dates :(
    'https://maryrosecook.com/blog/feed.xml',
    'https://www.more-magic.net/feed.atom',
    'https://jacobian.org/posts/index.xml',
    'https://euandre.org/feed.blog.en.atom',
]

for feed in feeds:
    print('FEED', feed)
    app.add_feed(feed)

app.sync_feeds()
