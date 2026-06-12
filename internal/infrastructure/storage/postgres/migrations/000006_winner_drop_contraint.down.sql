ALTER TABLE sessions
   ADD CONSTRAINT fk_sessions_winner
   FOREIGN KEY (winner) REFERENCES players(uuid);
