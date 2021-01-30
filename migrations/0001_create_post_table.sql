CREATE TABLE post (
    post_id SERIAL PRIMARY KEY,
    feed_id INTEGER NOT NULL REFERENCES feed(feed_id),
    title TEXT NOT NULL,
    updated TIMESTAMPTZ NOT NULL
)
