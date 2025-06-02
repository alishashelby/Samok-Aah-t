package queries

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type OrderQueries struct {
	db *sql.DB
}

func NewOrderQueries(db *sql.DB) *OrderQueries {
	return &OrderQueries{
		db: db,
	}
}

func (mq *OrderQueries) GetAllDataToCalculateTotalCostQuery() error {
	query := `
		SELECT
			b.booking_id,
			ms.price AS base_price,
			a_s.offer_price AS additional_service_price,
			p.percentage AS discount_percent,
			ll.cashback_percentage,
			ROUND((
				(ms.price + COALESCE(a_s.offer_price, 0)) * 
				(1 - COALESCE(p.percentage, 0) / 100.0)
			) * 1.1, 2) AS total_cost,
			ROUND((
				(ms.price + COALESCE(a_s.offer_price, 0)) * 
				(1 - COALESCE(p.percentage, 0) / 100.0)
			) * (ll.cashback_percentage / 100.0), 2) AS cashback_amount
		FROM booking b
		JOIN orders o ON o.booking_id = b.booking_id
		JOIN model_services ms ON b.model_service_id = ms.model_service_id
		LEFT JOIN additional_services a_s ON b.additional_service_id = a_s.additional_service_id
		LEFT JOIN promocodes p ON o.promocode_id = p.promocode_id
		JOIN clients cl ON b.client_id = cl.client_id
		JOIN loyalty_levels ll ON cl.loyalty_level_id = ll.level_id
		WHERE b.booking_id = 1234;
	`

	rows, err := mq.db.Query(query)
	if err != nil {
		return errors.New(fmt.Sprintf("error in order.go: GetAllDataToCalculateTotalCostQuery: %s", err))
	}
	defer rows.Close()

	log.Println("query GetAllDataToCalculateTotalCostQuery finished job")
	return nil
}

func (mq *OrderQueries) GetComplicatedOrdersQuery() error {
	query := `
		SELECT
			o.order_id,
			o.status,
			o.total_cost,
			o.created_at,
			cl.name AS client_name,
			m.name AS model_name,
			b.date_time,
			b.address->>'street' AS street,
			b.address->>'house' AS house,
			COUNT(r.review_id) AS review_count,
			AVG(r.rating) AS avg_rating,
			STRING_AGG(DISTINCT t.type, ', ') AS transaction_types
		FROM orders o
		JOIN booking b ON o.booking_id = b.booking_id
		JOIN clients cl ON b.client_id = cl.client_id
		JOIN model_services ms ON b.model_service_id = ms.model_service_id
		JOIN models m ON ms.model_id = m.model_id
		LEFT JOIN reviews r ON o.order_id = r.order_id
		LEFT JOIN transactions t ON o.order_id = t.order_id
		WHERE 
			o.status IN ('CANCELLED', 'REJECTED')
			OR EXISTS (
				SELECT 1 
				FROM transactions t2 
				WHERE t2.order_id = o.order_id 
				AND t2.status = 'FAILURE'
			)
			OR (o.created_at < NOW() - INTERVAL '1 hour' AND o.status = 'IN_PROCESS')
		GROUP BY o.order_id, o.status, o.total_cost, o.created_at, cl.name, m.name, b.date_time, street, house
		HAVING AVG(r.rating) < 3.0 OR COUNT(r.review_id) = 0
		ORDER BY o.created_at DESC
		LIMIT 40;
	`

	rows, err := mq.db.Query(query)
	if err != nil {
		return errors.New(fmt.Sprintf("error in order.go: GetComplicatedOrdersQuery: %s", err))
	}
	defer rows.Close()

	log.Println("query GetComplicatedOrdersQuery finished job")
	return nil
}
