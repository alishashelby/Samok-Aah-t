-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS models (
    model_id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    name VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS social_media (
    social_media_id SERIAL PRIMARY KEY,
    model_id INT NOT NULL REFERENCES models(model_id) ON DELETE CASCADE,
    platform VARCHAR(50) NOT NULL,
    url VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS portfolio_data (
    portfolio_id SERIAL PRIMARY KEY,
    model_id INT NOT NULL REFERENCES models(model_id) ON DELETE CASCADE,
    media_url VARCHAR(255) NOT NULL,
    description VARCHAR(500),
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_verified BOOLEAN DEFAULT FALSE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS portfolio_data;
DROP TABLE IF EXISTS social_media;
DROP TABLE IF EXISTS models;
-- +goose StatementEnd
