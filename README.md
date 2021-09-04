# bloggulus
A community for bloggers and readers

## Database
This project uses [PostgreSQL](https://www.postgresql.org/) for persistent backend storage and [Redis](https://redis.io/) as a task queue.
To develop locally, you'll an instance of both tools running somehow or another.
I find [Docker](https://www.docker.com/) to be a nice tool for this but you can do whatever works best.

The following commands start the necessary containers and define environment variables that the app will look for:
```
docker run -e POSTGRES_PASSWORD=postgres -p 5432:5432 --detach postgres
export BLOGGULUS_DATABASE_URL=postgresql://postgres:postgres@localhost:5432/postgres

docker run -p 6379:6379 --detach redis
export BLOGGULUS_REDIS_ADDR=localhost:6397
```

## Running
Assuming a recent version of Go is [installed](https://golang.org/dl/), simply run:
```
go run main.go
```
