# bloggulus
A website for avid blog readers

## Database
This project uses [PostgreSQL](https://www.postgresql.org/) for persistent storage.
To develop locally, you'll an instance of the database running somehow or another.
I find [Docker](https://www.docker.com/) to be a nice tool for this but you can do whatever works best.

The following command starts the necessary containers:
```
docker compose up -d
```

These containers can be stopped via:
```
docker compose down
```

## Running
Assuming a recent version of Go is [installed](https://golang.org/dl/), simply run:
```
go run main.go -conf internal/test/bloggulus.conf
```

## Testing
Tests can be ran after starting the necessary containers and applying database migrations:
```
go run main.go -migrate -conf internal/test/bloggulus.conf
go test -v ./...
```

Note that the tests will leave random test in the database so feel free to flush it out by restarting the containers:
```
docker compose down
docker compose up -d
```
