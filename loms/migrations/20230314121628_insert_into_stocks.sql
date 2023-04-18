-- +goose Up
-- +goose StatementBegin
insert into stocks (item_sku, warehouse_id, count) values
(1076963,221, 100),
(1076963,220, 100),
(1076963,223, 100),
(1148162,220, 100),
(1625903,221, 100),
(2618151,223, 100),
(2956315,220, 100),
(2958025,224, 100),
(3596599,225, 100);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM stocks
WHERE (item_sku = 1076963 and warehouse_id = 221) or
    (item_sku = 1076963 and warehouse_id = 220) or
    (item_sku = 1076963 and warehouse_id = 223) or
    (item_sku = 1148162 and warehouse_id = 220) or
    (item_sku = 1625903 and warehouse_id = 221) or
    (item_sku = 2618151 and warehouse_id = 223) or
    (item_sku = 2956315 and warehouse_id = 220) or
    (item_sku = 2958025 and warehouse_id = 224) or
    (item_sku = 3596599 and warehouse_id = 225);
-- +goose StatementEnd
