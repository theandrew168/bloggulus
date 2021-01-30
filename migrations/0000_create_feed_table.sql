CREATE TABLE feed (
    feed_id SERIAL PRIMARY KEY,
    url TEXT UNIQUE NOT NULL,
    site_url TEXT NOT NULL,
    title TEXT NOT NULL
)
