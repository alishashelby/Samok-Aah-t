-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS payment_system_integration (
    payment_system_id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS external_transactions (
    external_transaction_id SERIAL PRIMARY KEY,
    payment_system_id INT NOT NULL REFERENCES payment_system_integration(payment_system_id),
    failure_msg VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS transactions (
    transaction_id SERIAL PRIMARY KEY,
    amount DECIMAL(9,2) NOT NULL CHECK (amount >= 0),
    type VARCHAR(20) NOT NULL CHECK (
        type IN (
                 'ORDER_PAYMENT', 'ORDER_INCOME', 'CLIENT_DEPOSIT', 'MODEL_PAYOUT',
                 'REFUND_TO_CLIENT', 'REFUND_FROM_MODEL', 'ORDER_CANCELLATION',
                 'REFERRAL', 'CASHBACK'
                )
        ),
    order_id INT REFERENCES orders(order_id) ON DELETE CASCADE,
    external_transaction_id INT REFERENCES external_transactions(external_transaction_id),
    reason VARCHAR(255),
    status VARCHAR(20) NOT NULL CHECK (status IN ('FAILURE', 'PENDING', 'SUCCESS')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    processed_at TIMESTAMP,
    CHECK (
        (type IN ('CLIENT_DEPOSIT', 'MODEL_PAYOUT', 'REFERRAL') AND order_id IS NULL)
            OR
        (type NOT IN ('CLIENT_DEPOSIT', 'MODEL_PAYOUT', 'REFERRAL') AND order_id IS NOT NULL)
        ),
    CHECK (
        (type IN ('CLIENT_DEPOSIT', 'MODEL_PAYOUT') AND external_transaction_id IS NOT NULL)
            OR
        (type NOT IN ('CLIENT_DEPOSIT', 'MODEL_PAYOUT') AND external_transaction_id IS NULL)
        ),
    CHECK (
        (type IN ('REFUND_TO_CLIENT', 'REFUND_FROM_MODEL', 'ORDER_CANCELLATION') AND reason IS NOT NULL )
            OR
        (type NOT IN ('REFUND_TO_CLIENT', 'REFUND_FROM_MODEL', 'ORDER_CANCELLATION') AND reason IS NULL)
        )
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS external_transactions;
DROP TABLE IF EXISTS payment_system_integration;
-- +goose StatementEnd
