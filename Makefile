.POSIX:
.SUFFIXES:

.PHONY: default
default: build

.PHONY: build
build:
	go build -o bloggulus main.go

.PHONY: dist
dist: build
	rm -fr dist/
	mkdir dist/
	cp bloggulus dist/
	cp -r migrations dist/
	cp -r static dist/
	cp -r templates dist/

.PHONY: test
test:
	go test -count=1 -v ./...

.PHONY: clean
clean:
	rm -fr bloggulus dist/
