CREATE TABLE users (
    email VARCHAR(255) PRIMARY KEY,
    password BLOB NOT NULL,
    failed_attempts INT NOT NULL DEFAULT 0
);