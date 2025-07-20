-- +migrate Up
CREATE TABLE users (
    user_id VARCHAR(64) PRIMARY KEY,
    login VARCHAR(32) UNIQUE NOT NULL,
    password VARCHAR(128) NOT NULL
);
