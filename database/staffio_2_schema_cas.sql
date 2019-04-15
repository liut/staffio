
CREATE TABLE IF NOT EXISTS cas_ticket (
	id serial,
	type VARCHAR(5) NOT NULL,
	uid NAME NOT NULL , -- uid
	value VARCHAR(139) NOT NULL, -- Ticket value
	service VARCHAR(200) NOT NULL DEFAULT '',
	created timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (id)
) WITH (OIDS=FALSE);

CREATE INDEX IF NOT EXISTS idx_cas_ticket_uid ON cas_ticket (uid);
CREATE INDEX IF NOT EXISTS idx_cas_ticket_created ON cas_ticket (created);
CREATE INDEX IF NOT EXISTS idx_cas_ticket_value ON cas_ticket (value);
