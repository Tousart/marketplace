-- +migrate up
CREATE TABLE adverts (
    advert_id VARCHAR(64) PRIMARY KEY,
    user_id VARCHAR(64) NOT NULL,
    title VARCHAR(128) NOT NULL,
    text TEXT NOT NULL,
    url TEXT NOT NULL,
    price BIGINT NOT NULL,
    date TIMESTAMP NOT NULL
);