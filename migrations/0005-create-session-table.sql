CREATE TABLE session (
	id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
	account_id UUID NOT NULL REFERENCES account(id) ON DELETE CASCADE,
	token_hash TEXT NOT NULL UNIQUE,
	expires_at TIMESTAMPTZ NOT NULL,
	meta_created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	meta_updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	meta_version INTEGER NOT NULL DEFAULT 1
);

-- Used when querying for expired sessions.
CREATE INDEX session_expires_at_idx ON session(expires_at);