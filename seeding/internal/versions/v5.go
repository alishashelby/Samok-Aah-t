package versions

import (
	"encoding/json"
	"fmt"
	"github.com/alishashelby/Samok-Aah-t/seeding/internal/service"
	"log"
	"time"
)

type SeedingV5 struct {
	service *service.SeedingService
}

func NewSeedingV5(seedingService *service.SeedingService) *SeedingV5 {
	return &SeedingV5{
		service: seedingService,
	}
}

type Category struct {
	ID   int
	Name string
}

type Service struct {
	ID          int
	CategoryID  int
	Description string
}

type ModelService struct {
	ID        int
	ModelID   int
	ServiceID int
	Price     float32
}

type AdditionalService struct {
	ID          int
	Description string
	OfferPrice  float64
	Status      string
	UpdatedAt   time.Time
}

type Booking struct {
	ID                  int
	ClientID            int
	ModelServiceID      int
	DateTime            time.Time
	Duration            time.Time
	Address             map[string]interface{}
	AdditionalServiceID int
	Status              string
	CreatedAt           time.Time
}

func (s *SeedingV5) FakeV5() error {
	if !s.service.CheckTable("categories") ||
		!s.service.CheckTable("services") ||
		!s.service.CheckTable("model_services") ||
		!s.service.CheckTable("additional_services") ||
		!s.service.CheckTable("booking") {
		return nil
	}

	if !s.service.CheckIfFilled("categories") {
		log.Println("v5.go categories are not filled")
		if err := s.insertCategories(); err != nil {
			return err
		}
		log.Println("filled")
	}

	if !s.service.CheckIfFilled("services") {
		log.Println("v5.go services are not filled")
		if err := s.insertServices(); err != nil {
			return err
		}
		log.Println("filled")
	}

	var rez int
	if !s.service.CheckIfFilled("model_services") {
		log.Println("v5.go model_services are not filled")
		modelServicesCount, err := s.insertModelServices()
		if err != nil {
			return err
		}
		log.Println("filled")
		rez = modelServicesCount
	}

	if s.service.CheckIfFilled("additional_services") {
		return fmt.Errorf("v5.go additional_services are filled")
	}

	if !s.service.CheckIfFilled("booking") {
		log.Println("v5.go booking are not filled")
		if err := s.insertBooking(rez); err != nil {
			return err
		}
		log.Println("filled")
	}

	return nil
}

func (s *SeedingV5) insertCategories() error {
	tx, err := s.service.Db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO categories (name) VALUES ($1) RETURNING category_id")
	if err != nil {
		return err
	}
	defer stmt.Close()

	uniqueNames := make(map[string]bool)
	for range s.service.GetCategoryCount() {
		var name string
		for {
			name = s.service.Fake.Company().JobTitle()
			if !uniqueNames[name] {
				break
			}
		}
		uniqueNames[name] = true
		if len(name) > 70 {
			name = name[:70]
		}

		_, err = stmt.Exec(name)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *SeedingV5) insertServices() error {
	tx, err := s.service.Db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO services (category_id, description) " +
		"VALUES ($1, $2) RETURNING service_id")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for range s.service.GetServiceCount() {
		categoryID := s.service.Fake.IntBetween(1, s.service.GetCategoryCount()-1)
		description := s.service.Fake.Lorem().Paragraph(3)
		if len(description) > 255 {
			description = description[:255]
		}

		_, err = stmt.Exec(categoryID, description)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *SeedingV5) insertModelServices() (int, error) {
	tx, err := s.service.Db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	modelCount := 0
	stmt, err := tx.Prepare("INSERT INTO model_services (model_id, service_id, price) " +
		"VALUES ($1, $2, $3) RETURNING model_service_id")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	for i := range s.service.GetModelCount() {
		count := s.service.Fake.IntBetween(1, 4)
		for j := 0; j < count; j++ {
			ms := ModelService{
				ModelID:   i + 1,
				ServiceID: s.service.Fake.IntBetween(1, s.service.GetServiceCount()-1),
				Price:     s.service.Fake.Float32(2, 5000, 100000),
			}
			err := stmt.QueryRow(
				ms.ModelID,
				ms.ServiceID,
				ms.Price,
			).Scan(&ms.ID)
			if err != nil {
				return 0, err
			}

			modelCount++
		}
	}

	return modelCount, tx.Commit()
}

func (s *SeedingV5) insertBooking(modelServicesCount int) error {
	tx, err := s.service.Db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO booking (" +
		"client_id, model_service_id, date_time, duration, " +
		"address, status, created_at) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7)" +
		"RETURNING booking_id")
	if err != nil {
		return err
	}
	defer stmt.Close()

	stmtAdditional, err := tx.Prepare("INSERT INTO additional_services (" +
		"description, offer_price, status, updated_at) " +
		"VALUES ($1, $2, $3, $4) RETURNING additional_service_id")
	if err != nil {
		return err
	}
	defer stmtAdditional.Close()

	stmtUpdate, err := tx.Prepare("" +
		"UPDATE booking SET additional_service_id = $1 WHERE booking_id = $2")
	if err != nil {
		return err
	}
	defer stmtUpdate.Close()

	for range s.service.GetBookingCount() {
		clientID := s.service.Fake.IntBetween(1, s.service.GetClientCount()-1)
		modelServiceID := s.service.Fake.IntBetween(1, modelServicesCount-1)

		address := map[string]interface{}{
			"city":    s.service.Fake.Address().City(),
			"street":  s.service.Fake.Address().StreetAddress(),
			"details": s.service.Fake.Lorem().Sentence(6),
		}
		addressJSON, err := json.Marshal(address)
		if err != nil {
			return err
		}

		var additionalServiceID int

		if s.service.Fake.Bool() {
			description := s.service.Fake.Lorem().Paragraph(3)
			if len(description) > 1000 {
				description = description[:1000]
			}

			price := s.service.Fake.Float32(2, 5000, 100000)
			status := s.randomAdditionalStatus()
			updatedAt := s.service.Fake.Time().TimeBetween(
				time.Now().AddDate(-1, 0, 0), time.Now())

			err = stmtAdditional.QueryRow(
				description,
				price,
				status,
				updatedAt,
			).Scan(&additionalServiceID)
			if err != nil {
				return err
			}
		}

		var bookingID int

		err = stmt.QueryRow(
			clientID,
			modelServiceID,
			s.service.Fake.Time().TimeBetween(
				time.Now().AddDate(0, 0, 10),
				time.Now().AddDate(0, 3, 0)),
			fmt.Sprintf("%d hours", s.service.Fake.IntBetween(2, 12)),
			addressJSON,
			s.randomStatus(),
			s.service.Fake.Time().TimeBetween(
				time.Now().AddDate(0, 0, 0), time.Now()),
		).Scan(&bookingID)
		if err != nil {
			return err
		}

		if additionalServiceID != 0 {
			_, err = stmtUpdate.Exec(additionalServiceID, bookingID)
		}
	}

	return tx.Commit()
}

func (s *SeedingV5) randomAdditionalStatus() string {
	num := s.service.Fake.Float32(2, 0, 1)
	switch {
	case num < 0.5:
		return "APPROVED"
	case num < 0.6:
		return "CANCELLED"
	case num < 0.7:
		return "PENDING"
	case num < 0.8:
		return "HIGHER_PRICE"
	default:
		return "REJECTED"
	}
}

func (s *SeedingV5) randomStatus() string {
	num := s.service.Fake.Float32(2, 0, 1)
	switch {
	case num < 0.5:
		return "APPROVED"
	case num < 0.6:
		return "CANCELLED"
	case num < 0.7:
		return "PENDING"
	default:
		return "REJECTED"
	}
}
