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
    user_id SERIAL PRIMARY KEY,
    username TEXT NOT NULL,
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
    user_id INT NOT NULL,
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
    album_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (album_id) REFERENCES album(album_id) ON DELETE CASCADE
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
INSERT INTO users (username) VALUES
('Nguyen Tuan Nam'),
('user1'),
('user2'),
('user3');

-- Insert sample data into the playlist table
INSERT INTO playlist (playlist_name, user_id) VALUES
('asdasdasvasfasf', 1), -- From the old script
('Chill Vibes', 2),
('Workout Hits', 3),
('Classic Rock', 4);

-- Insert sample data into the album table
INSERT INTO album (album_name, release_date) VALUES
('mbdtf', '2010-11-22'), -- From the old script
('Album 1', '2020-01-01'),
('Album 2', '2021-06-15'),
('Album 3', '2022-03-22');

-- Insert sample data into the artist table
INSERT INTO artist (artist_name) VALUES
('kanye west'), -- From the old script
('Artist 1'),
('Artist 2'),
('Artist 3');

-- Insert sample data into the album_artist table
INSERT INTO album_artist (artist_id, album_id) VALUES
(1, 1), -- Kanye West and MBDTF
(2, 2),
(3, 3),
(4, 4);

-- Insert sample data into the song table
INSERT INTO song (song_name, album_id) VALUES
('devil in a new dress', 1), -- From the old script
('Song 1', 2),
('Song 2', 3),
('Song 3', 4);

-- Insert sample data into the playlist_song table
INSERT INTO playlist_song (playlist_id, song_id) VALUES
(1, 1), -- asdasdasvasfasf playlist and Devil in a New Dress
(2, 2),
(3, 3),
(4, 4);

-- Insert sample data into the artist_song table
INSERT INTO artist_song (artist_id, song_id) VALUES
(1, 1), -- Kanye West and Devil in a New Dress
(2, 2),
(3, 3),
(4, 4);
