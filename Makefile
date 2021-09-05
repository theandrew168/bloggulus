.POSIX:
.SUFFIXES:

.PHONY: default
default: build

.PHONY: build
build:
	go build -o bloggulus-migrate cmd/migrate/main.go
	go build -o bloggulus-web cmd/web/main.go
	go build -o bloggulus-worker cmd/worker/main.go
	go build -o bloggulus-scheduler cmd/scheduler/main.go

.PHONY: dist
dist: build
	rm -fr dist/
	mkdir dist/
	cp bloggulus-* dist/
	cp -r static dist/
	cp -r templates dist/

.PHONY: test
test:
	go test -count=1 -v ./...

.PHONY: format
format:
	go fmt ./...

.PHONY: clean
clean:
	rm -fr bloggulus-* dist/
