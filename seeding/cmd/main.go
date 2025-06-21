package main

import (
	"database/sql"
	"fmt"
	"github.com/alishashelby/Samok-Aah-t/seeding/internal/service"
	"github.com/alishashelby/Samok-Aah-t/seeding/internal/versions"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	if os.Getenv("APP_ENV") != "dev" {
		log.Println("Err: APP_ENV is", os.Getenv("APP_ENV"))
		return
	}

	seedCount, err := strconv.Atoi(os.Getenv("SEED_COUNT"))
	if err != nil {
		log.Println("Err: SEED_COUNT is", os.Getenv("SEED_COUNT"))
		return
	}

	if seedCount <= 0 {
		seedCount = 50
	}

	db, err := sql.Open("pgx", fmt.Sprintf(
		"postgres://%s:%s@db:5432/%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	))
	if err != nil {
		log.Println("Err:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Println("Err:", err)
	}

	seedingService := service.NewSeedingService(db, seedCount)
	handler := service.NewConfig(seedingService)

	initialize(handler)
	migrations, err := handler.TopSort()
	if err != nil {
		log.Println(err)
		return
	}

	order := make([]string, len(migrations))
	for _, m := range migrations {
		order = append(order, strconv.Itoa(m.Version))
	}

	log.Println("Order of versions:", strings.Join(order, ", "))

	for _, migration := range migrations {
		if err := migration.SeedFunction(); err != nil {
			log.Printf("Fatal err: %s", err)
			return
		}
	}

	log.Println("Seeding service is ready")
}

func initialize(config *service.Config) {

	s1 := versions.NewSeedingV1(config.SeedingService)
	s2 := versions.NewSeedingV2(config.SeedingService)
	s3 := versions.NewSeedingV3(config.SeedingService)
	s4 := versions.NewSeedingV4(config.SeedingService)
	s5 := versions.NewSeedingV5(config.SeedingService)
	s6 := versions.NewSeedingV6(config.SeedingService)
	s7 := versions.NewSeedingV7(config.SeedingService)
	s8 := versions.NewSeedingV8(config.SeedingService)
	s9 := versions.NewSeedingV9(config.SeedingService)

	config.AddSeed(1, []int{9}, s1.FakeV1)
	config.AddSeed(2, []int{1}, s2.FakeV2)
	config.AddSeed(3, []int{2}, s3.FakeV3)
	config.AddSeed(4, []int{3}, s4.FakeV4)
	config.AddSeed(5, []int{4}, s5.FakeV5)
	config.AddSeed(6, []int{5}, s6.FakeV6)
	config.AddSeed(7, []int{6}, s7.FakeV7)
	config.AddSeed(8, []int{7}, s8.FakeV8)
	config.AddSeed(9, []int{}, s9.FakeV9)
}
