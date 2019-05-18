
CREATE TABLE IF NOT EXISTS http_sessions (
	id BIGSERIAL ,
	key NAME NOT NULL,
	data BYTEA,
	created_on TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
	modified_on TIMESTAMPTZ,
	expires_on TIMESTAMPTZ,
	PRIMARY KEY(id)
);