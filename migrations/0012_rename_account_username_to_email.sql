ALTER TABLE account
RENAME COLUMN username TO email;

ALTER INDEX account_username_key
RENAME TO account_email_key;
