CREATE TABLE users (
   uuid TEXT PRIMARY KEY,
   login TEXT UNIQUE NOT NULL,
   password TEXT NOT NULL
);
