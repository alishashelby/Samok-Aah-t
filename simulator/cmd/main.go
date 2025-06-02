package main

import (
	"database/sql"
	"fmt"
	"github.com/alishashelby/Samok-Aah-t/simulator/internal/metrics"
	"github.com/alishashelby/Samok-Aah-t/simulator/internal/queries"
	"github.com/alishashelby/Samok-Aah-t/simulator/internal/services"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
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

	queryMetrics := metrics.NewQueryMetrics()
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(":6969", nil); err != nil {
			log.Fatalf("Error serving /metrics: %v", err)
		}

		log.Println("Serving /metrics on :6969")
	}()

	simulator := initialize(db, queryMetrics)

	simulator.Execute()
}

func initialize(db *sql.DB, queryMetrics *metrics.QueryMetrics) *services.Simulator {
	mq := queries.NewModelQueries(db)
	oq := queries.NewOrderQueries(db)
	sq := queries.NewStatisticsQueries(db)
	tq := queries.NewTransactionQueries(db)

	interval, err := time.ParseDuration(os.Getenv("INTERVAL_SECONDS"))
	if err != nil {
		log.Println("error: no interval in dot env:", err)
		interval = 10 * time.Second
	}

	simulator := services.NewSimulator(queryMetrics, interval)
	simulator.AddQuery("filter_models_v1", mq.FilterModelsV1Query)
	simulator.AddQuery("filter_models_v2", mq.FilterModelsV2Query)
	simulator.AddQuery("admin_get_complicated_orders", oq.GetComplicatedOrdersQuery)
	simulator.AddQuery("get_all_data_to_calculate_total_cost", oq.GetAllDataToCalculateTotalCostQuery)
	simulator.AddQuery("get_admin_statistics", sq.GetAdminStatisticsQuery)
	simulator.AddQuery("get_referral_statistics", sq.GetReferralStatisticsQuery)
	simulator.AddQuery("get_loyalty_statistics", sq.GetClientLoyaltyStatisticsQuery)
	simulator.AddQuery("calculate_model_financial", tq.CalculateModelFinancialQuery)
	simulator.AddQuery("get_all_transactions", tq.GetAllTransactionsQuery)

	return simulator
}
