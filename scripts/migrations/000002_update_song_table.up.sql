ALTER TABLE song
ADD COLUMN IF NOT EXISTS image_url TEXT NOT NULL,
ADD COLUMN IF NOT EXISTS duration INT NOT NULL;