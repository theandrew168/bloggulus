CREATE TABLE tag (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	name TEXT NOT NULL UNIQUE,
	meta_created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	meta_updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	meta_version INTEGER NOT NULL DEFAULT 1
);
