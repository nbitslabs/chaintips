-- +goose Up
-- +goose StatementBegin
CREATE TABLE endpoints (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    chain_id INTEGER NOT NULL,
    ip TEXT NOT NULL,
    port INTEGER NOT NULL,
    username TEXT NOT NULL,
    password TEXT NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT 1,
    FOREIGN KEY (chain_id) REFERENCES chains(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS endpoints;
-- +goose StatementEnd

