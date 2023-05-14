-- +goose Up
CREATE TABLE "user"
(
    id         bigserial                   PRIMARY KEY,
    username   text                        NOT NULL UNIQUE,
    email      text                        NOT NULL UNIQUE,
    password   text                        NOT NULL,
    role       text                        NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);


-- +goose Down
DROP TABLE IF EXISTS "user";
