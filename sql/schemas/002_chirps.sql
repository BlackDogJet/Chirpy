-- +goose Up
CREATE TABLE chirps(
    id UUID PRIMARY KEY,
    created_at TIMESTAMP not null,
    updated_at TIMESTAMP not null,
    body VARCHAR not null,
    user_id UUID not null
);

-- +goose Down
DROP TABLE chirps;