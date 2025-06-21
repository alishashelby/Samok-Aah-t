package services

import (
	"github.com/alishashelby/Samok-Aah-t/simulator/internal/metrics"
	"log"
	"time"
)

type Query struct {
	Name string
	Run  func() error
}

type Simulator struct {
	queries  []*Query
	metrics  *metrics.QueryMetrics
	interval time.Duration
}

func NewSimulator(metrics *metrics.QueryMetrics, interval time.Duration) *Simulator {
	return &Simulator{
		metrics:  metrics,
		interval: interval,
	}
}

func (s *Simulator) AddQuery(queryName string, function func() error) {
	s.queries = append(s.queries, &Query{
		Name: queryName,
		Run:  function,
	})
}

func (s *Simulator) Execute() {
	log.Println("Starting simulator with interval", s.interval)
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	s.executeQueries()

	for range ticker.C {
		s.executeQueries()
	}
}

func (s *Simulator) executeQueries() {
	for _, q := range s.queries {
		if err := s.metrics.Middleware(q.Name, q.Run); err != nil {
			log.Printf("error executing query %s: %v", q.Name, err)
		}
	}
}
