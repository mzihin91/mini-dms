-- Initial schema for device management system
-- Created: October 28, 2025

-- Device table
CREATE TABLE IF NOT EXISTS devices (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    device_type VARCHAR(50) NOT NULL,
    ip_address INET UNIQUE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'inactive',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    
    CONSTRAINT chk_device_type CHECK (device_type IN ('access_controller', 'face_reader', 'anpr')),
    CONSTRAINT chk_status CHECK (status IN ('active', 'inactive'))
);

-- Transaction table
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    device_id INTEGER NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    username VARCHAR(100),
    event_type VARCHAR(50) NOT NULL,
    payload JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    
    CONSTRAINT fk_device FOREIGN KEY (device_id) 
        REFERENCES devices(id) ON DELETE CASCADE,
    CONSTRAINT chk_timestamp_not_future CHECK (timestamp <= NOW())
);

-- Indexes for devices
CREATE INDEX IF NOT EXISTS idx_devices_status ON devices(status);
CREATE INDEX IF NOT EXISTS idx_devices_type ON devices(device_type);
CREATE UNIQUE INDEX IF NOT EXISTS idx_devices_name_lower ON devices(LOWER(name));

-- Indexes for transactions
CREATE INDEX IF NOT EXISTS idx_transactions_device_id ON transactions(device_id);
CREATE INDEX IF NOT EXISTS idx_transactions_timestamp ON transactions(timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_transactions_event_type ON transactions(event_type);
CREATE INDEX IF NOT EXISTS idx_transactions_payload ON transactions USING GIN(payload);

-- Trigger function for auto-updating updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger for devices table
DROP TRIGGER IF EXISTS update_devices_updated_at ON devices;
CREATE TRIGGER update_devices_updated_at
    BEFORE UPDATE ON devices
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
