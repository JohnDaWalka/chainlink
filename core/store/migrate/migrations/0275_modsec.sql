-- +goose Up
-- +goose StatementBegin

-- The modsec_specs table will hold the Modsec job specs.
CREATE TABLE modsec_specs(
    id BIGSERIAL PRIMARY KEY,

    -- The source chain ID.
    source_chain_id TEXT NOT NULL,

    -- The source chain family.
    source_chain_family TEXT NOT NULL,

    -- The destination chain ID.
    dest_chain_id TEXT NOT NULL,

    -- The destination chain family.
    dest_chain_family TEXT NOT NULL,

    -- The onramp address.
    on_ramp_address TEXT NOT NULL,

    -- The ccip message sent event signature.
    ccip_message_sent_event_sig TEXT NOT NULL,

    -- The offramp address.
    off_ramp_address TEXT NOT NULL,

    -- The created at timestamp.
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,

    -- The updated at timestamp.
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Add the modsec_spec_id column to the jobs table
-- and update the chk_specs constraint to include modsec_spec_id
ALTER TABLE jobs
    ADD COLUMN modsec_spec_id INT REFERENCES modsec_specs (id),
DROP CONSTRAINT chk_specs,
    ADD CONSTRAINT chk_specs CHECK (
        num_nonnulls(
            ocr_oracle_spec_id, ocr2_oracle_spec_id,
            direct_request_spec_id, flux_monitor_spec_id,
            keeper_spec_id, cron_spec_id, webhook_spec_id,
            vrf_spec_id, blockhash_store_spec_id,
            block_header_feeder_spec_id, bootstrap_spec_id,
            gateway_spec_id,
            legacy_gas_station_server_spec_id,
            legacy_gas_station_sidecar_spec_id,
            eal_spec_id,
            workflow_spec_id,
            standard_capabilities_spec_id,
            ccip_spec_id,
            ccip_bootstrap_spec_id,
            modsec_spec_id,
            CASE "type"
                WHEN 'stream'
                THEN 1
                ELSE NULL
            END -- 'stream' type lacks a spec but should not cause validation to fail
        ) = 1
    );

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Remove the modsec_spec_id column from the jobs table
-- and update the chk_specs constraint to remove modsec_spec_id
ALTER TABLE jobs
    DROP CONSTRAINT chk_specs,
    ADD CONSTRAINT chk_specs CHECK (
        num_nonnulls(
            ocr_oracle_spec_id, ocr2_oracle_spec_id,
            direct_request_spec_id, flux_monitor_spec_id,
            keeper_spec_id, cron_spec_id, webhook_spec_id,
            vrf_spec_id, blockhash_store_spec_id,
            block_header_feeder_spec_id, bootstrap_spec_id,
            gateway_spec_id,
            legacy_gas_station_server_spec_id,
            legacy_gas_station_sidecar_spec_id,
            eal_spec_id,
            workflow_spec_id,
            standard_capabilities_spec_id,
            ccip_spec_id,
            ccip_bootstrap_spec_id,
            CASE "type"
                WHEN 'stream'
                THEN 1
                ELSE NULL
            END -- 'stream' type lacks a spec but should not cause validation to fail
        ) = 1
    );

ALTER TABLE jobs
    DROP COLUMN modsec_spec_id;

DROP TABLE modsec_specs;

-- +goose StatementEnd
