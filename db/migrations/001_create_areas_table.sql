CREATE TABLE IF NOT EXISTS areas (
    area_id VARCHAR(255) PRIMARY KEY,
    urgency_level INTEGER NOT NULL CHECK (urgency_level BETWEEN 1 AND 5),
    required_resources JSONB NOT NULL,
    time_constraint INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_area_id UNIQUE (area_id)
); 