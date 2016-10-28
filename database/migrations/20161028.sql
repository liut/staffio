
BEGIN;
ALTER TABLE oauth_access_token ALTER access_token TYPE varchar(240), ALTER refresh_token TYPE varchar(240);

END;
