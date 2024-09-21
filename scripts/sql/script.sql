CREATE DATABASE playlist_manager;

\c playlist_manager;

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';


CREATE TABLE IF NOT EXISTS playlist (
    playlist_id SERIAL PRIMARY KEY,
    playlist_name TEXT NOT NULL,
    user_id TEXT NOT NULL,
    user_name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (playlist_name, user_id)
);

CREATE TRIGGER set_timestamp_playlist
BEFORE UPDATE ON playlist
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

CREATE TABLE IF NOT EXISTS album (
    album_id SERIAL PRIMARY KEY,
    album_name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (album_name)
);

CREATE TRIGGER set_timestamp_album
BEFORE UPDATE ON album
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

CREATE TABLE IF NOT EXISTS artist (
    artist_id SERIAL PRIMARY KEY,
    artist_name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (artist_name)
);

CREATE TRIGGER set_timestamp_artist
BEFORE UPDATE ON artist
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

CREATE TABLE IF NOT EXISTS artist_album (
    artist_id INT NOT NULL,
    album_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (artist_id, album_id),
    FOREIGN KEY (artist_id) REFERENCES artist(artist_id) ON DELETE CASCADE,
    FOREIGN KEY (album_id) REFERENCES album(album_id) ON DELETE CASCADE
);

CREATE TRIGGER set_timestamp_artist_album
BEFORE UPDATE ON artist_album
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

CREATE TABLE IF NOT EXISTS song (
    song_id SERIAL PRIMARY KEY,
    song_name TEXT NOT NULL,
    album_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (song_name, album_id),
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
    artist_insertion_order INT NOT NULL,
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


-- Insert sample data into the playlist table
INSERT INTO playlist (playlist_name, user_id, user_name) VALUES
('Chill Vibes', 'google-oauth2|117047339491229984586', 'Nguyen Tuan Nam'),
('Workout Hits', 'auth0|6699dbf5833b63066b68e209', 'Nguyen Tuan Nam'),
('Classic Rock', 'google-oauth2|106477541749932989337', 'Nguyen Tuan Nam');

-- Insert sample data into the album table
INSERT INTO album (album_name) VALUES
('Album 1'),
('Album 2'),
('Album 3');

-- Insert sample data into the artist table
INSERT INTO artist (artist_name) VALUES
('Artist 1'),
('Artist 2'),
('Artist 3');

-- Insert sample data into the artist_album table
INSERT INTO artist_album (artist_id, album_id) VALUES
(1, 1),
(2, 2),
(3, 3);

-- Insert sample data into the song table
INSERT INTO song (song_name, album_id) VALUES
('Song 1', 1),
('Song 2', 2),
('Song 3', 3);

-- Insert sample data into the playlist_song table
INSERT INTO playlist_song (playlist_id, song_id) VALUES
(1, 1),
(2, 2),
(3, 3);

-- Insert sample data into the artist_song table
INSERT INTO artist_song (artist_id, song_id, artist_insertion_order) VALUES
(1, 1, 0),
(2, 2, 0),
(3, 3, 0);
