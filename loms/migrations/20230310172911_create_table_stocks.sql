-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS stocks (
    item_sku bigint references items (sku) on delete cascade not null,
    warehouse_id bigint not null,
    "count" bigint not null,
    primary key (item_sku, warehouse_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS stocks;
-- +goose StatementEnd
