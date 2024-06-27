ALTER TABLE blog
ALTER COLUMN synced_at TYPE timestamp,
ALTER COLUMN created_at TYPE timestamp,
ALTER COLUMN updated_at TYPE timestamp;

ALTER TABLE post
ALTER COLUMN published_at TYPE timestamp,
ALTER COLUMN created_at TYPE timestamp,
ALTER COLUMN updated_at TYPE timestamp;

ALTER TABLE tag
ALTER COLUMN created_at TYPE timestamp,
ALTER COLUMN updated_at TYPE timestamp;

ALTER TABLE account
ALTER COLUMN created_at TYPE timestamp,
ALTER COLUMN updated_at TYPE timestamp;

ALTER TABLE token
ALTER COLUMN expires_at TYPE timestamp,
ALTER COLUMN created_at TYPE timestamp,
ALTER COLUMN updated_at TYPE timestamp;
