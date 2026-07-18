CREATE TABLE IF NOT EXISTS containers (
    id UUID PRIMARY KEY,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    capacity DOUBLE PRECISION NOT NULL,
    current_fill_level DOUBLE PRECISION DEFAULT 0.0,
    temperature DOUBLE PRECISION DEFAULT 25.0,
    gas_level_ppm DOUBLE PRECISION DEFAULT 0.0,
    battery_status DOUBLE PRECISION DEFAULT 100.0,
    last_emptied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Atık Olayları (Wastes) Tablosu
CREATE TABLE IF NOT EXISTS wastes (
    id UUID PRIMARY KEY,
    container_id UUID NOT NULL,
    type VARCHAR(50) NOT NULL,
    weight DOUBLE PRECISION NOT NULL,
    carbon_footprint DOUBLE PRECISION DEFAULT 0.0,
    fill_level_at_drop DOUBLE PRECISION NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_container
        FOREIGN KEY(container_id) 
        REFERENCES containers(id)
        ON DELETE CASCADE
);