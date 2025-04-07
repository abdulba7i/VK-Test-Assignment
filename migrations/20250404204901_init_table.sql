-- +goose Up
CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role_id INT NOT NULL REFERENCES roles(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS actors (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    gender VARCHAR(10) NOT NULL CHECK (gender IN ('male', 'female')),
    date_of_birth DATE NOT NULL
);

CREATE TABLE IF NOT EXISTS films (
    id SERIAL PRIMARY KEY,
    name VARCHAR(150) NOT NULL UNIQUE,
    description VARCHAR(1000) NOT NULL,
    release_date DATE NOT NULL,
    rating NUMERIC(2, 1) CHECK (rating >= 0 AND rating <= 10) NOT NULL
);

CREATE TABLE IF NOT EXISTS actor_film (
    film_id INT REFERENCES films(id) ON DELETE CASCADE,
    actor_id INT REFERENCES actors(id) ON DELETE CASCADE,
    PRIMARY KEY (film_id, actor_id)
);

-- Индексы для быстрого поиска
CREATE INDEX IF NOT EXISTS idx_film_name ON films (name);
CREATE INDEX IF NOT EXISTS idx_actor_name ON actors (name);
CREATE INDEX IF NOT EXISTS idx_film_rating ON films (rating);
CREATE INDEX IF NOT EXISTS idx_film_release_date ON films (release_date);

-- Добавляем роли
INSERT INTO roles (name) VALUES ('user'), ('admin') 
ON CONFLICT (name) DO NOTHING;

-- +goose Down
DROP TABLE IF EXISTS actor_film;
DROP TABLE IF EXISTS films;
DROP TABLE IF EXISTS actors;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS roles;
