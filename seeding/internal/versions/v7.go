package versions

import (
	"fmt"
	"github.com/alishashelby/Samok-Aah-t/seeding/internal/service"
	"log"
	"time"
)

type SeedingV7 struct {
	service *service.SeedingService
}

func NewSeedingV7(service *service.SeedingService) *SeedingV7 {
	return &SeedingV7{
		service: service,
	}
}

type PaymentSystem struct {
	ID   int
	Name string
}

type ExternalTransaction struct {
	ID              int
	PaymentSystemID int
	FailureMsg      *string
}

type Transaction struct {
	ID                    int
	Amount                float64
	Type                  string
	OrderID               *int
	ExternalTransactionID *int
	Reason                *string
	Status                string
	CreatedAt             time.Time
	ProcessedAt           *time.Time
}

func (s *SeedingV7) FakeV7() error {
	if !s.service.CheckTable("payment_system_integration") ||
		!s.service.CheckTable("external_transactions") ||
		!s.service.CheckTable("transactions") {
		return nil
	}

	var rez int
	if !s.service.CheckIfFilled("payment_system_integration") {
		log.Println("v7.go payment_system_integration are not filled")
		paymentCount, err := s.insertPaymentSystemIntegration()
		if err != nil {
			return err
		}
		log.Println("filled")
		rez = paymentCount
	}

	if s.service.CheckIfFilled("external_transactions") {
		return fmt.Errorf("v7.go external_transactions are filled")
	}

	if !s.service.CheckIfFilled("transactions") {
		log.Println("v7.go transactions are not filled")
		if err := s.insertTransactions(rez); err != nil {
			return err
		}
		log.Println("filled")
	}

	return nil
}

func (s *SeedingV7) insertPaymentSystemIntegration() (int, error) {
	tx, err := s.service.Db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	paymentSystems := []string{"Sberbank", "Tinkoff", "YandexBank", "Stripe", "PayPal"}

	stmt, err := tx.Prepare("INSERT INTO payment_system_integration (name)" +
		"VALUES ($1)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	for _, paymentSystem := range paymentSystems {
		if _, err := stmt.Exec(paymentSystem); err != nil {
			return 0, err
		}
	}

	return len(paymentSystems), tx.Commit()
}

func (s *SeedingV7) insertTransactions(paymentCount int) error {
	tx, err := s.service.Db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO transactions (" +
		"amount, type, order_id, external_transaction_id, " +
		"reason, status, created_at, processed_at) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8) ")
	if err != nil {
		return err
	}
	defer stmt.Close()

	transactionTypes := []string{
		"ORDER_PAYMENT", "ORDER_INCOME", "CLIENT_DEPOSIT", "MODEL_PAYOUT",
		"REFUND_TO_CLIENT", "REFUND_FROM_MODEL", "ORDER_CANCELLATION",
		"REFERRAL", "CASHBACK",
	}

	stmtExternal, err := tx.Prepare("INSERT INTO external_transactions " +
		"(payment_system_id, failure_msg) " +
		"VALUES ($1, $2) RETURNING external_transaction_id")
	if err != nil {
		return err
	}
	defer stmtExternal.Close()

	for range s.service.GetTransactionCount() {
		tType := transactionTypes[s.service.Fake.IntBetween(0, len(transactionTypes)-1)]
		transaction := Transaction{
			Amount: s.service.Fake.Float64(2, 100, 10000),
			Type:   tType,
			Status: s.randomTransactionStatus(),
			CreatedAt: s.service.Fake.Time().TimeBetween(
				time.Now().AddDate(-3, 0, 0), time.Now()),
		}

		switch tType {
		case "CLIENT_DEPOSIT", "MODEL_PAYOUT":
			var failureMsg *string
			if s.service.Fake.Float32(2, 0, 1) < 0.3 {
				msg := s.service.Fake.Lorem().Sentence(5)
				if len(msg) > 255 {
					msg = msg[:255]
				}
				failureMsg = &msg
			}

			err := stmtExternal.QueryRow(
				paymentCount,
				failureMsg,
			).Scan(&transaction.ExternalTransactionID)

			if err != nil {
				return err
			}
		case "ORDER_PAYMENT", "ORDER_INCOME", "CASHBACK":
			orderID := s.service.Fake.IntBetween(1, s.service.GetOrderCount()-1)
			transaction.OrderID = &orderID
		case "REFUND_TO_CLIENT", "REFUND_FROM_MODEL", "ORDER_CANCELLATION":
			orderID := s.service.Fake.IntBetween(1, s.service.GetOrderCount()-1)
			transaction.OrderID = &orderID

			reason := s.service.Fake.Lorem().Sentence(5)
			if len(reason) > 255 {
				reason = reason[:255]
			}
			transaction.Reason = &reason
		}

		if transaction.Status == "SUCCESS" {
			processing := transaction.CreatedAt.Add(time.Minute * 4)
			transaction.ProcessedAt = &processing
		}

		_, err = stmt.Exec(
			transaction.Amount,
			transaction.Type,
			transaction.OrderID,
			transaction.ExternalTransactionID,
			transaction.Reason,
			transaction.Status,
			transaction.CreatedAt,
			transaction.ProcessedAt,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *SeedingV7) randomTransactionStatus() string {
	num := s.service.Fake.Float32(2, 0, 1)
	switch {
	case num < 0.5:
		return "SUCCESS"
	case num < 0.8:
		return "PENDING"
	default:
		return "FAILURE"
	}
}
