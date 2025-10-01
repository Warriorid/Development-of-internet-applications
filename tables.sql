CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(150) NOT NULL UNIQUE,
    password VARCHAR(128) NOT NULL,
    is_staff BOOLEAN DEFAULT FALSE NOT NULL
);


CREATE TABLE materials (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL UNIQUE,
    coefficient FLOAT NOT NULL CHECK (coefficient > 0),
    image_url VARCHAR(255) NULL,
    description TEXT NOT NULL,
    is_deleted BOOLEAN DEFAULT FALSE NOT NULL
);

CREATE TABLE pits_calculations (
    id SERIAL PRIMARY KEY,
    status VARCHAR(20) NOT NULL CHECK (status IN ('draft', 'deleted', 'formed', 'completed', 'rejected')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    formed_at TIMESTAMP NULL,
    completed_at TIMESTAMP NULL,
    creator_id INTEGER NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    moderator_id INTEGER NULL REFERENCES users(id) ON DELETE RESTRICT,
    pit_length FLOAT NULL CHECK (pit_length > 0),
    pit_width FLOAT NULL CHECK (pit_width > 0),
    pit_depth FLOAT NULL CHECK (pit_depth > 0),
    CHECK (formed_at IS NULL OR formed_at >= created_at),
    CHECK (completed_at IS NULL OR completed_at >= formed_at)
);


CREATE UNIQUE INDEX one_draft_per_user 
ON pits_calculations (creator_id) 
WHERE status = 'draft';

CREATE TABLE calculation_materials (
    calculation_id INTEGER NOT NULL REFERENCES pits_calculations(id) ON DELETE RESTRICT,
    material_id INTEGER NOT NULL REFERENCES materials(id) ON DELETE RESTRICT,
    slope_angle INTEGER NULL CHECK (slope_angle >= 0 AND slope_angle <= 90) DEFAULT 90,
    volume_result FLOAT NULL,
    PRIMARY KEY (calculation_id, material_id)
);