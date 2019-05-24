BEGIN;
ALTER TABLE oauth_client
	ADD grant_types jsonb NOT NULL  DEFAULT '[]'::jsonb,
	ADD response_types jsonb NOT NULL  DEFAULT '[]'::jsonb,
	ADD scopes jsonb NOT NULL  DEFAULT '[]'::jsonb;

UPDATE oauth_client SET
  grant_types = to_jsonb(string_to_array(allowed_grant_types, ',')),
  response_types = to_jsonb(string_to_array(allowed_response_types, ',')),
  scopes = to_jsonb(string_to_array(allowed_scopes, ','))
;

ALTER TABLE oauth_client
  DROP allowed_grant_types,
  DROP allowed_response_types,
  DROP allowed_scopes
;
END;
