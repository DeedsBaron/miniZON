-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS cart (
    user_id bigint not null,
    item_sku bigint not null,
    "count" bigint not null,
    primary key (user_id, item_sku)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS cart;
-- +goose StatementEnd
