DROP DATABASE IF EXISTS playlist_manager;
CREATE DATABASE playlist_manager;

\c playlist_manager;

CREATE TABLE IF NOT EXISTS users
(
  user_id VARCHAR(50) NOT NULL,
  user_name VARCHAR(40) NOT NULL,
  PRIMARY KEY (user_id)
);

INSERT INTO users (user_id, user_name)
VALUES('google-oauth2|117047339491229984586', 'Nguyen Tuan Nam');
