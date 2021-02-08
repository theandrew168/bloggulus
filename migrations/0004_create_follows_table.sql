CREATE TABLE follows (
    account_id INTEGER NOT NULL REFERENCES accounts(account_id),
    blog_id INTEGER NOT NULL REFERENCES blogs(blog_id)
);
