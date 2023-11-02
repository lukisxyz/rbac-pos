CREATE TABLE IF NOT EXISTS account_roles (
    account_id bytea,
    role_id bytea,
    PRIMARY KEY (account_id, role_id),
    FOREIGN KEY (account_id) REFERENCES accounts(id),
    FOREIGN KEY (role_id) REFERENCES roles(id)
);