CREATE TABLE page (
	id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
	url TEXT NOT NULL UNIQUE,
	title TEXT NOT NULL,
	content TEXT NOT NULL,
	fts_data TSVECTOR GENERATED ALWAYS AS (to_tsvector('english', title || ' ' || content)) STORED,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX page_fts_data_idx ON page USING GIN(fts_data);
