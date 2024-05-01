-- +goose Up
-- +goose StatementBegin
ALTER TABLE endpoints RENAME COLUMN ip TO host;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE endpoints RENAME COLUMN host TO ip;
-- +goose StatementEnd
