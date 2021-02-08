CREATE TABLE posts (
    post_id SERIAL PRIMARY KEY,
    blog_id INTEGER NOT NULL REFERENCES blogs(blog_id),
    url TEXT NOT NULL UNIQUE,
    title TEXT NOT NULL,
    updated TIMESTAMPTZ NOT NULL
);

CREATE INDEX posts_blog_id_idx ON posts(blog_id);
