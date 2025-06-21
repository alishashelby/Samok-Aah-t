package queries

import (
	"database/sql"
	"fmt"
	"log"
)

type TransactionQueries struct {
	db *sql.DB
}

func NewTransactionQueries(db *sql.DB) *TransactionQueries {
	return &TransactionQueries{
		db: db,
	}
}

func (mq *TransactionQueries) CalculateModelFinancialQuery() error {
	query := `
		SELECT
			m.model_id,
			m.name AS model_name,
			COUNT(DISTINCT o.order_id) AS completed_orders,
			SUM(o.total_cost) AS total_revenue,
			SUM(o.total_cost * 0.9) AS model_income,
			COALESCE(SUM(CASE WHEN t.type = 'REFERRAL' THEN t.amount ELSE 0 END), 0) AS referral_bonuses,
			AVG(rv.rating) AS avg_rating,
			COUNT(DISTINCT rv.review_id) AS review_count
		FROM models m
		JOIN model_services ms ON m.model_id = ms.model_id
		JOIN booking b ON ms.model_service_id = b.model_service_id
		JOIN orders o ON b.booking_id = o.booking_id
		JOIN transactions t ON t.order_id = o.order_id AND t.type = 'MODEL_PAYOUT'
		LEFT JOIN reviews rv ON rv.to_user_id = m.user_id
		WHERE 
			o.status = 'COMPLETED'
			AND o.created_at BETWEEN '2023-01-01' AND '2023-12-31'
		GROUP BY m.model_id, m.name
		HAVING COUNT(DISTINCT o.order_id) > 5
		ORDER BY model_income DESC
		LIMIT 40;
	`

	rows, err := mq.db.Query(query)
	if err != nil {
		return fmt.Errorf("error in transaction.go: CalculateModelFinancialQuery: %s", err)
	}
	defer rows.Close()

	log.Println("query CalculateModelFinancialQuery finished job")
	return nil
}

func (mq *TransactionQueries) GetAllTransactionsQuery() error {
	query := `
		SELECT
			t.transaction_id,
			t.type,
			t.amount,
			t.status,
			t.created_at,
			o.order_id,
			o.total_cost,
			COALESCE(mu.user_id, cu.user_id) AS user_id,
			CASE 
				WHEN m.model_id IS NOT NULL THEN 'model' 
				WHEN cl.client_id IS NOT NULL THEN 'client'
				ELSE 'admin'
			END AS user_type,
			et.failure_msg,
			ps.name AS payment_system
		FROM transactions t
		LEFT JOIN orders o ON t.order_id = o.order_id
		LEFT JOIN booking b ON o.booking_id = b.booking_id
		LEFT JOIN model_services ms ON b.model_service_id = ms.model_service_id
		LEFT JOIN models m ON ms.model_id = m.model_id
		LEFT JOIN users mu ON m.user_id = mu.user_id

		LEFT JOIN clients cl ON b.client_id = cl.client_id
		LEFT JOIN users cu ON cl.user_id = cu.user_id

		LEFT JOIN external_transactions et ON t.external_transaction_id = et.external_transaction_id
		LEFT JOIN payment_system_integration ps ON et.payment_system_id = ps.payment_system_id
		WHERE t.created_at BETWEEN '2023-01-01' AND '2023-12-31'
		ORDER BY t.created_at DESC
		LIMIT 40;
	`

	rows, err := mq.db.Query(query)
	if err != nil {
		return fmt.Errorf("error in transaction.go: GetAllTransactionsQuery: %s", err)
	}
	defer rows.Close()

	log.Println("query GetAllTransactionsQuery finished job")
	return nil
}
