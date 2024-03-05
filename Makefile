.POSIX:
.SUFFIXES:

.PHONY: default
default: build

.PHONY: build
build: frontend backend

node_modules:
	npm install

.PHONY: frontend
frontend: node_modules
	npm run build

.PHONY: backend
backend: frontend
	go build -o bloggulus main.go

# run the backend and frontend concurrently (requires at least "-j2") 
.PHONY: run
run: run-frontend run-backend

.PHONY: run-frontend
run-frontend:
	npm run dev

.PHONY: run-backend
run-backend:
	DEBUG=1 go run main.go

.PHONY: migrate
migrate:
	go run main.go -migrate

.PHONY: test
test: migrate
	go test -count=1 ./...

.PHONY: release
release:
	goreleaser release --clean --snapshot

.PHONY: format
format: format-frontend format-backend

.PHONY: format-frontend
format-frontend: node_modules
	npm run format

.PHONY: format-backend
format-backend:
	gofmt -l -s -w .

.PHONY: update
update: update-frontend update-backend

.PHONY: update-frontend
update-frontend:
	npm update --save

.PHONY: update-backend
update-backend:
	go get -u ./...
	go mod tidy

.PHONY: clean
clean:
	rm -fr bloggulus c.out dist/ build/
