.POSIX:
.SUFFIXES:

.PHONY: default
default: build

.PHONY: css
css:
	tailwindcss -m -i tailwind.input.css -o internal/static/static/css/tailwind.min.css

.PHONY: build
build: css
	go build -o bloggulus cmd/web/main.go

.PHONY: web
web:
	ENV=dev go run cmd/web/main.go &
	tailwindcss --watch -m -i tailwind.input.css -o internal/static/static/css/tailwind.min.css

.PHONY: migrate
migrate:
	go run cmd/web/main.go -migrate

.PHONY: test
test: migrate
	go test -count=1 ./...

.PHONY: release
release:
	goreleaser release --snapshot --rm-dist

.PHONY: update
update:
	go get -u ./...
	go mod tidy

.PHONY: format
format:
	gofmt -l -s -w .

.PHONY: clean
clean:
	rm -fr bloggulus c.out dist/
