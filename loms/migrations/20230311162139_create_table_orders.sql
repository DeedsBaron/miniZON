-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders (
    id bigserial primary key not null,
    user_id bigint not null,
    status varchar(255) not null,
    created_at timestamp not null
);
CREATE INDEX idx_created_at ON orders (created_at);
CREATE INDEX idx_id ON orders (id);
CREATE INDEX idx_status ON orders (status);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
