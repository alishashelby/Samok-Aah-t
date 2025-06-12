-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_cities_name ON cities(name);
CREATE INDEX idx_users_gender ON users(gender);
CREATE INDEX idx_categories_name ON categories(name);
CREATE INDEX idx_model_services_price ON model_services(price);

CREATE INDEX idx_users_city_id ON users(city_id);
CREATE INDEX idx_services_category_id ON services(category_id);
CREATE INDEX idx_model_services_service_id ON model_services(service_id);
CREATE INDEX idx_model_services_model_id ON model_services(model_id);

CREATE INDEX idx_portfolio_data_model_id ON portfolio_data(model_id);

CREATE INDEX idx_ban_user_id ON ban(user_id);

CREATE INDEX idx_reviews_to_user_id ON reviews(to_user_id);

CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_reviews_created_at ON reviews(created_at);

CREATE INDEX idx_booking_model_service_id ON booking(model_service_id);
CREATE INDEX idx_orders_booking_id ON orders(booking_id);
CREATE INDEX idx_reviews_order_id ON reviews(order_id);

CREATE INDEX idx_users_city_id_gender ON users(city_id, gender);

CREATE INDEX idx_model_services_service_id_price ON model_services(service_id, price);

CREATE INDEX idx_portfolio_model_id_verified ON portfolio_data(model_id) WHERE is_verified = TRUE;

CREATE INDEX idx_booking_booking_id ON booking(booking_id);
CREATE INDEX idx_model_services_pricing ON model_services(model_service_id, price);
CREATE INDEX idx_clients_loyalty ON clients(client_id, loyalty_level_id);

CREATE INDEX idx_orders_status_created ON orders(status, created_at);
CREATE INDEX idx_orders_created_desc ON orders(created_at DESC);
CREATE INDEX idx_transactions_order_status ON transactions(order_id, status);
CREATE INDEX idx_reviews_order_rating ON reviews(order_id, rating);

CREATE INDEX idx_transactions_order_id ON transactions(order_id);
CREATE INDEX idx_clients_loyalty_level ON clients(loyalty_level_id);

CREATE INDEX idx_models_user_id ON models(user_id);
CREATE INDEX idx_clients_user_id ON clients(user_id);

CREATE INDEX idx_orders_created_month ON orders(DATE_TRUNC('month', created_at));
CREATE INDEX idx_transactions_type_order ON transactions(type, order_id);

CREATE INDEX idx_clients_user_id_client ON clients(user_id, client_id);
CREATE INDEX idx_reviews_from_order ON reviews(from_user_id, order_id, rating);
CREATE INDEX idx_booking_client_id_created ON booking(client_id);
CREATE INDEX idx_users_birth_date ON users(birth_date);

CREATE INDEX idx_users_referral_user_id ON users(referral_user_id);

CREATE INDEX idx_transactions_order_type ON transactions(order_id, type);
CREATE INDEX idx_model_services_model ON model_services(model_id, model_service_id);

CREATE INDEX idx_transactions_created_desc ON transactions(created_at DESC);
CREATE INDEX idx_external_transactions_payment ON external_transactions(payment_system_id, external_transaction_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_cities_name;
DROP INDEX IF EXISTS idx_users_gender;
DROP INDEX IF EXISTS idx_categories_name;
DROP INDEX IF EXISTS idx_model_services_price;

DROP INDEX IF EXISTS idx_users_city_id;
DROP INDEX IF EXISTS idx_services_category_id;
DROP INDEX IF EXISTS idx_model_services_service_id;
DROP INDEX IF EXISTS idx_model_services_model_id;

DROP INDEX IF EXISTS idx_portfolio_data_model_id;

DROP INDEX IF EXISTS idx_ban_user_id;

DROP INDEX IF EXISTS idx_reviews_to_user_id;

DROP INDEX IF EXISTS idx_orders_status;
DROP INDEX IF EXISTS idx_reviews_created_at;

DROP INDEX IF EXISTS idx_booking_model_service_id;
DROP INDEX IF EXISTS idx_orders_booking_id;
DROP INDEX IF EXISTS idx_reviews_order_id;

DROP INDEX IF EXISTS idx_users_city_id_gender;

DROP INDEX IF EXISTS idx_model_services_service_id_price;

DROP INDEX IF EXISTS idx_portfolio_model_id_verified;

DROP INDEX IF EXISTS idx_booking_booking_id;
DROP INDEX IF EXISTS idx_model_services_pricing;
DROP INDEX IF EXISTS idx_clients_loyalty;

DROP INDEX IF EXISTS idx_orders_status_created;
DROP INDEX IF EXISTS idx_orders_created_desc;
DROP INDEX IF EXISTS idx_transactions_order_status;
DROP INDEX IF EXISTS idx_reviews_order_rating;

DROP INDEX IF EXISTS idx_transactions_order_id;
DROP INDEX IF EXISTS idx_clients_loyalty_level;

DROP INDEX IF EXISTS idx_models_user_id;
DROP INDEX IF EXISTS idx_clients_user_id;

DROP INDEX IF EXISTS idx_orders_created_month;
DROP INDEX IF EXISTS idx_transactions_type_order;

DROP INDEX IF EXISTS idx_clients_user_id_client;
DROP INDEX IF EXISTS idx_reviews_from_order;
DROP INDEX IF EXISTS idx_booking_client_id_created;
DROP INDEX IF EXISTS idx_users_birth_date;

DROP INDEX IF EXISTS idx_users_referral_user_id;

DROP INDEX IF EXISTS idx_transactions_order_type;
DROP INDEX IF EXISTS idx_model_services_model;

DROP INDEX IF EXISTS idx_transactions_created_desc;
DROP INDEX IF EXISTS idx_external_transactions_payment;
-- +goose StatementEnd
