-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS reservation (
    id  bigserial  primary key not null,
    order_id bigint references orders (id) on delete cascade not null,
    item_sku bigint not null,
    warehouse_id bigint not null,

    FOREIGN KEY (item_sku, warehouse_id)
    REFERENCES stocks (item_sku, warehouse_id) on delete cascade,


    reserved_count bigint not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS reservation;
-- +goose StatementEnd
