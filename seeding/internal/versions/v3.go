package versions

import (
	"github.com/alishashelby/Samok-Aah-t/seeding/internal/service"
	"log"
)

type SeedingV3 struct {
	service *service.SeedingService
}

func NewSeedingV3(seedingService *service.SeedingService) *SeedingV3 {
	return &SeedingV3{
		service: seedingService,
	}
}

type LoyaltyLevel struct {
	ID                 int
	Name               string
	MinOrders          int
	CashbackPercentage int
}

type Client struct {
	ID             int
	UserID         int
	Name           string
	LoyaltyLevelID int
}

func (s *SeedingV3) FakeV3() error {
	if !s.service.CheckTable("loyalty_levels") ||
		!s.service.CheckTable("clients") {
		return nil
	}

	levels := []LoyaltyLevel{
		{Name: "BronzeBunny", MinOrders: 4, CashbackPercentage: 2},
		{Name: "SilverBunny", MinOrders: 8, CashbackPercentage: 3},
		{Name: "GoldBunny", MinOrders: 12, CashbackPercentage: 5},
		{Name: "PlatinumBunny", MinOrders: 20, CashbackPercentage: 10},
	}

	if !s.service.CheckIfFilled("loyalty_levels") {
		log.Println("v3.go loyalty_levels are not filled")
		if err := s.insertLoyaltyLevels(levels); err != nil {
			return err
		}
		log.Println("filled")
	}

	if !s.service.CheckIfFilled("clients") {
		log.Println("v3.go clients are not filled")
		if err := s.insertClients(len(levels)); err != nil {
			return err
		}
		log.Println("filled")
	}

	return nil
}

func (s *SeedingV3) insertLoyaltyLevels(levels []LoyaltyLevel) error {
	tx, err := s.service.Db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for i := range levels {
		err := tx.QueryRow(
			"INSERT INTO loyalty_levels (name, min_orders, cashback_percentage)"+
				"VALUES ($1, $2, $3) RETURNING level_id",
			levels[i].Name,
			levels[i].MinOrders,
			levels[i].CashbackPercentage,
		).Scan(&levels[i].ID)

		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *SeedingV3) insertClients(levelCount int) error {
	tx, err := s.service.Db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	rows, err := tx.Query("SELECT u.user_id FROM users u "+
		"WHERE NOT EXISTS (SELECT 1 FROM clients WHERE user_id = u.user_id) "+
		"AND NOT EXISTS (SELECT 1 FROM models WHERE user_id = u.user_id) "+
		"ORDER BY RANDOM() LIMIT $1", s.service.GetClientCount())
	if err != nil {
		return err
	}
	defer rows.Close()

	var userIDs []int
	for rows.Next() {
		var userID int
		if err := rows.Scan(&userID); err != nil {
			return err
		}
		userIDs = append(userIDs, userID)
	}
	if err := rows.Err(); err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO clients (user_id, name, loyalty_level_id) " +
		"VALUES ($1, $2, $3)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, userID := range userIDs {
		levelID := s.service.Fake.IntBetween(1, levelCount)
		_, err = stmt.Exec(userID, s.service.Fake.Person().Name(), levelID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
