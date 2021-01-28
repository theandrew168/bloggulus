CREATE VIRTUAL TABLE post_index USING fts5 (
    feed,
    title,
    content
)
