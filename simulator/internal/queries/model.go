package queries

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

type ModelQueries struct {
	db *sql.DB
}

func NewModelQueries(db *sql.DB) *ModelQueries {
	return &ModelQueries{
		db: db,
	}
}

func (mq *ModelQueries) FilterModelsV1Query() error {
	start := time.Now()
	defer func() {
		_ = time.Since(start)
	}()

	query := `
		SELECT
			m.model_id,
			m.name AS model_name,
			u.gender,
			ct.name AS city,
			r.name AS region,
			AVG(rv.rating) AS avg_rating,
			COUNT(rv.review_id) AS review_count,
			json_agg(DISTINCT cat.name) AS categories,
			MIN(ms.price) AS min_price
		FROM models m
		JOIN users u ON m.user_id = u.user_id
		JOIN cities ct ON u.city_id = ct.city_id
		JOIN regions r ON ct.region_id = r.region_id
		JOIN model_services ms ON m.model_id = ms.model_id
		JOIN services s ON ms.service_id = s.service_id
		JOIN categories cat ON s.category_id = cat.category_id
		LEFT JOIN reviews rv ON rv.to_user_id = u.user_id
		LEFT JOIN portfolio_data p ON p.model_id = m.model_id AND p.is_verified = TRUE
		WHERE 
			ct.name = 'Санкт-Петербург' 
			AND u.gender = 'FEMALE'
			AND cat.name = 'Полицейский(ая)'
			AND ms.price BETWEEN 5000 AND 20000
			AND EXISTS (
				SELECT 1 
				FROM portfolio_data pd 
				WHERE pd.model_id = m.model_id
			)
			AND NOT EXISTS (
				SELECT 1 
				FROM ban b 
				WHERE b.user_id = u.user_id
			)
		GROUP BY m.model_id, m.name, u.gender, ct.name, r.name
		HAVING AVG(rv.rating) >= 4.0
		ORDER BY avg_rating DESC, min_price
		LIMIT 40;
	`

	rows, err := mq.db.Query(query)
	if err != nil {
		return errors.New(fmt.Sprintf("error in model.go: FilterModelsV1Query: %s", err))
	}
	defer rows.Close()

	log.Println("query FilterModelsV1Query finished job")
	return nil
}

func (mq *ModelQueries) FilterModelsV2Query() error {
	query := `
		SELECT
			cat.name AS category,
			m.model_id,
			m.name AS model_name,
			AVG(r.rating) AS avg_rating,
			COUNT(r.review_id) AS review_count,
			COUNT(DISTINCT o.order_id) AS completed_orders,
			SUM(o.total_cost) AS total_income,
			PERCENT_RANK() OVER (PARTITION BY cat.category_id ORDER BY AVG(r.rating) DESC) AS rating_percentile
		FROM categories cat
		JOIN services s ON s.category_id = cat.category_id
		JOIN model_services ms ON ms.service_id = s.service_id
		JOIN models m ON ms.model_id = m.model_id
		JOIN booking b ON ms.model_service_id = b.model_service_id
		JOIN orders o ON b.booking_id = o.booking_id AND o.status = 'COMPLETED'
		JOIN reviews r ON o.order_id = r.order_id
		WHERE r.created_at > CURRENT_DATE - INTERVAL '6 months'
		GROUP BY cat.category_id, cat.name, m.model_id, m.name
		HAVING COUNT(r.review_id) >= 5
		ORDER BY cat.name, rating_percentile DESC
		LIMIT 40;
	`

	rows, err := mq.db.Query(query)
	if err != nil {
		return errors.New(fmt.Sprintf("error in model.go: FilterModelsV2Query: %s", err))
	}
	defer rows.Close()

	log.Println("query FilterModelsV2Query finished job")
	return nil
}
