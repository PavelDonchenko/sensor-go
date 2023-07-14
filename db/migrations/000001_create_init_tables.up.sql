CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE sensor_group (
    id int NOT NULL PRIMARY KEY,
    name text NOT NULL UNIQUE,
    sensors uuid [],
    created_at timestamp NOT NULL DEFAULT NOW()
);

CREATE TABLE sensor (
    id uuid  PRIMARY KEY DEFAULT uuid_generate_v4(),
    group_id int NOT NULL REFERENCES sensor_group (id) ON DELETE CASCADE,
    group_name text NOT NULL REFERENCES sensor_group (name) ON DELETE CASCADE,
    in_group_id int NOT NULL,
    data_output_rate int NOT NULL,
    temperature double precision,
    transparency int CHECK ( transparency between 0 and 100),
    x double precision NOT NULL,
    y double precision NOT NULL,
    z double precision NOT NULL,
    fishes uuid [],
    updated_at timestamp,
    created_at timestamp NOT NULL DEFAULT NOW()
);


CREATE TABLE detected_fish (
    id uuid  PRIMARY KEY DEFAULT uuid_generate_v4(),
    name text NOT NULL,
    count int NOT NULL,
    sensorID uuid REFERENCES sensor (id) ON DELETE CASCADE,
    created_at timestamp NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_fish_name ON detected_fish(name);

