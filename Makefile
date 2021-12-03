.POSIX:
.SUFFIXES:

.PHONY: default
default: build

.PHONY: build
build:
	go build -o bloggulus main.go

.PHONY: run
run:
	go run main.go -conf internal/test/bloggulus.conf

.PHONY: test
test:
	go run main.go -conf internal/test/bloggulus.conf -migrate
	go test -count=1 -v ./...

.PHONY: race
race:
	go run main.go -conf internal/test/bloggulus.conf -migrate
	go test -race -count=1 ./...

.PHONY: cover
cover:
	go run main.go -conf internal/test/bloggulus.conf -migrate
	go test -coverprofile=c.out -coverpkg=./... -count=1 ./...
	go tool cover -html=c.out

.PHONY: format
format:
	go fmt ./...

.PHONY: clean
clean:
	rm -fr bloggulus c.out dist/
