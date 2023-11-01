CREATE TABLE IF NOT EXISTS role_permissions (
    permission_id bytea,
    role_id bytea,
    PRIMARY KEY (permission_id, role_id),
    FOREIGN KEY (permission_id) REFERENCES permissions(id),
    FOREIGN KEY (role_id) REFERENCES roles(id)
);