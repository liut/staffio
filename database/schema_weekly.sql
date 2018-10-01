
CREATE TABLE IF NOT EXISTS  weekly_report (
	id serial,
	uid NAME NOT NULL,
	iso_year smallint  NOT NULL,
	iso_week smallint  NOT NULL,
	content jsonb NOT NULL,
	up_count int NOT NULL DEFAULT 0,
	created timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	UNIQUE(uid,iso_year,iso_week),
	PRIMARY KEY (id)
);

CREATE INDEX IF NOT EXISTS idx_weekly_report_yw ON weekly_report (iso_year, iso_week);
CREATE INDEX IF NOT EXISTS idx_weekly_report_created ON weekly_report (created);
CREATE INDEX IF NOT EXISTS idx_weekly_report_up_count ON weekly_report (up_count, id);

CREATE TABLE IF NOT EXISTS  weekly_report_up (
	id serial,
	report_id INT  NOT NULL,
	uid NAME NOT NULL,
	created timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	UNIQUE(report_id, uid),
	FOREIGN KEY (report_id) REFERENCES weekly_report (id),
	PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS  weekly_problem (
	id serial,
	uid NAME NOT NULL,
	iso_year smallint  NOT NULL,
	iso_week smallint  NOT NULL,
	content TEXT NOT NULL,
	created timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	UNIQUE(uid,iso_year,iso_week),
	PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS  weekly_status (
	id serial,
	uid NAME NOT NULL,
	iso_year smallint  NOT NULL,
	iso_week smallint  NOT NULL,
	status smallint NOT NULL, -- 1=vacation，2=ignore,
	created timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	UNIQUE(uid, iso_year, iso_week),
	PRIMARY KEY (id)
);
