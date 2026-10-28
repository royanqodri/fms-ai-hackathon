-- UUD
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Schema Operations
CREATE SCHEMA IF NOT EXISTS operations;


CREATE TABLE IF NOT EXISTS operations.trx_telemetries (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    truck_id     UUID NOT NULL, 
    latitude     DECIMAL(10,8),
    longitude    DECIMAL(11,8),
    temperature  DECIMAL(5,2),
    weight_load  DECIMAL(10,2),
    vibration    DECIMAL(5,2),
    recorded_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexing
CREATE INDEX IDX_trx_telemetries_truck ON operations.trx_telemetries(truck_id);
CREATE INDEX IDX_trx_telemetries_time ON operations.trx_telemetries(recorded_at);

CREATE TABLE IF NOT EXISTS operations.log_alerts (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    truck_id    UUID NOT NULL,
    alert_type  VARCHAR(50) NOT NULL, -- 'OVERHEAT', 'OVERLOAD', 'ACCIDENT'
    ai_analysis TEXT,
    is_resolved BOOLEAN DEFAULT FALSE,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IDX_log_alerts_truck ON operations.log_alerts(truck_id);