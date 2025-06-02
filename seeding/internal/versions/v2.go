package versions

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/alishashelby/Samok-Aah-t/seeding/internal/service"
	"github.com/google/uuid"
	"log"
	"strconv"
	"time"
)

type SeedingV2 struct {
	service *service.SeedingService
}

func NewSeedingV2(seedingService *service.SeedingService) *SeedingV2 {
	return &SeedingV2{
		service: seedingService,
	}
}

func (s *SeedingV2) getEmail() string {
	for {
		email := s.service.Fake.Internet().Email()
		var exists bool
		err := s.service.Db.QueryRow(
			"SELECT email FROM auth WHERE email = $1", email).Scan(&exists)
		if errors.Is(err, sql.ErrNoRows) {
			return email
		}
	}
}

func (s *SeedingV2) getPhone() string {
	for {
		phone := s.service.Fake.Phone().E164Number()
		var exists bool
		err := s.service.Db.QueryRow(
			"SELECT phone FROM auth WHERE phone = $1", phone).Scan(&exists)
		if errors.Is(err, sql.ErrNoRows) {
			return phone
		}
	}
}

type Auth struct {
	ID           int
	Email        string
	Phone        string
	PasswordHash string
	CreatedAt    time.Time
}

type Admin struct {
	ID          int
	AuthID      int
	Permissions string
}

type User struct {
	ID                int
	AuthID            int
	BirthDate         time.Time
	Gender            string
	CityID            int
	PassportSeries    string
	PassportNumber    string
	PassportIssueDate time.Time
	PassportVerified  bool
	ReferralCode      uuid.UUID
	ReferralUserID    int
	ReferralCount     int
	IsBanned          bool
}

type Ban struct {
	ID        int
	AdminID   int
	UserID    int
	Reason    string
	CreatedAt time.Time
}

func (s *SeedingV2) FakeV2() error {
	if !s.service.CheckTable("auth") ||
		!s.service.CheckTable("users") ||
		!s.service.CheckTable("ban") ||
		!s.service.CheckTable("admins") {
		return nil
	}

	if s.service.CheckIfFilled("auth") {
		return fmt.Errorf("err: v2.go auth is filled")
	}

	if !s.service.CheckIfFilled("admins") {
		log.Printf("v2.go admins are not filled")
		if err := s.insertAdmins(); err != nil {
			return err
		}
		log.Println("filled")
	}

	if !s.service.CheckIfFilled("users") {
		log.Printf("v2.go users are not filled")
		if err := s.insertUsers(); err != nil {
			return err
		}
		log.Println("filled")
	}

	if !s.service.CheckIfFilled("ban") {
		log.Printf("v2.go ban are not filled")
		if err := s.insertBans(); err != nil {
			return err
		}
		log.Println("filled")
	}

	return nil
}

func (s *SeedingV2) provideAuth() (Auth, error) {
	var auth Auth

	err := s.service.Db.QueryRow(
		"INSERT INTO auth (email, phone, password_hash, created_at) "+
			"VALUES ($1, $2, $3, $4) RETURNING auth_id",
		s.getEmail(),
		s.getPhone(),
		s.service.Fake.Hash().SHA256(),
		s.service.Fake.Time().TimeBetween(
			time.Now().AddDate(-3, 0, 0),
			time.Now().AddDate(0, -3, 0),
		),
	).Scan(&auth.ID)

	return auth, err
}

func (s *SeedingV2) insertAdmins() error {
	tx, err := s.service.Db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
        INSERT INTO admins (auth_id, permissions) 
        VALUES ($1, $2) 
        RETURNING admin_id`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for range s.service.GetAdminCount() {

		auth, err := s.provideAuth()
		if err != nil {
			return err
		}

		permissions := []string{
			`{"check passport": true}`,
			`{"check portfolio data": true}`,
			`{"check social media": true}`,
		}

		_, err = stmt.Exec(
			auth.ID,
			permissions[s.service.Fake.IntBetween(0, len(permissions)-1)],
		)

		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *SeedingV2) insertUsers() error {
	tx, err := s.service.Db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for range s.service.GetUserCount() {
		auth, err := s.provideAuth()
		if err != nil {
			return err
		}

		user := User{
			AuthID: auth.ID,
			BirthDate: s.service.Fake.Time().TimeBetween(
				time.Now().AddDate(-98, 0, 0),
				time.Now().AddDate(-18, 0, 0),
			),
			Gender:         s.service.Fake.Gender().Name(),
			CityID:         s.service.Fake.IntBetween(1, s.service.GetCityCount()-1),
			PassportSeries: strconv.Itoa(s.service.Fake.RandomNumber(4)),
			PassportNumber: strconv.Itoa(s.service.Fake.RandomNumber(6)),
			ReferralCode:   uuid.New(),
		}

		user.PassportIssueDate = user.BirthDate.AddDate(14, s.service.Fake.IntBetween(0, 11), s.service.Fake.IntBetween(1, 28))

		user.PassportVerified = s.service.Fake.Bool()

		err = s.service.Db.QueryRow(
			"INSERT INTO users ("+
				"auth_id, birth_date, gender, city_id, "+
				"passport_series, passport_number, passport_issue_date, "+
				"passport_verified, referral_code) "+
				"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING user_id",
			user.AuthID, user.BirthDate, user.Gender, user.CityID,
			user.PassportSeries, user.PassportNumber, user.PassportIssueDate,
			user.PassportVerified, user.ReferralCode,
		).Scan(&user.ID)
		if err != nil {
			return err
		}

		if s.service.Fake.Float32(2, 0, 1) < 0.4 {
			referrerID := s.service.Fake.IntBetween(1, s.service.GetUserCount()-1)
			rez, err := tx.Exec(
				"UPDATE users "+
					"SET referral_user_count = referral_user_count + 1 "+
					"WHERE user_id = $1 AND referral_user_count < 7",
				referrerID,
			)
			if err != nil {
				return err
			}

			updatedRows, err := rez.RowsAffected()
			if err != nil {
				return err
			}
			if updatedRows > 0 {
				user.ReferralUserID = referrerID
				_, err = tx.Exec(
					"UPDATE users "+
						"SET referral_user_id = $1 WHERE user_id = $2",
					referrerID, user.ID)
				if err != nil {
					return err
				}
			}
		}
	}

	return tx.Commit()
}

func (s *SeedingV2) insertBans() error {
	tx, err := s.service.Db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for range s.service.GetBanCount() {
		userID := s.service.Fake.IntBetween(1, s.service.GetUserCount()-1)
		adminID := s.service.Fake.IntBetween(1, s.service.GetAdminCount()-1)

		_, err := s.service.Db.Exec(
			"INSERT INTO ban ("+
				"admin_id, user_id, reason, created_at"+
				") VALUES ($1, $2, $3, $4)",
			adminID, userID, s.service.Fake.Lorem().Sentence(15),
			s.service.Fake.Time().TimeBetween(
				time.Now().AddDate(0, -7, 0),
				time.Now(),
			),
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
