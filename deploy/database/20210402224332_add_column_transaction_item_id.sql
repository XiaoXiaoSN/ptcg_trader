-- +goose Up
ALTER TABLE trader.transactions	ADD item_id bigint DEFAULT 0 NOT NULL;

-- +goose Down
ALTER TABLE trader.transactions DROP COLUMN item_id;
