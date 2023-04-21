DROP PROCEDURE IF EXISTS append_event;

DROP INDEX IF EXISTS stream_events_stream_version_idx;

DROP TABLE IF EXISTS events;

DROP TABLE IF EXISTS streams;

CREATE TABLE streams (
    stream_id CHAR(36) PRIMARY KEY,
    stream_version INTEGER NOT NULL
);

CREATE TABLE events (
    event_id BIGSERIAL PRIMARY KEY,
    stream_id CHAR(36) NOT NULL REFERENCES streams(stream_id),
    stream_version INTEGER NOT NULL,
    event_type TEXT NOT NULL,
    event_encoding TEXT NOT NULL,
    event_source TEXT NOT NULL,
    event_data BYTEA NOT NULL,
    event_ts TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX stream_events_stream_version_idx ON events(stream_id, stream_version);

CREATE PROCEDURE append_event(
    IN p_stream_id CHAR(36),
    IN p_expected_version INTEGER,
    IN p_event_type TEXT,
    IN p_event_encoding TEXT,
    IN p_event_source TEXT,
    IN p_event_data BYTEA
)
LANGUAGE plpgsql
AS $$
DECLARE
  new_version INTEGER;
BEGIN
  BEGIN
    SELECT stream_version INTO STRICT new_version FROM streams WHERE stream_id = p_stream_id FOR UPDATE;
  EXCEPTION
  WHEN NO_DATA_FOUND THEN
    new_version := 0;
    INSERT INTO streams (stream_id, stream_version) VALUES (p_stream_id, new_version);
  END;

  IF new_version = p_expected_version THEN
    new_version := p_expected_version + 1;
    INSERT INTO events (stream_id, stream_version, event_type, event_encoding, event_source, event_data, event_ts)
    VALUES (p_stream_id, new_version, p_event_type, p_event_encoding, p_event_source, p_event_data, NOW());
    UPDATE streams SET stream_version = new_version WHERE stream_id = p_stream_id;
    RETURN;
  END IF;

  RAISE EXCEPTION 'Expected stream_version % but got %', p_expected_version, new_version;
END;
$$;
