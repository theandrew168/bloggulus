.POSIX:
.SUFFIXES:

.PHONY: default
default: build

.PHONY: build
build:
	lein with-profiles web:worker:scheduler uberjar
	cp target/uberjar/*.jar .

.PHONY: dist
dist: build
	rm -fr dist/
	mkdir dist/
	cp *.jar dist/
	cp -r migrations dist/
	cp -r static dist/
	cp -r templates dist/

.PHONY: test
test:
	lein test

.PHONY: clean
clean:
	rm -fr *.jar dist/ target/
