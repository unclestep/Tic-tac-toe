BEGIN;

ALTER TABLE sessions
   ALTER COLUMN winner DROP NOT NULL;

UPDATE sessions SET winner = NULL
WHERE winner NOT IN (SELECT uuid FROM players);

ALTER TABLE sessions
   ADD CONSTRAINT fk_sessions_winner
   FOREIGN KEY (winner) REFERENCES players(uuid);

COMMIT;
