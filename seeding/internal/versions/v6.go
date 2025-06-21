package versions

import (
	"database/sql"
	"github.com/alishashelby/Samok-Aah-t/seeding/internal/service"
	"log"
	"strings"
	"time"
)

type SeedingV6 struct {
	service *service.SeedingService
}

func NewSeedingV6(service *service.SeedingService) *SeedingV6 {
	return &SeedingV6{
		service: service,
	}
}

type Promocode struct {
	ID         int
	Code       string
	Percentage int
	StartTime  time.Time
	FinishTime time.Time
	IsAlways   bool
}

type Order struct {
	ID           int
	BookingID    int
	PromocodeID  *int
	PlatformFee  float64
	TotalCost    float64
	Status       string
	SecurityCode string
	CreatedAt    time.Time
}

type Review struct {
	ID          int
	OrderID     int
	FromUserID  int
	ToUserID    int
	Rating      int
	Description string
	CreatedAt   time.Time
}

// nolint:dupl
func (s *SeedingV6) FakeV6() error {
	if !s.service.CheckTable("promocodes") ||
		!s.service.CheckTable("orders") ||
		!s.service.CheckTable("reviews") {
		return nil
	}

	if !s.service.CheckIfFilled("promocodes") {
		log.Println("v6.go promocodes are not filled")
		if err := s.insertPromocodes(); err != nil {
			return err
		}
		log.Println("filled")
	}

	if !s.service.CheckIfFilled("orders") {
		log.Println("v6.go orders are not filled")
		if err := s.insertOrders(); err != nil {
			return err
		}
		log.Println("filled")
	}

	if !s.service.CheckIfFilled("reviews") {
		log.Println("v6.go reviews are not filled")
		if err := s.insertReviews(); err != nil {
			return err
		}
		log.Println("filled")
	}

	return nil
}

func (s *SeedingV6) insertPromocodes() error {
	tx, err := s.service.Db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		//nolint:errcheck
		tx.Rollback()
	}()

	stmt, err := tx.Prepare("INSERT INTO promocodes (" +
		"code, percentage, start_time, finish_time, is_always)" +
		"VALUES ($1, $2, $3, $4, $5)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	usedCodes := make(map[string]bool)
	for range s.service.GetPromocodeCount() {
		var code string
		for {
			code = strings.ToUpper(s.service.Fake.RandomStringWithLength(6))
			if !usedCodes[code] {
				usedCodes[code] = true
				break
			}
		}

		promocode := &Promocode{
			Code:       code,
			Percentage: s.service.Fake.IntBetween(3, 15),
			IsAlways:   s.service.Fake.Bool(),
		}

		if !promocode.IsAlways {
			startTime := s.service.Fake.Time().TimeBetween(
				time.Now().AddDate(0, -1, 0), time.Now().AddDate(0, 2, 0))
			promocode.StartTime = startTime

			finishTime := startTime.AddDate(0, 0, s.service.Fake.IntBetween(5, 30))
			promocode.FinishTime = finishTime
		}

		_, err = stmt.Exec(
			promocode.Code,
			promocode.Percentage,
			promocode.StartTime,
			promocode.FinishTime,
			promocode.IsAlways)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *SeedingV6) insertOrders() error {
	tx, err := s.service.Db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		//nolint:errcheck
		tx.Rollback()
	}()

	stmt, err := tx.Prepare("INSERT INTO orders (" +
		"booking_id, promocode_id, platform_fee, total_cost, status, security_code, created_at) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	approvedBookings, err := s.getBookings()
	if err != nil {
		return err
	}

	for _, book := range approvedBookings {
		order := &Order{
			BookingID:    book.BookingID,
			Status:       s.randomOrderStatus(),
			SecurityCode: strings.ToUpper(s.service.Fake.RandomStringWithLength(6)),
			TotalCost:    book.Price,
			PlatformFee:  book.Price * 0.2,
			CreatedAt: s.service.Fake.Time().TimeBetween(
				time.Now().AddDate(0, -1, -10), time.Now()),
		}

		if s.service.Fake.Float32(2, 0, 1) < 0.4 {
			promocodeID := s.service.Fake.IntBetween(1, s.service.GetPromocodeCount()-1)
			order.PromocodeID = &promocodeID
		}
		order.TotalCost += order.PlatformFee

		_, err := stmt.Exec(
			order.BookingID,
			order.PromocodeID,
			order.PlatformFee,
			order.TotalCost,
			order.Status,
			order.SecurityCode,
			order.CreatedAt)
		if err != nil {
			return err
		}

		s.service.SetOrderCount(s.service.GetOrderCount() + 1)
	}

	return tx.Commit()
}

func (s *SeedingV6) getBookings() ([]struct {
	BookingID int
	Price     float64
}, error) {
	rows, err := s.service.Db.Query("SELECT b.booking_id, ms.price + COALESCE(ad.offer_price, 0) " +
		"FROM booking b JOIN model_services ms " +
		"ON ms.model_service_id = b.model_service_id " +
		"LEFT JOIN additional_services ad " +
		"ON b.additional_service_id = ad.additional_service_id " +
		"WHERE b.status = 'APPROVED' AND ad.status IN ('APPROVED', 'CANCELLED')")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []struct {
		BookingID int
		Price     float64
	}
	for rows.Next() {
		var booking struct {
			BookingID int
			Price     float64
		}
		if err := rows.Scan(&booking.BookingID, &booking.Price); err != nil {
			return nil, err
		}
		bookings = append(bookings, booking)
	}

	return bookings, nil
}

func (s *SeedingV6) insertReviews() error {
	tx, err := s.service.Db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		//nolint:errcheck
		tx.Rollback()
	}()

	stmt, err := tx.Prepare("INSERT INTO reviews (" +
		"order_id, from_user_id, to_user_id, rating, description, created_at) " +
		"VALUES ($1, $2, $3, $4, $5, $6)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	orders, err := s.getOrders()
	if err != nil {
		return err
	}

	for _, order := range orders {
		err := s.formReview(
			order.OrderID,
			order.ClientID,
			order.ModelID,
			stmt,
		)
		if err != nil {
			return err
		}

		err = s.formReview(
			order.OrderID,
			order.ModelID,
			order.ClientID,
			stmt,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *SeedingV6) formReview(orderID, fromID, toID int, stmt *sql.Stmt) error {
	description := s.service.Fake.Lorem().Sentence(5)
	if len(description) > 500 {
		description = description[:500]
	}

	_, err := stmt.Exec(
		orderID,
		fromID,
		toID,
		s.service.Fake.IntBetween(1, 5),
		description,
		s.service.Fake.Time().TimeBetween(
			time.Now().AddDate(0, -1, -10), time.Now()),
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *SeedingV6) getOrders() ([]struct {
	OrderID  int
	ClientID int
	ModelID  int
}, error) {
	rows, err := s.service.Db.Query("SELECT o.order_id, b.client_id, ms.model_id " +
		"FROM orders o JOIN booking b ON o.booking_id = b.booking_id " +
		"JOIN model_services ms ON b.model_service_id = ms.model_service_id " +
		"WHERE o.status IN ('COMPLETED', 'CANCELLED')")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []struct {
		OrderID  int
		ClientID int
		ModelID  int
	}

	for rows.Next() {
		var order struct {
			OrderID  int
			ClientID int
			ModelID  int
		}
		if err := rows.Scan(&order.OrderID, &order.ClientID, &order.ModelID); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (s *SeedingV6) randomOrderStatus() string {
	num := s.service.Fake.Float32(2, 0, 1)
	switch {
	case num < 0.5:
		return "COMPLETED"
	case num < 0.6:
		return cancelled
	case num < 0.7:
		return "IN_PROCESS"
	case num < 0.8:
		return "IN_TRANSIT"
	default:
		return "CONFIRMED"
	}
}
