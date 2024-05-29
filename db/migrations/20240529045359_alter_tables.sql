-- migrate:up

ALTER TABLE movie ADD COLUMN title_tsvector TSVECTOR 
  AS (to_tsvector('english', title)) STORED;

CREATE INDEX ON movie USING GIN (title_tsvector);

-- migrate:down
