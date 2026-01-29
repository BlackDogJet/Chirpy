-- +goose Up
CREATE TABLE users(
    id UUID PRIMARY KEY,
    created_at TIMESTAMP not null,
    updated_at TIMESTAMP not null,
    email VARCHAR(100) not null UNIQUE
);

-- +goose Down
DROP TABLE users;