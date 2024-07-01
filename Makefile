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

.PHONY: frontend-watch
frontend-watch: frontend/node_modules
	cd frontend && npm run watch

.PHONY: backend
backend: frontend
	go build -tags embed -o bloggulus main.go

# run the backend and frontend concurrently (requires "-j" to work correctly)
.PHONY: run
run: run-frontend run-backend

.PHONY: run-frontend
run-frontend:
	cd frontend && npm run dev

.PHONY: run-backend
run-backend:
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
