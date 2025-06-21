package versions

import (
	"github.com/alishashelby/Samok-Aah-t/seeding/internal/service"
	"log"
)

type SeedingV9 struct {
	service *service.SeedingService
}

func NewSeedingV9(seedingService *service.SeedingService) *SeedingV9 {
	return &SeedingV9{
		service: seedingService,
	}
}

func (s *SeedingV9) FakeV9() error {
	if !s.service.CheckTable("regions") {
		return nil
	}

	if !s.service.CheckIfFilled("regions") {
		log.Println("v9.go regions are not filled")
		if err := s.insertNewRegion(); err != nil {
			return err
		}
		log.Println("filled")
	}

	return nil
}

func (s *SeedingV9) insertNewRegion() error {
	tx, err := s.service.Db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		//nolint:errcheck
		tx.Rollback()
	}()

	for range s.service.GetRegionCount() {
		name := s.service.Fake.Address().State()
		if len(name) > 5 {
			name = name[:5]
		}
		_, err := tx.Exec(
			"INSERT INTO regions (name) VALUES ($1) RETURNING region_id",
			name,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
