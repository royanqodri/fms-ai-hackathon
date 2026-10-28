-- For UUID
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Schema Master
CREATE SCHEMA IF NOT EXISTS master;

CREATE TABLE IF NOT EXISTS master.mst_customers (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_no VARCHAR(20) NOT NULL UNIQUE,
    name        VARCHAR(100) NOT NULL,
    type        VARCHAR(20) NOT NULL, -- 'COMPANY' OR 'INDIVIDUAL'
    phone_no    VARCHAR(15),
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS master.mst_users (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id  UUID REFERENCES master.mst_customers(id) ON DELETE CASCADE,
    telegram_id  VARCHAR(50) UNIQUE NOT NULL,
    name         VARCHAR(100) NOT NULL,
    role         VARCHAR(20) DEFAULT 'MANAGER',
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS master.mst_trucks (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id  UUID REFERENCES master.mst_customers(id) ON DELETE CASCADE,
    plate_number VARCHAR(20) UNIQUE NOT NULL,
    cargo_type   VARCHAR(50) NOT NULL, -- 'SAND', 'SUGARCANE', 'TIMBER', 'WATERMELON' etc
    status       VARCHAR(20) DEFAULT 'ACTIVE', 
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexing
CREATE INDEX IDX_mst_trucks_plate ON master.mst_trucks(plate_number);
CREATE INDEX IDX_mst_trucks_cargo ON master.mst_trucks(cargo_type);