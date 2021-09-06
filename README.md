# bloggulus
A community for bloggers and readers

## Database
This project uses [PostgreSQL](https://www.postgresql.org/) for persistent storage and [Redis](https://redis.io/) as a task queue.
To develop locally, you'll an instance of both tools running somehow or another.
I find [Docker](https://www.docker.com/) to be a nice tool for this but you can do whatever works best.

The following commands start the necessary containers and define environment variables that the app will look for:
```
docker compose up -d
export BLOGGULUS_DATABASE_URL=postgresql://postgres:postgres@localhost:5432/postgres
export BLOGGULUS_REDIS_URL=redis://localhost:6397
```

## Running
Assuming a recent version of Go is [installed](https://golang.org/dl/), simply run:
```
go run cmd/migrate/main.go
go run cmd/web/main.go
```
