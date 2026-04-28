-- +migrate Up
CREATE TABLE feedbacks (
    id SERIAL PRIMARY KEY,
    data TEXT NOT NULL
);