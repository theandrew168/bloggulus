.POSIX:
.SUFFIXES:

.PHONY: default
default: build

.PHONY: build
build:
	go build -o bloggulus main.go

# use wgo to watch for code changes and subsequently rebuild the app
.PHONY: run
run:
	DEBUG=1 go run github.com/bokwoon95/wgo@latest run -file .html -file .css main.go

# run the app using the local-only config file
.PHONY: run-local
run-local:
	DEBUG=1 go run github.com/bokwoon95/wgo@latest run -file .html -file .css main.go -conf bloggulus.local.conf

.PHONY: migrate
migrate:
	go run main.go -conf bloggulus.conf -migrate
	go run main.go -conf bloggulus.test.conf -migrate

.PHONY: test
test: migrate
	go test -count=1 -shuffle=on -race -vet=all -failfast ./...

.PHONY: cover
cover:
	go test -coverprofile=c.out -coverpkg=./... -count=1 ./...
	go tool cover -html=c.out

.PHONY: release
release:
	goreleaser release --clean --snapshot

.PHONY: deploy
deploy: release
	scp dist/bloggulus_linux_amd64_v1/bloggulus derz@bloggulus.com:/tmp/bloggulus
	ssh -t derz@bloggulus.com sudo install /tmp/bloggulus /usr/local/bin/bloggulus
	ssh -t derz@bloggulus.com sudo systemctl restart bloggulus

format:
	gofmt -l -s -w .

.PHONY: update
update: update-deps update-htmx update-alpine

.PHONY: update-deps
update-deps:
	go get -u ./...
	go mod tidy

# https://htmx.org/docs/#installing
.PHONY: update-htmx
update-htmx:
	curl -L -s -o public/js/htmx.min.js https://unpkg.com/htmx.org@2.x.x/dist/htmx.min.js

# https://alpinejs.dev/essentials/installation
.PHONY: update-alpine
update-alpine:
	curl -L -s -o public/js/alpine.min.js https://cdn.jsdelivr.net/npm/alpinejs@3.x.x/dist/cdn.min.js

.PHONY: clean
clean:
	rm -fr bloggulus c.out dist/
