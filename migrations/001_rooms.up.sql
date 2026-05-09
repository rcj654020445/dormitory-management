-- Migration: Create rooms table
-- Up

CREATE TABLE IF NOT EXISTS rooms (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    building_id   UUID NOT NULL REFERENCES buildings(id) ON DELETE CASCADE,
    number        VARCHAR(20) NOT NULL,
    floor         INTEGER NOT NULL DEFAULT 1,
    type          VARCHAR(20) NOT NULL DEFAULT 'standard',
    beds_total    INTEGER NOT NULL DEFAULT 4,
    beds_used     INTEGER NOT NULL DEFAULT 0,
    has_bathroom  BOOLEAN NOT NULL DEFAULT false,
    has_ac        BOOLEAN NOT NULL DEFAULT false,
    status        VARCHAR(20) NOT NULL DEFAULT 'available',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ DEFAULT NULL
);

CREATE UNIQUE INDEX idx_rooms_building_number ON rooms(building_id, number) WHERE deleted_at IS NULL;
CREATE INDEX idx_rooms_building_id ON rooms(building_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_rooms_status ON rooms(status) WHERE deleted_at IS NULL;

-- Trigger to auto-update updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_rooms_updated_at
    BEFORE UPDATE ON rooms
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
