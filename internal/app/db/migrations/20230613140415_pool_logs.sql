-- +goose Up
-- +goose StatementBegin
CREATE TABLE public.pool_logs
(
    id bigserial NOT NULL,
    pool_address text,
    txn_id text,
    block_number bigint,
    token0_balance bigint,
    token1_balance bigint,
    token0_delta bigint,
    token1_delta bigint,
    tick bigint,
    created_at timestamp without time zone NOT NULL DEFAULT 'NOW()',
    PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
