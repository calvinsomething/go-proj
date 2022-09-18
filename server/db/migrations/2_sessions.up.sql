CREATE TABLE sessions (
    id CHAR(36) NOT NULL PRIMARY KEY,
    data BLOB NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);