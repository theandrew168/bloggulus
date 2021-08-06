.POSIX:
.SUFFIXES:

.PHONY: default
default: build

.PHONY: build
build:
	go build -o bloggulus cmd/web/main.go
#	go build -o bloggulus-worker cmd/worker/main.go
#	go build -o bloggulus-scheduler cmd/scheduler/main.go

.PHONY: dist
dist: build
	rm -fr dist/
	mkdir dist/
	cp bloggulus dist/
#	cp bloggulus-worker dist/
#	cp bloggulus-scheduler dist/
	cp -r migrations dist/
	cp -r static dist/
	cp -r templates dist/

.PHONY: run
run: build
	./bloggulus

.PHONY: test
test:
	go test -count=1 -v ./...

.PHONY: clean
clean:
	rm -fr bloggulus dist/

.PHONY: reload
reload: build
	./bloggulus -addblog http://lucumr.pocoo.org/feed.atom
	./bloggulus -addblog https://eli.thegreenplace.net/feeds/all.atom.xml
	./bloggulus -addblog https://eliasdaler.github.io/feed.xml
	./bloggulus -addblog https://nullprogram.com/feed/
	./bloggulus -addblog https://erikbern.com/atom.xml
	./bloggulus -addblog http://stevehanov.ca/blog/?atom
	./bloggulus -addblog https://stuartsierra.com/feed
	./bloggulus -addblog https://yogthos.net/feed.xml
	./bloggulus -addblog http://fabiensanglard.net/rss.xml
	./bloggulus -addblog https://ruslanspivak.com/feeds/all.atom.xml
	./bloggulus -addblog http://antirez.com/rss
	./bloggulus -addblog http://charlesleifer.com/blog/rss/
	./bloggulus -addblog https://dave.cheney.net/feed
	./bloggulus -addblog https://www.alexedwards.net/static/feed.rss
	./bloggulus -addblog https://www.joshmcguigan.com/rss.xml
	./bloggulus -addblog https://rachelbythebay.com/w/atom.xml
	./bloggulus -addblog https://floooh.github.io/feed.xml
	./bloggulus -addblog https://blog.notryan.com/rss.xml
	./bloggulus -addblog https://hookrace.net/blog/feed/
	./bloggulus -addblog http://peter.michaux.ca/feed/atom.xml
	./bloggulus -addblog http://jakob.space/feed.xml
	./bloggulus -addblog https://blog.regehr.org/feed
	./bloggulus -addblog https://maryrosecook.com/blog/feed.xml
	./bloggulus -addblog https://www.more-magic.net/feed.atom
	./bloggulus -addblog https://jacobian.org/posts/index.xml
	./bloggulus -addblog https://caseymuratori.com/blog_atom.rss
	./bloggulus -addblog https://shallowbrooksoftware.com/posts/index.xml
	./bloggulus -addblog https://www.bruceeckel.com/index.xml

.PHONY: sync
sync: build
	./bloggulus -syncblogs
