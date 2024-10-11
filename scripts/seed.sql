INSERT INTO account (
	username,
	password_hash,
	is_admin
) VALUES (
	'admin',
	'$2y$10$aJRzokJLFlU8dUG1NND8LO2MTVQXK16Kq2yhC3.WX5x3FX/ygYL..',
	true
);

INSERT INTO blog (
	feed_url,
	site_url,
	title
) VALUES (
	'https://shallowbrooksoftware.com/posts/index.xml',
	'https://shallowbrooksoftware.com/',
	'posts on Shallow Brook Software'
);

INSERT INTO blog (
	feed_url,
	site_url,
	title
) VALUES (
	'https://nickherrig.com/index.xml',
	'https://nickherrig.com/',
	'Nick Herrig'
);

INSERT INTO account_blog
SELECT
	account.id,
	blog.id
FROM account
CROSS JOIN blog;
