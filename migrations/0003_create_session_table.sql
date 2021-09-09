CREATE TABLE session (
    session_id TEXT PRIMARY KEY,
    expiry TIMESTAMPTZ NOT NULL,
    account_id INTEGER NOT NULL REFERENCES account(account_id)
);
