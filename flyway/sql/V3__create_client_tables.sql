CREATE TABLE IF NOT EXISTS loyalty_levels (
    level_id SERIAL PRIMARY KEY,
    name VARCHAR(30) NOT NULL UNIQUE,
    min_orders INT NOT NULL CHECK (min_orders >= 0),
    cashback_percentage INT NOT NULL CHECK (cashback_percentage BETWEEN 0 AND 100)
);

CREATE TABLE IF NOT EXISTS clients (
    client_id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    name VARCHAR(50) NOT NULL,
    loyalty_level_id INT NOT NULL REFERENCES loyalty_levels(level_id) ON DELETE SET NULL
);