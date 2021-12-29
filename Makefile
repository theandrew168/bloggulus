.POSIX:
.SUFFIXES:

.PHONY: default
default: build

.PHONY: css
css:
	tailwindcss -m -i static/css/tailwind.input.css -o static/css/tailwind.min.css

.PHONY: build
build: css
	go build -o bloggulus main.go

.PHONY: run
run:
	go run main.go -conf internal/test/bloggulus.conf &
	tailwindcss --watch -m -i static/css/tailwind.input.css -o static/css/tailwind.min.css

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
