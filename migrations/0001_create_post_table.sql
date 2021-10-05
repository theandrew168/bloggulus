CREATE TABLE post (
    post_id SERIAL PRIMARY KEY,
    url TEXT NOT NULL UNIQUE,
    title TEXT NOT NULL,
    updated TIMESTAMPTZ NOT NULL,

    body TEXT NOT NULL,
    content_index TSVECTOR
        GENERATED ALWAYS AS (to_tsvector('english', title || ' ' || body)) STORED,

    blog_id INTEGER NOT NULL REFERENCES blog(blog_id)
);

CREATE INDEX post_content_index_idx ON post USING GIN(content_index);
CREATE INDEX post_blog_id_idx ON post(blog_id);
