# bloggulus
Custom RSS aggregator and index powered by [Go](https://golang.org/), [PostgreSQL](https://www.postgresql.org/), and [FTS](https://www.postgresql.org/docs/current/textsearch-intro.html).

## Database
This project uses [PostgreSQL](https://www.postgresql.org/) for persistent backend storage.
To develop locally, you'll an instance of PostgreSQL running somehow or another.
I find [Docker](https://www.docker.com/) to be a nice tool for this but you can do whatever works best.

The following commands start the database container and define an environment variable that the app will look for:
```
docker run -e POSTGRES_PASSWORD=postgres -p 5432:5432 --detach postgres
export BLOGGULUS_DATABASE_URL=postgresql://postgres:postgres@localhost:5432/postgres
```

## Running
Assuming a recent version of Go is [installed](https://golang.org/dl/), simply run:
```
go run main.go
```
