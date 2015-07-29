-- init

INSERT INTO oauth_scope(name,label,description,is_default) VALUES('basic', 'Basic', 'Read your Uid (login name)', true);
INSERT INTO oauth_scope(name,label,description) VALUES('user_info', 'Personal Information', 'Read your GivenName, Surname, Email, etc.');
