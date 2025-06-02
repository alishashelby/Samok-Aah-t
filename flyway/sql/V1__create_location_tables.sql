CREATE TABLE IF NOT EXISTS regions (
    region_id SERIAL PRIMARY KEY,
    name varchar(40) NOT NULL
);

CREATE TABLE IF NOT EXISTS cities (
    city_id SERIAL PRIMARY KEY,
    name varchar(20) NOT NULL,
    region_id INT NOT NULL REFERENCES regions(region_id) ON DELETE CASCADE
);