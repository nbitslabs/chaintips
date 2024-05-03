-- +goose Up
-- +goose StatementBegin
CREATE TABLE blocks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    height INTEGER NOT NULL,
    hash TEXT NOT NULL,
    chain_id INTEGER NOT NULL,
    FOREIGN KEY (chain_id) REFERENCES chains(id)
);

CREATE UNIQUE INDEX idx_blocks_unique ON blocks (height, hash, chain_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS blocks;
-- +goose StatementEnd

