
ALTER TABLE operations.trx_telemetries DROP COLUMN IF EXISTS fuel_level;
ALTER TABLE operations.trx_telemetries DROP COLUMN IF EXISTS oil_pressure;
ALTER TABLE operations.trx_telemetries DROP COLUMN IF EXISTS engine_rpm;