-- +goose Up
ALTER TABLE endpoints ADD COLUMN protocol TEXT NOT NULL DEFAULT 'http';

-- +goose Down
-- SQLite does not support dropping a column via ALTER TABLE
-- To roll back this migration, you would need to create a new table without the protocol column, copy data over, and drop the old table

