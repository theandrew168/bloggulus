CREATE TABLE post (
    post_id SERIAL PRIMARY KEY,
    feed_id INTEGER NOT NULL REFERENCES feed(feed_id),
    url TEXT UNIQUE NOT NULL,
    title TEXT NOT NULL,
    updated TIMESTAMPTZ NOT NULL
)
