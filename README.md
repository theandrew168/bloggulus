# bloggulus

Bloggulus is a web application for aggregating and indexing your favorite blogs.
I wrote it to serve as a less engaging and more personalized version of sites like Hacker News or Reddit.

## Local Development

While the primary [Bloggulus website](https://bloggulus.com) represents my own personal collection of blogs, it is designed to be easily self-hostable.
Check out the [releases page](https://github.com/theandrew168/bloggulus/releases) for pre-built binaries.

### Setup

This project depends on the [Go](https://golang.org/dl/) programming language.

### Database

This project uses [PostgreSQL](https://www.postgresql.org/) for persistent storage.
To develop locally, you'll an instance of the database running somehow or another.
I find [Docker](https://www.docker.com/) to be a nice tool for this but you can do whatever works best.

The following command starts the necessary containers:

```bash
docker compose up -d
```

These containers can be stopped via:

```bash
docker compose down
```

### Running

Run the application (with automatic restarts via [wgo](https://github.com/bokwoon95/wgo)):

```bash
make run
```

### Testing

Tests can be ran after starting the necessary containers:

```bash
make test
```
