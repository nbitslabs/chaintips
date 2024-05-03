-- +goose Up
ALTER TABLE blocks ADD COLUMN version TEXT;
ALTER TABLE blocks ADD COLUMN merkleroot TEXT;
ALTER TABLE blocks ADD COLUMN time TEXT;
ALTER TABLE blocks ADD COLUMN mediantime TEXT;
ALTER TABLE blocks ADD COLUMN nonce TEXT;
ALTER TABLE blocks ADD COLUMN bits TEXT;
ALTER TABLE blocks ADD COLUMN difficulty TEXT;
ALTER TABLE blocks ADD COLUMN chainwork TEXT;
ALTER TABLE blocks ADD COLUMN previousblockhash TEXT;

-- +goose Down
-- SQLite does not support dropping a column via ALTER TABLE
-- Rollback would require a more complex strategy if needed

