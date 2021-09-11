CREATE TABLE account (
    account_id SERIAL PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    email TEXT NOT NULL,
    verified BOOLEAN NOT NULL
);
