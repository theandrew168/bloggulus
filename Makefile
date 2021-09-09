.POSIX:
.SUFFIXES:

.PHONY: default
default: build

.PHONY: build
build:
	go build -o bloggulus main.go

.PHONY: test
test:
	go test -count=1 -v ./...

.PHONY: format
format:
	go fmt ./...

.PHONY: clean
clean:
	rm -fr bloggulus
