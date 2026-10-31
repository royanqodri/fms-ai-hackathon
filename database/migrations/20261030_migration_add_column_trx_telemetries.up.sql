
ALTER TABLE operations.trx_telemetries ADD COLUMN IF NOT EXISTS fuel_level DECIMAL(5,2) DEFAULT 0;
ALTER TABLE operations.trx_telemetries ADD COLUMN IF NOT EXISTS oil_pressure DECIMAL(5,2) DEFAULT 0;
ALTER TABLE operations.trx_telemetries ADD COLUMN IF NOT EXISTS engine_rpm DECIMAL(7,2) DEFAULT 0;