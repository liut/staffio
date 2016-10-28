
BEGIN;

CREATE TABLE IF NOT EXISTS oauth_client
(
	id serial,
	code varchar(80) NOT NULL, -- client_id
	name varchar(120) NOT NULL,
	secret varchar(40) NOT NULL,
	redirect_uri varchar(255) NOT NULL DEFAULT '',
	userdata jsonb NOT NULL DEFAULT '{}'::jsonb,
	created timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	allowed_grant_types varchar(255) DEFAULT '',
	allowed_response_types varchar(255) DEFAULT '',
	allowed_scopes varchar(255) DEFAULT '',
	UNIQUE (code),
	PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS oauth_access_token
(
	id serial,
	client_id varchar(120) NOT NULL,
	username varchar(120) NOT NULL DEFAULT '',
	access_token varchar(240) NOT NULL,
	refresh_token varchar(240) NOT NULL DEFAULT '',
	expires_in int NOT NULL DEFAULT 86400,
	scopes varchar(255) NOT NULL DEFAULT '',
	is_frozen BOOLEAN NOT NULL DEFAULT false,
	created timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
	UNIQUE (access_token),
	PRIMARY KEY (id)
);

CREATE INDEX idx_access_refresh ON oauth_access_token (refresh_token);

CREATE TABLE IF NOT EXISTS oauth_authorization_code
(
	id serial,
	code varchar(140) NOT NULL,
	client_id varchar(120) NOT NULL,
	username varchar(120) NOT NULL DEFAULT '',
	redirect_uri varchar(255) NOT NULL DEFAULT '',
	expires_in int NOT NULL DEFAULT 86400,
	scopes varchar(255) NOT NULL DEFAULT '',
	-- token_id int NOT NULL DEFAULT '',
	created timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
	UNIQUE (code),
	PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS oauth_scope
(
	id serial,
	name varchar(64) NOT NULL, -- ascii code
	label varchar(120) NOT NULL,
	description varchar(255) NOT NULL DEFAULT '',
	is_default BOOLEAN  NOT NULL DEFAULT false,
	UNIQUE (name),
	PRIMARY KEY (id)
);




END;
