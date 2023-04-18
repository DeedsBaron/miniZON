-- +goose Up
-- +goose StatementBegin
insert into outbox (order_id, old_status, new_status, changed_at, is_sent) values
(1,'none','new','2023-03-31 22:45:41.141054','pending'),
(1,'new','awaiting_payment','2023-03-31 22:45:41.169015','pending'),
(2,'none','new','2023-03-31 22:45:45.519288','pending'),
(2,'new','awaiting_payment','2023-03-31 22:45:45.540105','pending'),
(1,'awaiting_payment','payed','2023-03-31 22:46:03.590366','pending'),
(2,'awaiting_payment','cancelled','2023-03-31 22:46:11.036483','pending'),
(1,'none','new','2023-03-31 22:45:41.141054','pending'),
(1,'new','awaiting_payment','2023-03-31 22:45:41.169015','pending'),
(2,'none','new','2023-03-31 22:45:45.519288','pending'),
(2,'new','awaiting_payment','2023-03-31 22:45:45.540105','pending'),
(1,'awaiting_payment','payed','2023-03-31 22:46:03.590366','pending'),
(2,'awaiting_payment','cancelled','2023-03-31 22:46:11.036483','pending');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
delete from outboux where order_id in (1,2)
-- +goose StatementEnd
