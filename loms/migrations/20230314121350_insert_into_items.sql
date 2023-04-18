-- +goose Up
-- +goose StatementBegin
insert into items (sku) values
    (1076963),
    (1148162),
    (1625903),
    (2618151),
    (2956315),
    (2958025),
    (3596599);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM items
WHERE sku IN (
      1076963,
      1148162,
      1625903,
      2618151,
      2956315,
      2958025,
      3596599,
);
-- +goose StatementEnd
