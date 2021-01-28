CREATE TABLE post (
    post_id INTEGER PRIMARY KEY NOT NULL,
    feed_id INTEGER NOT NULL REFERENCES feed(feed_id),
    title TEXT NOT NULL,
    updated DATETIME NOT NULL    
)
