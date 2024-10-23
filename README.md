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

### OAuth Services

For authentication, this project relies on OAuth social sign ins (from GitHub and Google).
To work on the auth system, you'll need to create a `bloggulus.local.conf` file that contains the necessary OAuth credentials for each service.
If you need these credentials, feel free to reach out.

Then, you can run the app using the local config file with:

```bash
make run-local
```

Otherwise, you can simply run the application normally (without OAuth configured) and use the local-only debug sign in.
This is automatically enabled (via the `ENABLE_DEBUG_AUTH` environment variable) for `run` and `run-local`.

### Testing

Tests can be ran after starting the necessary containers:

```bash
make test
```
