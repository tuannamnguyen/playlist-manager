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
INSERT INTO song (song_name, album_id, image_url, duration) VALUES
('Song 1', 1, 'https://picsum.photos/200/300', 180),
('Song 2', 2, 'https://picsum.photos/200/300', 180),
('Song 3', 3, 'https://picsum.photos/200/300', 180);

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
