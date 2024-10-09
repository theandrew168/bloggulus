CREATE TABLE account_page (
	account_id UUID NOT NULL REFERENCES account(id) ON DELETE CASCADE,
	page_id UUID NOT NULL REFERENCES page(id) ON DELETE CASCADE,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	PRIMARY KEY (account_id, page_id)
);
