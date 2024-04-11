.POSIX:
.SUFFIXES:

.PHONY: default
default: build

.PHONY: build
build: frontend backend

frontend/node_modules:
	cd frontend && npm install

.PHONY: frontend
frontend: frontend/node_modules
	cd frontend && npm run build

.PHONY: backend
backend: frontend
	go build -tags embed -o bloggulus main.go

# run the backend and frontend concurrently (requires at least "-j2") 
.PHONY: run
run: run-frontend run-backend

.PHONY: run-frontend
run-frontend:
	cd frontend && npm run dev

.PHONY: run-backend
run-backend:
	DEBUG=1 go run main.go

.PHONY: migrate
migrate:
	go run main.go -migrate

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

.PHONY: format
format: format-frontend format-backend

.PHONY: format-frontend
format-frontend: frontend/node_modules
	cd frontend && npm run format

.PHONY: format-backend
format-backend:
	gofmt -l -s -w .

.PHONY: update
update: update-frontend update-backend

.PHONY: update-frontend
update-frontend:
	cd frontend && npm update --save

.PHONY: update-backend
update-backend:
	go get -u ./...
	go mod tidy

.PHONY: clean
clean:
	rm -fr bloggulus c.out dist/ frontend/build/
