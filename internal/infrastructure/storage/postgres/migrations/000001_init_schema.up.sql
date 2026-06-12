CREATE TABLE sessions (
    uuid TEXT PRIMARY KEY,
    board JSONB NOT NULL,
    turn INT NOT NULL DEFAULT 0,
    winner TEXT NOT NULL DEFAULT '',
    state TEXT NOT NULL,
    seed BIGINT NOT NULL DEFAULT 0
);

CREATE TABLE rules (
      session_uuid TEXT PRIMARY KEY REFERENCES sessions(uuid),
      board_width  INT NOT NULL,
      board_height INT NOT NULL,
      win_length   INT NOT NULL
);

CREATE TABLE players (
   uuid TEXT PRIMARY KEY,
   session_uuid TEXT NOT NULL REFERENCES sessions(uuid) ON DELETE CASCADE,
   name TEXT NOT NULL,
   mark VARCHAR(1) NOT NULL,
   UNIQUE(session_uuid, mark)
);
