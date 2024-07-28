DROP DATABASE IF EXISTS playlist_manager;

CREATE DATABASE playlist_manager;

\c playlist_manager;

CREATE TABLE IF NOT EXISTS users (
  user_id TEXT NOT NULL,
  user_name TEXT NOT NULL,
  PRIMARY KEY (user_id)
);

INSERT INTO
  users (user_id, user_name)
VALUES
  (
    'google-oauth2|117047339491229984586',
    'Nguyen Tuan Nam'
  );

CREATE TABLE IF NOT EXISTS playlist (
  playlist_id TEXT NOT NULL,
  playlist_name TEXT NOT NULL,
  user_id TEXT NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  created_at TIMESTAMP NOT NULL,
  PRIMARY KEY (playlist_id),
  FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS song (
  song_id TEXT NOT NULL,
  song_name TEXT NOT NULL,
  artist_id TEXT NOT NULL,
  album_id TEXT NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  created_at TIMESTAMP NOT NULL,
  PRIMARY KEY (song_id)
);


CREATE TABLE IF NOT EXISTS playlist_song (
  playlist_id TEXT NOT NULL,
  song_id TEXT NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  created_at TIMESTAMP NOT NULL,
  PRIMARY KEY (playlist_id, song_id),
  FOREIGN KEY (playlist_id) REFERENCES playlist(playlist_id) ON DELETE CASCADE,
  FOREIGN KEY (song_id) REFERENCES song(song_id) ON DELETE CASCADE
);

INSERT INTO playlist (playlist_id, playlist_name, user_id, updated_at, created_at)
VALUES ('asdasdasdsaasd', 'asdasdasvasfasf', 'google-oauth2|117047339491229984586', '2024-07-27 10:12:00', '2024-07-27 10:12:00');

INSERT INTO song (song_id, song_name, artist_id, album_id, updated_at, created_at)
VALUES ('asiuasubfasuifaufb', 'devil in a new dress', 'kanye west', 'mbdtf', '2024-07-27 10:12:00', '2024-07-27 10:12:00');

INSERT INTO playlist_song (playlist_id, song_id, updated_at, created_at)
VALUES ('asdasdasdsaasd', 'asiuasubfasuifaufb', '2024-07-27 10:12:00', '2024-07-27 10:12:00');
