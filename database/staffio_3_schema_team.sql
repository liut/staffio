
-- sync from wxqiye(exwechat)
CREATE TABLE IF NOT EXISTS "department" (
  id SERIAL,
  name VARCHAR(90) NOT NULL,
  parent_id int NOT NULL DEFAULT 0,
  position int NOT NULL DEFAULT 0,
  created timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated timestamptz ,
  UNIQUE (parent_id, name),
  PRIMARY KEY (id)
);


CREATE TABLE IF NOT EXISTS teams (
	id serial,
	name VARCHAR(120) NOT NULL UNIQUE,
	leaders jsonb NOT NULL, -- leader uid
	members jsonb NOT NULL DEFAULT '[]'::jsonb,
	created timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS team_leader (
	id serial,
	team_id INT  NOT NULL,
	leader NAME NOT NULL, -- leader uid
	created timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	UNIQUE(team_id, leader),
	FOREIGN KEY (team_id) REFERENCES teams (id),
	PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS team_member (
	id serial,
	team_id INT  NOT NULL,
	uid NAME NOT NULL,
	created timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	UNIQUE(team_id, uid),
	FOREIGN KEY (team_id) REFERENCES teams (id),
	PRIMARY KEY (id)
);

