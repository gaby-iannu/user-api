CREATE TABLE IF NOT EXISTS failed_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL,
    event_type VARCHAR(50) NOT NULL,
    user_id UUID NOT NULL,
    payload JSONB NOT NULL,
    error TEXT NOT NULL,
    attempts INTEGER NOT NULL DEFAULT 3,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_error TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_failed_events_user_id ON failed_events(user_id);
CREATE INDEX idx_failed_events_event_type ON failed_events(event_type);
CREATE INDEX idx_failed_events_created_at ON failed_events(created_at);
