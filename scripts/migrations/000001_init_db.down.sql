-- Drop triggers in reverse order of creation
DROP TRIGGER IF EXISTS set_timestamp_artist_song ON artist_song;
DROP TRIGGER IF EXISTS set_timestamp_playlist_song ON playlist_song;
DROP TRIGGER IF EXISTS set_timestamp_song ON song;
DROP TRIGGER IF EXISTS set_timestamp_artist_album ON artist_album;
DROP TRIGGER IF EXISTS set_timestamp_artist ON artist;
DROP TRIGGER IF EXISTS set_timestamp_album ON album;
DROP TRIGGER IF EXISTS set_timestamp_playlist ON playlist;

-- Drop tables in reverse order to maintain referential integrity
DROP TABLE IF EXISTS artist_song;
DROP TABLE IF EXISTS playlist_song;
DROP TABLE IF EXISTS song;
DROP TABLE IF EXISTS artist_album;
DROP TABLE IF EXISTS artist;
DROP TABLE IF EXISTS album;
DROP TABLE IF EXISTS playlist;

-- Drop the function used by the triggers
DROP FUNCTION IF EXISTS update_updated_at_column;
