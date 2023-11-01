CREATE TABLE IF NOT EXISTS roles (
    id bytea PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL
);
