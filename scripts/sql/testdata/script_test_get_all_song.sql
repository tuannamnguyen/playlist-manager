CREATE DATABASE playlist_manager;

\c playlist_manager;

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';

CREATE TABLE IF NOT EXISTS users (
    user_id TEXT PRIMARY KEY,
    user_name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER set_timestamp_users
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

CREATE TABLE IF NOT EXISTS playlist (
    playlist_id SERIAL PRIMARY KEY,
    playlist_name TEXT NOT NULL,
    user_id TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE TRIGGER set_timestamp_playlist
BEFORE UPDATE ON playlist
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

CREATE TABLE IF NOT EXISTS album (
    album_id SERIAL PRIMARY KEY,
    album_name TEXT NOT NULL,
    release_date DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER set_timestamp_album
BEFORE UPDATE ON album
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

CREATE TABLE IF NOT EXISTS artist (
    artist_id SERIAL PRIMARY KEY,
    artist_name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER set_timestamp_artist
BEFORE UPDATE ON artist
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

CREATE TABLE IF NOT EXISTS album_artist (
    artist_id INT NOT NULL,
    album_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (artist_id, album_id),
    FOREIGN KEY (artist_id) REFERENCES artist(artist_id) ON DELETE CASCADE,
    FOREIGN KEY (album_id) REFERENCES album(album_id) ON DELETE CASCADE
);

CREATE TRIGGER set_timestamp_album_artist
BEFORE UPDATE ON album_artist
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

CREATE TABLE IF NOT EXISTS song (
    song_id SERIAL PRIMARY KEY,
    song_name TEXT NOT NULL,
    artist_id TEXT NOT NULL,
    album_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (album_id) REFERENCES album(album_id) ON DELETE CASCADE,
    UNIQUE (song_name, artist_id, album_id)
);

CREATE TRIGGER set_timestamp_song
BEFORE UPDATE ON song
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

CREATE TABLE IF NOT EXISTS playlist_song (
    playlist_id INT NOT NULL,
    song_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (playlist_id, song_id),
    FOREIGN KEY (playlist_id) REFERENCES playlist(playlist_id) ON DELETE CASCADE,
    FOREIGN KEY (song_id) REFERENCES song(song_id) ON DELETE CASCADE
);

CREATE TRIGGER set_timestamp_playlist_song
BEFORE UPDATE ON playlist_song
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

CREATE TABLE IF NOT EXISTS artist_song (
    artist_id INT NOT NULL,
    song_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (artist_id, song_id),
    FOREIGN KEY (artist_id) REFERENCES artist(artist_id) ON DELETE CASCADE,
    FOREIGN KEY (song_id) REFERENCES song(song_id) ON DELETE CASCADE
);

CREATE TRIGGER set_timestamp_artist_song
BEFORE UPDATE ON artist_song
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

-- Insert sample data into the users table
INSERT INTO users (user_id, user_name) VALUES
('google-oauth2|117047339491229984586', 'Nguyen Tuan Nam'),
('user1', 'User One'),
('user2', 'User Two'),
('user3', 'User Three');

-- Insert sample data into the playlist table
INSERT INTO playlist (playlist_id, playlist_name, user_id, updated_at, created_at) VALUES
(1, 'asdasdasvasfasf', 'google-oauth2|117047339491229984586', '2024-07-27 10:12:00', '2024-07-27 10:12:00'),
(2, 'Chill Vibes', 'user1', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(3, 'Workout Hits', 'user2', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(4, 'Classic Rock', 'user3', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Insert sample data into the album table
INSERT INTO album (album_id, album_name, release_date, created_at, updated_at) VALUES
(1, 'mbdtf', '2010-11-22', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(2, 'Album 1', '2020-01-01', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(3, 'Album 2', '2021-06-15', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(4, 'Album 3', '2022-03-22', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Insert sample data into the artist table
INSERT INTO artist (artist_id, artist_name, created_at, updated_at) VALUES
(1, 'kanye west', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(2, 'Artist 1', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(3, 'Artist 2', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(4, 'Artist 3', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Insert sample data into the album_artist table
INSERT INTO album_artist (artist_id, album_id, created_at, updated_at) VALUES
(1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(2, 2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(3, 3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(4, 4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Insert sample data into the song table
INSERT INTO song (song_id, song_name, artist_id, album_id, created_at, updated_at) VALUES
(1, 'devil in a new dress', 1, 1, '2024-07-27 10:12:00', '2024-07-27 10:12:00'),
(2, 'Song 1', 2, 2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(3, 'Song 2', 3, 3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(4, 'Song 3', 4, 4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Insert sample data into the playlist_song table
INSERT INTO playlist_song (playlist_id, song_id, created_at, updated_at) VALUES
(1, 1, '2024-07-27 10:12:00', '2024-07-27 10:12:00'),
(2, 2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(3, 3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(4, 4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Insert sample data into the artist_song table
INSERT INTO artist_song (artist_id, song_id, created_at, updated_at) VALUES
(1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(2, 2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(3, 3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(4, 4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
