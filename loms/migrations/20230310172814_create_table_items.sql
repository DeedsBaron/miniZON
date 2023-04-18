-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS items (
     sku bigserial primary key not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS items;
-- +goose StatementEnd
