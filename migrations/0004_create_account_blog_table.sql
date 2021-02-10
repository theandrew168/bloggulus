CREATE TABLE account_blog (
    account_id INTEGER REFERENCES account(account_id),
    blog_id INTEGER REFERENCES blog(blog_id),
    PRIMARY KEY (account_id, blog_id)
);

CREATE INDEX account_blog_account_id_idx ON account_blog(account_id);
CREATE INDEX account_blog_blog_id_idx ON account_blog(blog_id);
