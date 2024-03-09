-- create new tables w/ UUID PKs
CREATE TABLE blog_new (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    feed_url TEXT NOT NULL UNIQUE,
    site_url TEXT NOT NULL,
    title TEXT NOT NULL,
	old_id INTEGER
);

CREATE TABLE post_new (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    url TEXT NOT NULL UNIQUE,
    title TEXT NOT NULL,
    updated TIMESTAMPTZ NOT NULL,

    body TEXT NOT NULL,
    content_index TSVECTOR
        GENERATED ALWAYS AS (to_tsvector('english', title || ' ' || body)) STORED,

    blog_id UUID NOT NULL
);


-- fill new tables w/ existing data
INSERT INTO blog_new (feed_url, site_url, title, old_id)
SELECT blog.feed_url, blog.site_url, blog.title, blog.id
FROM blog;

INSERT INTO post_new (url, title, updated, body, blog_id)
SELECT post.url, post.title, post.updated, post.body, blog_new.id
FROM post
INNER JOIN blog_new
	ON blog_new.old_id = post.blog_id;


-- drop temp cols from new tables
ALTER TABLE blog_new DROP COLUMN old_id;


-- drop old tables
DROP TABLE post;
DROP TABLE blog;


-- rename new tables
ALTER TABLE blog_new RENAME TO blog;
ALTER TABLE post_new RENAME TO post;


-- rename new table indexes
ALTER INDEX blog_new_pkey RENAME TO blog_pkey;
ALTER INDEX blog_new_feed_url_key RENAME TO blog_feed_url_key;
ALTER INDEX post_new_pkey RENAME TO post_pkey;
ALTER INDEX post_new_url_key RENAME TO post_url_key;


-- add indexes to new tables
CREATE INDEX post_content_index_idx ON post USING GIN(content_index);
CREATE INDEX post_blog_id_idx ON post(blog_id);


-- add FK constraint from post->blog
ALTER TABLE post ADD CONSTRAINT blog_id_fkey FOREIGN KEY (blog_id) REFERENCES blog(id) ON DELETE CASCADE;


-- update other tables in place
ALTER TABLE tag
	ALTER COLUMN id DROP DEFAULT, 
	ALTER COLUMN id SET DATA TYPE UUID USING (gen_random_uuid()), 
	ALTER COLUMN id SET DEFAULT gen_random_uuid();

ALTER TABLE migration
	ALTER COLUMN id DROP DEFAULT, 
	ALTER COLUMN id SET DATA TYPE UUID USING (gen_random_uuid()), 
	ALTER COLUMN id SET DEFAULT gen_random_uuid();
