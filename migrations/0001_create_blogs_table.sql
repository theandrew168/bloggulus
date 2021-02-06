CREATE TABLE blogs (
    blog_id SERIAL PRIMARY KEY,
    feed_url TEXT UNIQUE NOT NULL,
    site_url TEXT UNIQUE NOT NULL,
    title TEXT NOT NULL
);
