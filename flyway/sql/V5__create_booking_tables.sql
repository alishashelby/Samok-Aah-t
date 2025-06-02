CREATE TABLE IF NOT EXISTS categories (
    category_id SERIAL PRIMARY KEY,
    name VARCHAR(70) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS services (
    service_id SERIAL PRIMARY KEY,
    category_id INT NOT NULL REFERENCES categories(category_id),
    description VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS model_services (
    model_service_id SERIAL PRIMARY KEY,
    model_id INT NOT NULL REFERENCES models(model_id) ON DELETE CASCADE,
    service_id INT NOT NULL REFERENCES services(service_id),
    price DECIMAL(9,2) NOT NULL CHECK (price >= 0)
);

CREATE TABLE IF NOT EXISTS additional_services (
    additional_service_id SERIAL PRIMARY KEY,
    description VARCHAR(1000) NOT NULL,
    offer_price DECIMAL(9,2) NOT NULL CHECK (offer_price >= 0),
    status VARCHAR(20) NOT NULL CHECK (
        status IN ('PENDING', 'REJECTED', 'APPROVED', 'CANCELLED', 'HIGHER_PRICE')
        ),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS booking (
    booking_id SERIAL PRIMARY KEY,
    client_id INT NOT NULL REFERENCES clients(client_id) ON DELETE CASCADE,
    model_service_id INT NOT NULL REFERENCES model_services(model_service_id) ON DELETE SET NULL,
    date_time TIMESTAMP NOT NULL,
    duration INTERVAL NOT NULL,
    address json NOT NULL,
    additional_service_id INT REFERENCES additional_services(additional_service_id),
    status VARCHAR(20) NOT NULL CHECK (status IN ('PENDING', 'REJECTED', 'APPROVED', 'CANCELLED')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
