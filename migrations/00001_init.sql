-- +goose Up
CREATE TABLE users (
    id serial primary key,
    email text,
    password  text,
    createdAt timestamp,
    updatedAt timestamp
);

-- +goose Down
DROP TABLE users;
