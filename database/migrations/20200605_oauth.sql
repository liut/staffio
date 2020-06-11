BEGIN;

ALTER TABLE oauth_client
	RENAME userdata to meta;

UPDATE oauth_client SET meta = json_build_object('name', name,
	'grant_types', grant_types,
	'response_types', response_types,
	'scopes', scopes);

ALTER TABLE oauth_client
  DROP name,
  DROP grant_types,
  DROP response_types,
  DROP scopes
;

ALTER TABLE oauth_client ALTER id SET DATA TYPE varchar(30);
UPDATE oauth_client SET id = code;
ALTER TABLE oauth_client DROP code;
ALTER TABLE oauth_client ALTER id DROP DEFAULT;

ALTER TABLE oauth_access_token
	ADD userdata jsonb NOT NULL DEFAULT '{}'::jsonb;

UPDATE oauth_access_token SET userdata = json_build_object('name', username);

ALTER TABLE oauth_access_token DROP COLUMN IF EXISTS username;

ALTER TABLE oauth_authorization_code
	ADD userdata jsonb NOT NULL DEFAULT '{}'::jsonb;

UPDATE oauth_authorization_code SET userdata = json_build_object('name', username);

ALTER TABLE oauth_authorization_code DROP COLUMN IF EXISTS username;

ALTER TABLE oauth_access_token
	ADD COLUMN authorize_code varchar(140) NOT NULL DEFAULT '',
	ADD COLUMN previous varchar(240) NOT NULL DEFAULT '';

END;
