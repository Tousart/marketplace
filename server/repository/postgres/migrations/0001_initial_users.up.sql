-- +migrate Up
CREATE TABLE users (
    user_id VARCHAR(64) PRIMARY KEY,
    login VARCHAR(32) UNIQUE,
    password VARCHAR(128)
);
