CREATE TABLE blogs (
    blog_id SERIAL PRIMARY KEY,
    feed_url TEXT NOT NULL UNIQUE,
    site_url TEXT NOT NULL UNIQUE,
    title TEXT NOT NULL
);
