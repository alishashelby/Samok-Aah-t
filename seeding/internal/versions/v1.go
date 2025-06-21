package versions

import (
	"github.com/alishashelby/Samok-Aah-t/seeding/internal/service"
	"log"
)

type SeedingV1 struct {
	service *service.SeedingService
}

func NewSeedingV1(seedingService *service.SeedingService) *SeedingV1 {
	return &SeedingV1{
		service: seedingService,
	}
}

type Region struct {
	ID   int
	Name string
}

type City struct {
	ID       int
	Name     string
	RegionID int
}

func (s *SeedingV1) FakeV1() error {
	if !s.service.CheckTable("regions") ||
		!s.service.CheckTable("cities") {
		return nil
	}

	if !s.service.CheckIfFilled("regions") {
		log.Println("v1.go regions are not filled")
		if err := s.insertRegion(); err != nil {
			return err
		}
		log.Println("filled")
	}

	if !s.service.CheckIfFilled("cities") {
		log.Println("v1.go cities are not filled")
		if err := s.insertCities(); err != nil {
			return err
		}
		log.Println("filled")
	}

	return nil
}

func (s *SeedingV1) insertRegion() error {
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
		if len(name) > 40 {
			name = name[:40]
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

func (s *SeedingV1) insertCities() error {
	tx, err := s.service.Db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		//nolint:errcheck
		tx.Rollback()
	}()

	stmt, err := tx.Prepare("INSERT INTO cities (region_id, name) VALUES ($1, $2)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for range s.service.GetCityCount() {
		name := s.service.Fake.Address().City()
		regionID := s.service.Fake.IntBetween(1, s.service.GetRegionCount()-1)

		if len(name) > 20 {
			name = name[:20]
		}

		if _, err := stmt.Exec(regionID, name); err != nil {
			return err
		}
	}

	return tx.Commit()
}
