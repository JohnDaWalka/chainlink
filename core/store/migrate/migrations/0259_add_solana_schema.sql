-- +goose Up
-- +goose StatementBegin
-- This removes evm from the search path, to avoid risk of an unqualified query intended to act on a Solana table accidentally
-- acting on the corresponding evm table. This makes schema qualification mandatory for all chain-specific queries.
SET search_path TO public;
CREATE SCHEMA solana

CREATE TABLE solana.log_poller_filters (
    id BIGSERIAL
    chain_id TEXT
    name TEXT NOT NULL
    address BYTEA NOT NULL
    event_name TEXT NOT NULL
    event_sig BYTEA NOT NULL
    starting_block BIGINT NOT NULL
    event_idl TEXT
    sub_key_paths TEXT[][] -- A list of subkeys to be indexed, represented by their json paths in the event struct
    retention BIGINT NOT NULL DEFAULT 0 -– we don’t have to implement this initially, but good to include it in the schema
    max_logs_kept BIGINT NOT NULL DEFAULT 0 -- same as retention, no need to implement yet
)

CREATE TABLE solana.logs (
    id               BIGSERIAL
    filter_id        BIGINT NOT NULL REFERENCES solana.log_poller_filters (id) ON DELETE CASCADE
    chain_id         TEXT                      not null
    log_index        bigint                    not null
    block_hash       bytea                     not null
    block_number     bigint                    not null
    block_timestamp  timestamp with time zone  not null
    address          bytea                     not null
    event_sig        bytea                     not null
    subkey_values    bytea[]                   not null
    tx_hash          bytea                     not null
    data             bytea                     not null
    created_at       timestamp with time zone  not null
    expires_at       timestamp with time zone  not null
    sequence_num     bigint                    not null
)

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS solana.logs
DROP TABLE IF EXISTS solana.log_poller_filters

DROP SCHEMA solana;
-- +goose StatementEnd
