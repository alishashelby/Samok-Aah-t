CREATE TABLE IF NOT EXISTS promocodes (
    promocode_id SERIAL PRIMARY KEY,
    code VARCHAR(30) NOT NULL UNIQUE,
    percentage SMALLINT NOT NULL CHECK (percentage BETWEEN 1 AND 100),
    start_time TIMESTAMP,
    finish_time TIMESTAMP,
    is_always BOOLEAN DEFAULT FALSE,
    CHECK (finish_time > start_time OR is_always)
);

CREATE TABLE IF NOT EXISTS orders (
    order_id SERIAL PRIMARY KEY,
    booking_id INT NOT NULL UNIQUE REFERENCES booking(booking_id) ON DELETE CASCADE,
    promocode_id INT REFERENCES promocodes(promocode_id) ON DELETE SET NULL,
    platform_fee DECIMAL(9,2) NOT NULL CHECK (platform_fee >= 0),
    total_cost DECIMAL(9,2) NOT NULL CHECK (total_cost >= 0),
    status VARCHAR(20) NOT NULL CHECK (
        status IN ('IN_PROCESS', 'CONFIRMED', 'IN_TRANSIT', 'COMPLETED', 'CANCELLED')
        ),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    security_code VARCHAR(6) NOT NULL
);

CREATE TABLE IF NOT EXISTS reviews (
    review_id SERIAL PRIMARY KEY,
    order_id INT NOT NULL REFERENCES orders(order_id) ON DELETE CASCADE,
    from_user_id INT NOT NULL REFERENCES users(user_id),
    to_user_id INT NOT NULL REFERENCES users(user_id),
    rating SMALLINT NOT NULL CHECK (rating BETWEEN 1 AND 5),
    description VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);