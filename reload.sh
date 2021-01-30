#!/bin/bash -e

go build -o bloggulus main.go
./bloggulus -addfeed http://lucumr.pocoo.org/feed.atom https://lucumr.pocoo.org
./bloggulus -addfeed https://eli.thegreenplace.net/feeds/all.atom.xml https://eli.thegreenplace.net
./bloggulus -addfeed https://eliasdaler.github.io/feed.xml https://eliasdaler.github.io
./bloggulus -addfeed https://nullprogram.com/feed/ https://nullprogram.com
./bloggulus -addfeed https://erikbern.com/atom.xml https://erikbern.com
./bloggulus -addfeed http://stevehanov.ca/blog/?atom http://stevehanov.ca
./bloggulus -addfeed https://stuartsierra.com/feed https://stuartsierra.com
./bloggulus -addfeed http://headerphile.com/feed/ http://headerphile.com
./bloggulus -addfeed https://yogthos.net/feed.xml https://yogthos.net
./bloggulus -addfeed http://fabiensanglard.net/rss.xml http://fabiensanglard.net
./bloggulus -addfeed https://ruslanspivak.com/feeds/all.atom.xml https://ruslanspivak.com
./bloggulus -addfeed http://antirez.com/rss http://antirez.com
./bloggulus -addfeed http://charlesleifer.com/blog/rss/ http://charlesleifer.com
./bloggulus -addfeed https://dave.cheney.net/feed https://dave.cheney.net
./bloggulus -addfeed https://www.alexedwards.net/static/feed.rss https://www.alexedwards.net
./bloggulus -addfeed https://www.joshmcguigan.com/rss.xml https://www.joshmcguigan.com
./bloggulus -addfeed https://rachelbythebay.com/w/atom.xml https://rachelbythebay.com
./bloggulus -addfeed https://floooh.github.io/feed.xml https://floooh.github.io
./bloggulus -addfeed https://blog.notryan.com/rss.xml https://blog.notryan.com
./bloggulus -addfeed https://hookrace.net/blog/feed/ https://hookrace.net
./bloggulus -addfeed http://peter.michaux.ca/feed/atom.xml http://peter.michaux.ca
./bloggulus -addfeed http://jakob.space/feed.xml http://jakob.space
./bloggulus -addfeed https://blog.regehr.org/feed https://blog.regehr.org
./bloggulus -addfeed https://maryrosecook.com/blog/feed.xml https://maryrosecook.com
./bloggulus -addfeed https://www.more-magic.net/feed.atom https://www.more-magic.net
./bloggulus -addfeed https://jacobian.org/posts/index.xml https://jacobian.org
./bloggulus -addfeed http://www.aaronsw.com/2002/feeds/pgessays.rss http://www.paulgraham.com
./bloggulus -addfeed https://caseymuratori.com/blog_atom.rss https://caseymuratori.com/
./bloggulus -addfeed https://shallowbrooksoftware.com/posts/index.xml https://shallowbrooksoftware.com
./bloggulus -addfeed https://www.bruceeckel.com/index.xml https://www.bruceeckel.com
