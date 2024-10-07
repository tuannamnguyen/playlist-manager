ALTER TABLE playlist
ADD CONSTRAINT playlist_playlist_name_user_id_key UNIQUE (playlist_name, user_id);
