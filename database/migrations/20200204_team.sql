
ALTER TABLE teams
   DROP CONSTRAINT IF EXISTS teams_name_key,
   ADD COLUMN IF NOT EXISTS parent_id int NOT NULL DEFAULT 0;

ALTER TABLE teams
	ADD UNIQUE (parent_id, name);
