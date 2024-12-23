-- +goose Up
ALTER TABLE evm.log_poller_blocks
    ADD COLUMN parent_block_hash bytea;


-- +goose Down
ALTER TABLE evm.log_poller_blocks
    DROP COLUMN parent_block_hash;
