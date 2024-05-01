-- +goose Up
-- +goose StatementBegin
INSERT INTO chains (identifier, title, icon)
VALUES ('dogecoin-mainnet', 'Dogecoin', 'dogecoin.svg');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM chains WHERE chain_identifier = 'dogecoin-mainnet';
-- +goose StatementEnd

