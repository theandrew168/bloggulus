.POSIX:
.SUFFIXES:

.PHONY: default
default: build

.PHONY: build
build: css
	go build -o bloggulus main.go

.PHONY: css
css:
	tailwindcss -o public/css/tailwind.min.css --minify

# run the application and tailwind watcher concurrently (requires "-j" to work correctly)
.PHONY: run
run: run-app run-css

.PHONY: run-css
run-css:
	tailwindcss -o public/css/tailwind.min.css --minify --watch

.PHONY: run-app
run-app:
	go run main.go

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
update:
	go get -u ./...
	go mod tidy

.PHONY: clean
clean:
	rm -fr bloggulus c.out
