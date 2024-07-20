DROP DATABASE IF EXISTS playlist_manager;

CREATE DATABASE playlist_manager;

\c playlist_manager;

CREATE TABLE IF NOT EXISTS users (
  user_id VARCHAR(50) NOT NULL,
  user_name VARCHAR(50) NOT NULL,
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
  playlist_id VARCHAR(50) NOT NULL,
  playlist_name VARCHAR(50) NOT NULL,
  user_id VARCHAR(50) NOT NULL,
  PRIMARY KEY (playlist_id),
  FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS song (
  song_id VARCHAR(50) NOT NULL,
  song_name VARCHAR(50) NOT NULL,
  artist_id VARCHAR(50) NOT NULL,
  album_id VARCHAR(50) NOT NULL,
  PRIMARY KEY (song_id)
);

CREATE TABLE IF NOT EXISTS album (
  album_id VARCHAR(50) NOT NULL,
  album_name VARCHAR(50) NOT NULL,
  PRIMARY KEY (album_id)
);

CREATE TABLE IF NOT EXISTS artist (
  artist_id VARCHAR(50) NOT NULL,
  artist_name VARCHAR(50) NOT NULL,
  PRIMARY KEY (artist_id)
);

CREATE TABLE IF NOT EXISTS playlist_song (
  playlist_id VARCHAR(50) NOT NULL,
  song_id VARCHAR(50) NOT NULL,
  PRIMARY KEY (playlist_id, song_id),
  FOREIGN KEY (playlist_id) REFERENCES playlist(playlist_id),
  FOREIGN KEY (song_id) REFERENCES song(song_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS album_song (
  song_id VARCHAR(50) NOT NULL,
  album_id VARCHAR(50) NOT NULL,
  PRIMARY KEY (song_id, album_id),
  FOREIGN KEY (song_id) REFERENCES song(song_id),
  FOREIGN KEY (album_id) REFERENCES album(album_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS artist_album (
  artist_id VARCHAR(50) NOT NULL,
  album_id VARCHAR(50) NOT NULL,
  PRIMARY KEY (artist_id, album_id),
  FOREIGN KEY (artist_id) REFERENCES artist(artist_id),
  FOREIGN KEY (album_id) REFERENCES album(album_id) ON DELETE CASCADE
);
