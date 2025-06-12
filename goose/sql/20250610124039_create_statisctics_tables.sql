-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS daily_statistics (
    stat_id SERIAL PRIMARY KEY,
    date DATE NOT NULL UNIQUE DEFAULT CURRENT_DATE,
    total_clients INT NOT NULL CHECK (total_clients >= 0),
    total_models INT NOT NULL CHECK (total_models >= 0),
    total_orders INT NOT NULL CHECK (total_orders >= 0),
    completed_orders INT NOT NULL CHECK (completed_orders >= 0 AND completed_orders <= total_orders),
    avg_order_cost DECIMAL(9,2) CHECK (avg_order_cost >= 0),
    total_referrals INT CHECK (total_referrals >= 0),
    referral_bonuses DECIMAL(9,2) CHECK (referral_bonuses >= 0)
);

CREATE OR REPLACE PROCEDURE update_statistics()
    LANGUAGE plpgsql AS $$
BEGIN
    INSERT INTO daily_statistics (
        date,
        total_clients,
        total_models,
        total_orders,
        completed_orders,
        avg_order_cost,
        total_referrals,
        referral_bonuses
    )
    SELECT
        CURRENT_DATE,
        (SELECT COUNT(*) FROM clients),
        (SELECT COUNT(*) FROM models),
        (SELECT COUNT(*) FROM orders),
        (SELECT COUNT(*) FROM orders WHERE status = 'COMPLETED'),
        (SELECT AVG(total_cost) FROM orders),
        (SELECT COUNT(*) FROM transactions WHERE type = 'REFERRAL'),
        (SELECT SUM(amount) FROM transactions WHERE type = 'REFERRAL')
    ON CONFLICT (date) DO UPDATE SET
                                     total_clients = excluded.total_clients,
                                     total_models = excluded.total_models,
                                     total_orders = excluded.total_orders,
                                     completed_orders = excluded.completed_orders,
                                     avg_order_cost = excluded.avg_order_cost,
                                     total_referrals = excluded.total_referrals,
                                     referral_bonuses = excluded.referral_bonuses;
END;
$$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS daily_statistics;
DROP PROCEDURE update_statistics;
-- +goose StatementEnd
