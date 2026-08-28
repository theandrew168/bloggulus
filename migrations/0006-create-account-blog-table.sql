CREATE TABLE account_blog (
	account_id UUID NOT NULL REFERENCES account(id) ON DELETE CASCADE,
	blog_id UUID NOT NULL REFERENCES blog(id) ON DELETE CASCADE,
	meta_created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	meta_updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	meta_version INTEGER NOT NULL DEFAULT 1,
	PRIMARY KEY (account_id, blog_id)
);