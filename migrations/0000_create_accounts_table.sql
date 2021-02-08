CREATE TABLE accounts (
    account_id SERIAL PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    email TEXT,
    verified BOOLEAN NOT NULL DEFAULT FALSE
);
