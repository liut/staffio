
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
	grant_types jsonb NOT NULL DEFAULT '[]'::jsonb,
	response_types jsonb NOT NULL DEFAULT '[]'::jsonb,
	scopes jsonb NOT NULL DEFAULT '[]'::jsonb,
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

CREATE INDEX IF NOT EXISTS idx_access_created ON oauth_access_token (created);
CREATE INDEX IF NOT EXISTS idx_access_refresh ON oauth_access_token (refresh_token);

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

CREATE INDEX idx_authorize_created ON oauth_authorization_code (created);

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


CREATE TABLE IF NOT EXISTS oauth_client_user_authorized
(
	id serial,
	client_id varchar(120) NOT NULL,
	username varchar(120) NOT NULL DEFAULT '',
	created timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
	UNIQUE (client_id, username),
	PRIMARY KEY (id)
);

CREATE SEQUENCE IF NOT EXISTS staff_id_seq START 1027;


CREATE TABLE IF NOT EXISTS articles
(
	id serial,
	title varchar(64) NOT NULL,
	content text NOT NULL,
	author varchar(64) NOT NULL,
	created timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated timestamptz,
	PRIMARY KEY (id)
);

CREATE INDEX IF NOT EXISTS idx_article_created ON articles (created);


CREATE TABLE IF NOT EXISTS links
(
	id serial,
	title varchar(64) NOT NULL,
	url varchar(128) NOT NULL UNIQUE,
	author varchar(64) NOT NULL,
	position smallint NOT NULL DEFAULT 0,
	created timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (id)
);

CREATE INDEX IF NOT EXISTS idx_links_created ON links (created);
CREATE INDEX IF NOT EXISTS idx_links_position ON links (position);


CREATE TABLE IF NOT EXISTS password_reset (
	id serial,
	uid name NOT NULL , -- uid
	type_id smallint NOT NULL, -- 2=email/3=phone
	target varchar(50) NOT NULL , -- phone_number/email_address
	code_hash bigint NOT NULL DEFAULT 0, -- value in crc64
	life_seconds int NOT NULL DEFAULT 3600,
	created timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
	UNIQUE (type_id, target),
	PRIMARY KEY (id)
) WITH (OIDS=FALSE);

CREATE INDEX IF NOT EXISTS idx_password_reset_uid ON password_reset (uid, created);
CREATE INDEX IF NOT EXISTS idx_password_reset_created ON password_reset (created);


CREATE TABLE IF NOT EXISTS user_log (
	id serial,
	uid name NOT NULL , -- uid
	subject name NOT NULL,
	body text NOT NULL DEFAULT '',
	created timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (id)
) WITH (OIDS=FALSE);

CREATE INDEX IF NOT EXISTS idx_user_log_uid ON user_log (uid);


END;
