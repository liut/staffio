
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
