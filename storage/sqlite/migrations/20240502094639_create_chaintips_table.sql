-- +goose Up
-- +goose StatementBegin
CREATE TABLE chaintips (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    chain_id INTEGER NOT NULL,
    endpoint_id INTEGER NOT NULL,
    height INTEGER NOT NULL,
    hash TEXT NOT NULL,
    branchlen INTEGER NOT NULL,
    status TEXT NOT NULL,
    inserted_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (chain_id) REFERENCES chains(id),
    FOREIGN KEY (endpoint_id) REFERENCES endpoints(id)
);

CREATE INDEX idx_chain_id ON chaintips (chain_id);
CREATE INDEX idx_endpoint_id ON chaintips (endpoint_id);
CREATE INDEX idx_height ON chaintips (height);
CREATE INDEX idx_hash ON chaintips (hash);
CREATE INDEX idx_branchlen ON chaintips (branchlen);
CREATE INDEX idx_status ON chaintips (status);
CREATE INDEX idx_inserted_at ON chaintips (inserted_at);
CREATE UNIQUE INDEX idx_chaintips_unique ON chaintips (hash, height, chain_id, endpoint_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS chaintips;
-- +goose StatementEnd

