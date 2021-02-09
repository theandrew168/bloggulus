CREATE TABLE session (
    session_id TEXT PRIMARY KEY,
    account_id INTEGER NOT NULL REFERENCES account(account_id),
    expiry TIMESTAMPTZ NOT NULL
);
