CREATE TABLE IF NOT EXISTS trucks (
    truck_id VARCHAR(255) PRIMARY KEY,
    available_resources JSONB NOT NULL,
    travel_time_to_area JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_truck_id UNIQUE (truck_id)
); 