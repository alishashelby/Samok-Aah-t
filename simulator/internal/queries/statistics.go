package queries

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

type StatisticsQueries struct {
	db *sql.DB
}

func NewStatisticsQueries(db *sql.DB) *StatisticsQueries {
	return &StatisticsQueries{
		db: db,
	}
}

func (mq *StatisticsQueries) GetAdminStatisticsQuery() error {
	query := `
		SELECT
			DATE_TRUNC('month', o.created_at) AS month,
			COUNT(DISTINCT o.order_id) AS total_orders,
			COUNT(DISTINCT CASE WHEN o.status = 'COMPLETED' THEN o.order_id END) AS completed_orders,
			SUM(o.total_cost) AS total_revenue,
			AVG(o.total_cost) AS avg_order_value,
			COUNT(DISTINCT CASE WHEN u.gender = 'MALE' THEN u.user_id END) AS male_users,
			COUNT(DISTINCT CASE WHEN u.gender = 'FEMALE' THEN u.user_id END) AS female_users,
			COUNT(DISTINCT CASE WHEN bn.ban_id IS NOT NULL THEN u.user_id END) AS banned_users,
			SUM(CASE WHEN t.type = 'REFERRAL' THEN t.amount ELSE 0 END) AS referral_payouts,
			SUM(CASE WHEN t.type = 'CASHBACK' THEN t.amount ELSE 0 END) AS cashback_payouts
		FROM orders o
		JOIN booking bk ON o.booking_id = bk.booking_id
		LEFT JOIN users u ON (
			EXISTS (SELECT 1 FROM models m WHERE m.user_id = u.user_id) OR 
			EXISTS (SELECT 1 FROM clients cl WHERE cl.user_id = u.user_id))
		LEFT JOIN transactions t ON t.order_id = o.order_id AND t.type IN ('REFERRAL', 'CASHBACK')
		LEFT JOIN ban bn ON bn.user_id = u.user_id
		GROUP BY month
		ORDER BY month DESC
		LIMIT 40;
	`

	rows, err := mq.db.Query(query)
	if err != nil {
		return errors.New(fmt.Sprintf("error in statistics.go: GetAdminStatisticsQuery: %s", err))
	}
	defer rows.Close()

	log.Println("query GetAdminStatisticsQuery finished job")
	return nil
}

func (mq *StatisticsQueries) GetClientLoyaltyStatisticsQuery() error {
	query := `
		SELECT
			cl.client_id,
			u.gender,
			EXTRACT(YEAR FROM AGE(CURRENT_DATE, u.birth_date)) AS age,
			ct.name AS city,
			ll.name AS loyalty_level,
			COUNT(DISTINCT o.order_id) AS total_orders,
			SUM(o.total_cost) AS total_spent,
			AVG(o.total_cost) AS avg_order_value,
			COALESCE(SUM(t.amount), 0) AS total_cashback,
			MAX(o.created_at) AS last_order_date,
			COUNT(DISTINCT r.review_id) AS reviews_given,
			AVG(r.rating) AS avg_rating_given
		FROM clients cl
		JOIN users u ON cl.user_id = u.user_id
		JOIN cities ct ON u.city_id = ct.city_id
		JOIN loyalty_levels ll ON cl.loyalty_level_id = ll.level_id
		JOIN booking b ON cl.client_id = b.client_id
		JOIN orders o ON b.booking_id = o.booking_id
		LEFT JOIN transactions t ON t.order_id = o.order_id AND t.type = 'CASHBACK'
		LEFT JOIN reviews r ON o.order_id = r.order_id AND r.from_user_id = u.user_id
		WHERE o.status = 'COMPLETED'
		GROUP BY cl.client_id, u.gender, u.birth_date, ct.name, ll.name
		HAVING COUNT(DISTINCT o.order_id) >= 3
		ORDER BY total_spent DESC
		LIMIT 40;
	`

	rows, err := mq.db.Query(query)
	if err != nil {
		return errors.New(fmt.Sprintf("error in statistics.go: GetClientLoyaltyStatisticsQuery: %s", err))
	}
	defer rows.Close()

	log.Println("query GetClientLoyaltyStatisticsQuery finished job")
	return nil
}

func (mq *StatisticsQueries) GetReferralStatisticsQuery() error {
	start := time.Now()
	defer func() {
		_ = time.Since(start)
	}()

	query := `
		WITH referral_data AS (
            SELECT
                u.user_id,
                u.referral_code,
                COUNT(invited.user_id) AS invited_count,
                SUM(CASE WHEN invited.is_banned = FALSE THEN 1 ELSE 0 END) AS active_invited,
                SUM(
                    CASE 
                        WHEN invited.referral_count <= 3 THEN 500.0
                        WHEN invited.referral_count <= 7 THEN 300.0
                        ELSE 0
                    END
                ) AS total_bonuses
            FROM users u
            LEFT JOIN (
                SELECT 
                    invited.*,
                    ROW_NUMBER() OVER (PARTITION BY invited.referral_user_id ORDER BY invited.user_id) AS referral_count
                FROM users invited
                WHERE invited.referral_user_id IS NOT NULL
            ) invited ON invited.referral_user_id = u.user_id
            GROUP BY u.user_id, u.referral_code
        ),
        user_activity AS (
			SELECT
				u.user_id,
				CASE 
					WHEN m.model_id IS NOT NULL THEN 'model'
					WHEN c.client_id IS NOT NULL THEN 'client'
					ELSE 'unknown'
				END AS user_type,
				COUNT(o.order_id) AS completed_orders,
				SUM(o.total_cost) AS total_spent
			FROM users u
			LEFT JOIN models m ON m.user_id = u.user_id
			LEFT JOIN clients c ON c.user_id = u.user_id
			LEFT JOIN booking b ON b.client_id = c.client_id
			LEFT JOIN orders o ON o.booking_id = b.booking_id AND o.status = 'COMPLETED'
			GROUP BY u.user_id, user_type
		)
        SELECT
            rd.user_id,
            rd.referral_code,
            rd.invited_count,
            LEAST(7, rd.invited_count) AS effective_invites,
            rd.active_invited,
            rd.total_bonuses,
            ua.user_type,
            ua.completed_orders,
            ua.total_spent
        FROM referral_data rd
        JOIN user_activity ua ON rd.user_id = ua.user_id
        WHERE rd.invited_count > 0
        ORDER BY rd.total_bonuses DESC
		LIMIT 40;
	`

	rows, err := mq.db.Query(query)
	if err != nil {
		return errors.New(fmt.Sprintf("error in statistics.go: GetReferralStatisticsQuery: %s", err))
	}
	defer rows.Close()

	log.Println("query GetReferralStatisticsQuery finished job")
	return nil
}
