CREATE TABLE IF NOT EXISTS requests (
    id SERIAL PRIMARY KEY,

    title VARCHAR(255) NOT NULL,

    description TEXT,

    status VARCHAR(50) NOT NULL DEFAULT 'new',

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);