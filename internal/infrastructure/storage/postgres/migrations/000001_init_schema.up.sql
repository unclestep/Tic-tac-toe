CREATE TABLE rules (
      uuid        TEXT PRIMARY KEY,
      board_width  INT NOT NULL,
      board_height INT NOT NULL,
      win_length   INT NOT NULL
);

CREATE TABLE sessions (
    uuid TEXT PRIMARY KEY,
    rules_uuid TEXT NOT NULL REFERENCES rules(uuid),
    board JSONB NOT NULL,
    bots INT NOT NULL DEFAULT 0,
    turn INT NOT NULL DEFAULT 0,
    winner TEXT NOT NULL DEFAULT '',
    state TEXT NOT NULL,
    seed BIGINT NOT NULL DEFAULT 0
);

CREATE TABLE players (
   uuid TEXT,
   session_uuid TEXT NOT NULL REFERENCES sessions(uuid) ON DELETE CASCADE,
   name TEXT NOT NULL,
   mark VARCHAR(1) NOT NULL,
   bot BOOLEAN NOT NULL,
   PRIMARY KEY(session_uuid, mark)
);
