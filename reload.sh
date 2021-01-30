#!/bin/bash -e

go build -o bloggulus main.go
./bloggulus -addblog http://lucumr.pocoo.org/feed.atom https://lucumr.pocoo.org
./bloggulus -addblog https://eli.thegreenplace.net/feeds/all.atom.xml https://eli.thegreenplace.net
./bloggulus -addblog https://eliasdaler.github.io/feed.xml https://eliasdaler.github.io
./bloggulus -addblog https://nullprogram.com/feed/ https://nullprogram.com
./bloggulus -addblog https://erikbern.com/atom.xml https://erikbern.com
./bloggulus -addblog http://stevehanov.ca/blog/?atom http://stevehanov.ca
./bloggulus -addblog https://stuartsierra.com/feed https://stuartsierra.com
./bloggulus -addblog http://headerphile.com/feed/ http://headerphile.com
./bloggulus -addblog https://yogthos.net/feed.xml https://yogthos.net
./bloggulus -addblog http://fabiensanglard.net/rss.xml http://fabiensanglard.net
./bloggulus -addblog https://ruslanspivak.com/feeds/all.atom.xml https://ruslanspivak.com
./bloggulus -addblog http://antirez.com/rss http://antirez.com
./bloggulus -addblog http://charlesleifer.com/blog/rss/ http://charlesleifer.com
./bloggulus -addblog https://dave.cheney.net/feed https://dave.cheney.net
./bloggulus -addblog https://www.alexedwards.net/static/feed.rss https://www.alexedwards.net
./bloggulus -addblog https://www.joshmcguigan.com/rss.xml https://www.joshmcguigan.com
./bloggulus -addblog https://rachelbythebay.com/w/atom.xml https://rachelbythebay.com
./bloggulus -addblog https://floooh.github.io/feed.xml https://floooh.github.io
./bloggulus -addblog https://blog.notryan.com/rss.xml https://blog.notryan.com
./bloggulus -addblog https://hookrace.net/blog/feed/ https://hookrace.net
./bloggulus -addblog http://peter.michaux.ca/feed/atom.xml http://peter.michaux.ca
./bloggulus -addblog http://jakob.space/feed.xml http://jakob.space
./bloggulus -addblog https://blog.regehr.org/feed https://blog.regehr.org
./bloggulus -addblog https://maryrosecook.com/blog/feed.xml https://maryrosecook.com
./bloggulus -addblog https://www.more-magic.net/feed.atom https://www.more-magic.net
./bloggulus -addblog https://jacobian.org/posts/index.xml https://jacobian.org
./bloggulus -addblog http://www.aaronsw.com/2002/feeds/pgessays.rss http://www.paulgraham.com
./bloggulus -addblog https://caseymuratori.com/blog_atom.rss https://caseymuratori.com/
./bloggulus -addblog https://shallowbrooksoftware.com/posts/index.xml https://shallowbrooksoftware.com
./bloggulus -addblog https://www.bruceeckel.com/index.xml https://www.bruceeckel.com
