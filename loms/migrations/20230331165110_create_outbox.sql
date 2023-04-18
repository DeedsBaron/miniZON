-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS outbox (
    id bigserial primary key not null,
    order_id bigint not null,
    old_status varchar(255),
    new_status varchar(255) not null,
    changed_at timestamp not null,
    is_sent varchar(10) default 'pending'
);
CREATE INDEX idx_changed_at ON outbox (changed_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS outbox;
-- +goose StatementEnd
