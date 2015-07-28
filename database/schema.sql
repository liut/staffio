
BEGIN;

CREATE TABLE oauth_client
(
	id serial,
	code varchar(120) NOT NULL, -- client_id
	name varchar(255) NOT NULL,
	secret varchar(255) NOT NULL,
	redirect_uri varchar(255) NOT NULL DEFAULT '',
	userdata json NOT NULL DEFAULT '{}'::json,
	created timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	UNIQUE (code),
	PRIMARY KEY (id)
);
-- INSERT INTO oauth_client VALUES(1, '1234', 'test client', 'aabbccdd', 'http://localhost:3000/appauth', '{}', now());

CREATE TABLE oauth_access_token
(
	id serial,
	client_id varchar(120) NOT NULL,
	username varchar(120) NOT NULL DEFAULT '',
	access_token varchar(40) NOT NULL,
	expires_in int NOT NULL DEFAULT 86400,
	scopes varchar(255) NOT NULL DEFAULT '',
	is_frozen BOOLEAN NOT NULL DEFAULT false,
	created timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	UNIQUE (access_token),
	PRIMARY KEY (id)
);

CREATE TABLE oauth_authorization_code
(
	id serial,
	code varchar(40) NOT NULL,
	client_id varchar(120) NOT NULL,
	username varchar(120) NOT NULL DEFAULT '',
	redirect_uri varchar(255) NOT NULL DEFAULT '',
	expires_in int NOT NULL DEFAULT 86400,
	scopes varchar(255) NOT NULL DEFAULT '',
	-- token_id int NOT NULL DEFAULT '',
	created timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	UNIQUE (code),
	PRIMARY KEY (id)
);

CREATE TABLE oauth_refresh_token
(
	id serial,
	client_id varchar(120) NOT NULL,
	username varchar(120) NOT NULL DEFAULT '',
	refresh_token varchar(40) NOT NULL,
	expires_in int NOT NULL DEFAULT 86400,
	scopes varchar(255) NOT NULL DEFAULT '',
	is_frozen BOOLEAN NOT NULL DEFAULT false,
	created timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	UNIQUE (refresh_token),
	PRIMARY KEY (refresh_token)
);

CREATE TABLE oauth_scope
(
	scope varchar(255) NOT NULL DEFAULT '',
	is_default BOOLEAN  NOT NULL DEFAULT false,
	PRIMARY KEY (scope)
);


CREATE TABLE oauth_user
(
	id serial,
	username VARCHAR(120) NOT NULL,
	password VARCHAR(200),
	first_name VARCHAR(255),
	last_name VARCHAR(255),
	created timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	UNIQUE (username),
	PRIMARY KEY (id)
);


CREATE TABLE oauth_public_key
(
	id serial,
	client_id varchar(120) NOT NULL,
	public_key TEXT,
	private_key TEXT,
	encryption_algorithm VARCHAR(80) DEFAULT 'RS256',
	created timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (id)
);

END;
