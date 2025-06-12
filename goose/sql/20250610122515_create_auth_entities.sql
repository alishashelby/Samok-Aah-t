-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS auth (
    auth_id SERIAL PRIMARY KEY,
    email varchar(320) NOT NULL UNIQUE,
    phone varchar(15) NOT NULL UNIQUE,
    password_hash varchar(255) NOT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS admins (
    admin_id SERIAL PRIMARY KEY,
    auth_id INT NOT NULL REFERENCES auth(auth_id) ON DELETE CASCADE,
    permissions json NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
    user_id SERIAL PRIMARY KEY,
    auth_id INT NOT NULL REFERENCES auth(auth_id) ON DELETE CASCADE,
    birth_date date CHECK (age(CURRENT_DATE, birth_date) >= '18 years'),
    gender varchar(10),
    city_id INT NOT NULL REFERENCES cities(city_id) ON DELETE SET NULL,
    passport_series varchar(4),
    passport_number varchar(6),
    passport_issue_date DATE,
    passport_verified BOOLEAN DEFAULT FALSE,
    referral_code uuid UNIQUE,
    referral_user_id INT REFERENCES users(user_id) ON DELETE SET NULL,
    referral_user_count SMALLINT DEFAULT 0 CHECK (referral_user_count <= 7),
    is_banned BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS ban (
    ban_id SERIAL PRIMARY KEY,
    admin_id INT NOT NULL REFERENCES admins(admin_id) ON DELETE SET NULL,
    user_id INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    reason varchar(500) NOT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE FUNCTION process_ban_status() RETURNS TRIGGER AS $ban_status$
BEGIN
    UPDATE users SET is_banned = TRUE WHERE user_id = NEW.user_id;
    RETURN NEW;
END;
$ban_status$ LANGUAGE plpgsql;

CREATE TRIGGER ban_status
    AFTER INSERT ON ban
    FOR EACH ROW EXECUTE PROCEDURE process_ban_status();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS ban;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS admins;
DROP TABLE IF EXISTS auth;
DROP FUNCTION IF EXISTS process_ban_status;
DROP TRIGGER IF EXISTS ban_status ON ban;
-- +goose StatementEnd
