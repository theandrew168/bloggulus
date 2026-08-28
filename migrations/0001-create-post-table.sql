CREATE TABLE post (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	blog_id uuid NOT NULL REFERENCES blog(id) ON DELETE CASCADE,
	url TEXT NOT NULL UNIQUE,
	title TEXT NOT NULL,
	published_at TIMESTAMPTZ NOT NULL,
	content TEXT,
	fts_data TSVECTOR GENERATED ALWAYS AS (to_tsvector('english', title || ' ' || content)) STORED,
	meta_created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	meta_updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	meta_version INTEGER NOT NULL DEFAULT 1
);

-- Used when querying for posts by blog.
CREATE INDEX post_blog_id_idx ON post(blog_id);
-- Used when searching for posts.
CREATE INDEX post_fts_data_idx ON post USING GIN(fts_data);
