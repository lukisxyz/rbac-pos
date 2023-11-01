CREATE TABLE IF NOT EXISTS permissions (
    id bytea PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    url VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL
);
