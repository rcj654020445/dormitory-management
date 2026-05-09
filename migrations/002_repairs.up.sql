-- Migration: Create repairs table
-- Up

CREATE TABLE IF NOT EXISTS repairs (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id       UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    reporter_id   UUID NOT NULL,
    repairer_id   UUID DEFAULT NULL,
    type          VARCHAR(20) NOT NULL DEFAULT 'plumbing',
    description   TEXT NOT NULL,
    status        VARCHAR(20) NOT NULL DEFAULT 'pending',
    priority      VARCHAR(10) NOT NULL DEFAULT 'normal',
    scheduled_at  TIMESTAMPTZ DEFAULT NULL,
    completed_at  TIMESTAMPTZ DEFAULT NULL,
    cost          DECIMAL(10,2) DEFAULT NULL,
    rating        INTEGER DEFAULT NULL CHECK (rating >= 1 AND rating <= 5),
    remark        TEXT DEFAULT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ DEFAULT NULL
);

CREATE INDEX idx_repairs_room_id ON repairs(room_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_repairs_reporter_id ON repairs(reporter_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_repairs_repairer_id ON repairs(repairer_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_repairs_status ON repairs(status) WHERE deleted_at IS NULL;

CREATE TRIGGER update_repairs_updated_at
    BEFORE UPDATE ON repairs
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
