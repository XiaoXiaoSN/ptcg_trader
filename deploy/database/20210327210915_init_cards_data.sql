-- +goose Up
INSERT INTO trader.items (id, name, image_url, creator_id, updated_at, created_at) VALUES (DEFAULT, 'Pikachu', 'https://imgur.com/NTSEJxX', 0, DEFAULT, DEFAULT);
INSERT INTO trader.items (id, name, image_url, creator_id, updated_at, created_at) VALUES (DEFAULT, 'Charmander', 'https://imgur.com/JNyYH8F', 0, DEFAULT, DEFAULT);
INSERT INTO trader.items (id, name, image_url, creator_id, updated_at, created_at) VALUES (DEFAULT, 'Bulbasaur', 'https://imgur.com/KqZO7Wv', 0, DEFAULT, DEFAULT);
INSERT INTO trader.items (id, name, image_url, creator_id, updated_at, created_at) VALUES (DEFAULT, 'Squirtle', 'https://imgur.com/sOX4WkY', 0, DEFAULT, DEFAULT);

-- +goose Down
DELETE FROM trader.items WHERE name IN ('Pikachu', 'Charmander', 'Bulbasaur', 'Squirtle');
