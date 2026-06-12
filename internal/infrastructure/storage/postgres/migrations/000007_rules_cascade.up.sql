ALTER TABLE rules DROP CONSTRAINT rules_session_uuid_fkey;
ALTER TABLE rules ADD CONSTRAINT rules_session_uuid_fkey
    FOREIGN KEY (session_uuid) REFERENCES sessions(uuid) ON DELETE CASCADE;
