# bloggulus

Bloggulus is a web application for aggregating and indexing your favorite blogs.
I wrote it to serve as a less engaging and more personalized version of sites like Hacker News or Reddit.

## Local Development

While the primary [Bloggulus website](https://bloggulus.com) represents my own personal collection of blogs, it is designed to be easily self-hostable.
Check out the [releases page](https://github.com/theandrew168/bloggulus/releases) for pre-built binaries.

### Setup

This project depends on [Go](https://golang.org/dl/) and the [Tailwind CLI](https://tailwindcss.com/blog/standalone-cli).

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

Run the application and Tailwind CSS watcher at the same time:

```bash
make -j run
```

### Testing

Tests can be ran after starting the necessary containers:

```bash
make test
```
