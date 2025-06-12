package main

import (
	"database/sql"
	"fmt"
	"github.com/alishashelby/Samok-Aah-t/test_migrations/internal/service"
	_ "github.com/jackc/pgx/v4/stdlib"
	"log"
	"os"
)

func main() {
	url := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)
	db, err := sql.Open("pgx", url)
	if err != nil {
		log.Println("error connecting to database:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Println(err)
	}

	migrationTester := service.NewMigrationTester(db, "/goose/sql", url)
	err = migrationTester.Initialize()
	if err != nil {
		log.Println(err)
	}

	migrations, err := migrationTester.GetMigrations()
	if err != nil {
		log.Println(err)
	}
	log.Printf("found %d migrations\n", len(migrations))

	for _, migration := range migrations {
		if err = migrationTester.TestMigration(migration); err != nil {
			log.Fatal(err)
		}
	}

	log.Println("migrations tested successfully and are idempotent")
}
