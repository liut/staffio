-- init

INSERT INTO oauth_scope(name,label,description,is_default) VALUES('basic', 'Basic', 'Read your Uid (login name) and Nickname', true);
INSERT INTO oauth_scope(name,label,description) VALUES('openid', 'OpenID Connect', 'Read your ID Token after authenticated.');
INSERT INTO oauth_scope(name,label,description) VALUES('profile', 'Personal Information', 'Read your GivenName, Surname, Email, etc.');
