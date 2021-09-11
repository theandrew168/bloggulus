CREATE TABLE post (
    post_id SERIAL PRIMARY KEY,
    url TEXT NOT NULL UNIQUE,
    title TEXT NOT NULL,
    author TEXT NOT NULL,
    body TEXT NOT NULL,
    updated TIMESTAMPTZ NOT NULL,

    body_index TSVECTOR GENERATED ALWAYS AS
        (TO_TSVECTOR('english', title || ' ' || author || ' ' || body)) STORED,
    blog_id INTEGER NOT NULL REFERENCES blog(blog_id)
);

CREATE INDEX post_body_index_idx ON post USING GIN(body_index);
CREATE INDEX post_blog_id_idx ON post(blog_id);
