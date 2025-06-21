package versions

import (
	"github.com/alishashelby/Samok-Aah-t/seeding/internal/service"
	"log"
	"strconv"
	"time"
)

type SeedingV4 struct {
	service *service.SeedingService
}

func NewSeedingV4(seedingService *service.SeedingService) *SeedingV4 {
	return &SeedingV4{
		service: seedingService,
	}
}

type Model struct {
	ID     int
	UserID int
	Name   string
}

type SocialMedia struct {
	ModelID  int
	Platform string
	URL      string
}

type PortfolioData struct {
	ModelID     int
	MediaURL    string
	Description string
	UploadedAt  time.Time
}

// nolint:dupl
func (s *SeedingV4) FakeV4() error {
	if !s.service.CheckTable("models") ||
		!s.service.CheckTable("social_media") ||
		!s.service.CheckTable("portfolio_data") {
		return nil
	}

	if !s.service.CheckIfFilled("models") {
		log.Println("v4.go models are not filled")
		if err := s.insertModels(); err != nil {
			return err
		}
		log.Println("filled")
	}

	if !s.service.CheckIfFilled("social_media") {
		log.Println("v4.go social_media is not filled")
		if err := s.insertSocialMedia(); err != nil {
			return err
		}
		log.Println("filled")
	}

	if !s.service.CheckIfFilled("portfolio_data") {
		log.Println("v4.go portfolio_data is not filled")
		if err := s.insertPortfolioData(); err != nil {
			return err
		}
		log.Println("filled")
	}

	return nil
}

func (s *SeedingV4) insertModels() error {
	tx, err := s.service.Db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		//nolint:errcheck
		tx.Rollback()
	}()

	rows, err := tx.Query("SELECT u.user_id FROM users u "+
		"WHERE NOT EXISTS (SELECT 1 FROM clients WHERE user_id = u.user_id) "+
		"AND NOT EXISTS (SELECT 1 FROM models WHERE user_id = u.user_id) "+
		"ORDER BY RANDOM() LIMIT $1", s.service.GetModelCount())
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

	stmt, err := tx.Prepare("INSERT INTO models (user_id, name) " +
		"VALUES ($1, $2)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, userID := range userIDs {
		name := s.service.Fake.Person().Name()
		if len(name) > 50 {
			name = name[:50]
		}
		_, err = stmt.Exec(userID, name)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *SeedingV4) insertSocialMedia() error {
	platforms := []string{"Instagram", "Telegram", "Vk",
		"TikTok", "Facebook", "OnlyFans", "Twitter", "PornHub"}

	tx, err := s.service.Db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		//nolint:errcheck
		tx.Rollback()
	}()

	stmt, err := tx.Prepare("INSERT INTO social_media(model_id, platform, url) " +
		"VALUES ($1, $2, $3)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for ID := range s.service.GetModelCount() {
		count := s.service.Fake.IntBetween(1, 10)
		for i := 0; i < count; i++ {
			platform := platforms[s.service.Fake.IntBetween(0, len(platforms)-1)]
			url := s.service.Fake.Internet().URL()
			if len(url) > 255 {
				url = url[:255]
			}

			if _, err = stmt.Exec(ID+1, platform, url); err != nil {
				return err
			}

		}
	}

	return tx.Commit()
}

func (s *SeedingV4) insertPortfolioData() error {
	tx, err := s.service.Db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		//nolint:errcheck
		tx.Rollback()
	}()

	stmt, err := tx.Prepare("INSERT INTO portfolio_data(model_id, media_url," +
		"description, uploaded_at, is_verified) " +
		"VALUES ($1, $2, $3, $4, $5)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for ID := range s.service.GetModelCount() {
		count := s.service.Fake.IntBetween(1, 10)
		for i := 0; i < count; i++ {
			description := s.service.Fake.Lorem().Sentence(20)
			mediaURL := s.service.Fake.Internet().URL() + "/media/" + strconv.Itoa(ID+1)
			if len(mediaURL) > 255 {
				mediaURL = mediaURL[:255]
			}
			uploadedAt := s.service.Fake.Time().TimeBetween(
				time.Now().AddDate(-2, 0, 0), time.Now())
			isVerified := s.service.Fake.Bool()

			if _, err := stmt.Exec(ID+1, mediaURL, description, uploadedAt, isVerified); err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}
