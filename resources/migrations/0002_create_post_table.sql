CREATE TABLE post (
    post_id SERIAL PRIMARY KEY,
    blog_id INTEGER NOT NULL REFERENCES blog(blog_id),
    url TEXT NOT NULL UNIQUE,
    title TEXT NOT NULL,
    preview TEXT NOT NULL DEFAULT '',
    updated TIMESTAMPTZ NOT NULL
);

CREATE INDEX post_blog_id_idx ON post(blog_id);
