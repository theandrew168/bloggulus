CREATE TABLE account (
	id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
	username TEXT NOT NULL UNIQUE,
	is_admin BOOLEAN NOT NULL DEFAULT FALSE,
	meta_created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	meta_updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	meta_version INTEGER NOT NULL DEFAULT 1
);