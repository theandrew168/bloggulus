CREATE TABLE sessions (
    session_id TEXT PRIMARY KEY,
    account_id INTEGER NOT NULL REFERENCES accounts(account_id),
    expiry TIMESTAMPTZ NOT NULL
);
