# bloggulus

[![MIT](https://img.shields.io/github/license/theandrew168/bloggulus)](https://img.shields.io/github/license/theandrew168/bloggulus)

Bloggulus is a web application for aggregating and indexing your favorite blogs.
I wrote it to serve as a less engaging and more personalized version of sites like Hacker News or Reddit.

## Local Development

While the primary [Bloggulus website](https://bloggulus.com) represents my own personal collection of blogs, it is designed to be easily self-hostable.
Check out the [releases page](https://github.com/theandrew168/bloggulus/releases) for pre-built binaries.

### Setup

This project depends on the [Go programming language](https://golang.org/dl/) and the [TailwindCSS CLI](https://tailwindcss.com/blog/standalone-cli).

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

If actively working on frontend templates, set `DEBUG=1` to tell the server to reload templates from the filesystem on every page load.
Run the web server (in one terminal) and let Tailwind watch for CSS changes (in a second terminal):

```bash
# make -j run
DEBUG=1 go run cmd/web/main.go
tailwindcss --watch -m -i tailwind.input.css -o backend/static/static/css/tailwind.min.css
```

### Testing

Tests can be ran after starting the necessary containers and applying database migrations:

```bash
# make test
go run cmd/web/main.go -migrate
go test -v ./...
```

Note that the tests will leave random test in the database so feel free to flush it out by restarting the containers:

```bash
docker compose down
docker compose up -d
```
