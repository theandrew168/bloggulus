CREATE TABLE blog (
    blog_id SERIAL PRIMARY KEY,
    feed_url TEXT NOT NULL UNIQUE,
    site_url TEXT NOT NULL,
    title TEXT NOT NULL
);
