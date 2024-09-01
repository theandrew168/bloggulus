CREATE TABLE session (
	id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
	account_id UUID NOT NULL REFERENCES account(id) ON DELETE CASCADE,
	hash TEXT NOT NULL UNIQUE,
	expires_at TIMESTAMPTZ NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Used when querying for expired sessions.
CREATE INDEX session_expires_at_idx ON session(expires_at);
