CREATE TABLE follows (
    account_id INTEGER REFERENCES account(account_id),
    blog_id INTEGER REFERENCES blog(blog_id),
    PRIMARY KEY (account_id, blog_id)
);

CREATE INDEX follows_account_id_idx ON follows(account_id);
CREATE INDEX follows_blog_id_idx ON follows(blog_id);
