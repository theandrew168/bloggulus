CREATE TABLE account_blog (
	account_id UUID NOT NULL REFERENCES account(id) ON DELETE CASCADE,
	blog_id UUID NOT NULL REFERENCES blog(id) ON DELETE CASCADE,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	PRIMARY KEY (account_id, blog_id)
);
