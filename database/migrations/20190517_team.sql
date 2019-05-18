
ALTER TABLE teams
   DROP COLUMN leader,
   ADD COLUMN leaders jsonb NOT NULL DEFAULT '[]'::jsonb;
