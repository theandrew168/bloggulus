.POSIX:
.SUFFIXES:

.PHONY: default
default: build

.PHONY: css
css:
	tailwindcss -m -i tailwind.input.css -o static/css/tailwind.min.css

.PHONY: build
build: css
	go build -o bloggulus main.go

.PHONY: run
run:
	ENV=dev go run main.go &
	tailwindcss --watch -m -i tailwind.input.css -o static/css/tailwind.min.css

.PHONY: migrate
migrate:
	go run main.go -migrate

.PHONY: test
test: migrate
	go test -count=1 ./...

.PHONY: race
race: migrate
	go test -race -count=1 ./...

.PHONY: cover
cover: migrate
	go test -coverprofile=c.out -coverpkg=./... -count=1 ./...
	go tool cover -html=c.out

.PHONY: release
release:
	goreleaser release --snapshot --rm-dist

.PHONY: format
format:
	go fmt ./...

.PHONY: clean
clean:
	rm -fr bloggulus c.out dist/
