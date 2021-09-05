CREATE TABLE follow (
    account_id INTEGER REFERENCES account(account_id),
    blog_id INTEGER REFERENCES blog(blog_id),
    PRIMARY KEY (account_id, blog_id)
);

CREATE INDEX follow_account_id_idx ON follow(account_id);
CREATE INDEX follow_blog_id_idx ON follow(blog_id);
