-- create new tables w/ UUID PKs
CREATE TABLE blog_new (
	id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
	feed_url TEXT NOT NULL UNIQUE,
	site_url TEXT NOT NULL,
	title TEXT NOT NULL,
	etag TEXT NOT NULL DEFAULT '',
	last_modified TEXT NOT NULL DEFAULT '',
	synced_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	-- temporary
	old_id INTEGER
);

CREATE TABLE post_new (
	id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
	blog_id UUID NOT NULL,
	url TEXT NOT NULL UNIQUE,
	title TEXT NOT NULL,
	content TEXT NOT NULL,
	published_at TIMESTAMPTZ NOT NULL,
	fts_data TSVECTOR GENERATED ALWAYS AS (to_tsvector('english', title || ' ' || content)) STORED,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);


-- fill new tables w/ existing data
INSERT INTO blog_new
	(feed_url, site_url, title, etag, last_modified, old_id)
SELECT
	blog.feed_url,
	blog.site_url,
	blog.title,
	blog.etag,
	blog.last_modified,
	blog.id
FROM blog;

INSERT INTO post_new
	(blog_id, url, title, content, published_at)
SELECT
	blog_new.id,
	post.url,
	post.title,
	post.body,
	post.updated
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
CREATE INDEX post_fts_data_idx ON post USING GIN(fts_data);
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


-- add metadata columns to tag table
ALTER TABLE tag
	ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT now();
