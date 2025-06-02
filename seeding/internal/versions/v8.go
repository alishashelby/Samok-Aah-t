package versions

import (
	"github.com/alishashelby/Samok-Aah-t/seeding/internal/service"
)

type SeedingV8 struct {
	service *service.SeedingService
}

func NewSeedingV8(service *service.SeedingService) *SeedingV8 {
	return &SeedingV8{
		service: service,
	}
}

func (s *SeedingV8) FakeV8() error {
	if !s.service.CheckTable("daily_statistics") {
		return nil
	}

	if _, err := s.service.Db.Exec("CALL update_statistics()"); err != nil {
		return err
	}

	return nil
}
