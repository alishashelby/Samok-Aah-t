package service

import (
	"database/sql"
	"fmt"
	"github.com/jaswdr/faker/v2"
	"sort"
)

type Node struct {
	Version      int
	DependsOn    []int
	SeedFunction func() error
}

type Config struct {
	SeedingService *SeedingService
	Node           map[int]Node
}

func NewConfig(service *SeedingService) *Config {
	return &Config{
		SeedingService: service,
		Node:           make(map[int]Node),
	}
}

func (s *Config) AddSeed(version int, deps []int, handler func() error) {
	s.Node[version] = Node{
		Version:      version,
		DependsOn:    deps,
		SeedFunction: handler,
	}
}

func (s *Config) dfs(version int, visited, temp map[int]bool, topSort *[]Node) error {
	if temp[version] {
		return fmt.Errorf("err: cycle %d", version)
	}
	if visited[version] {
		return nil
	}

	temp[version] = true
	seed, exists := s.Node[version]
	if !exists {
		return fmt.Errorf("err: no seed node for version %d", version)
	}

	for _, dep := range seed.DependsOn {
		if err := s.dfs(dep, visited, temp, topSort); err != nil {
			return err
		}
	}

	temp[version] = false
	visited[version] = true
	*topSort = append(*topSort, seed)
	return nil
}

func (s *Config) TopSort() ([]Node, error) {
	visited := make(map[int]bool)
	temp := make(map[int]bool)
	topSort := make([]Node, 0)

	versions := make([]int, len(s.Node))
	for v := range s.Node {
		versions = append(versions, v)
	}
	sort.Ints(versions)

	for _, v := range versions {
		if !visited[v] {
			if err := s.dfs(v, visited, temp, &topSort); err != nil {
				return nil, err
			}
		}
	}

	return topSort, nil
}

type SeedingService struct {
	Db    *sql.DB
	Fake  faker.Faker
	Count SeedCount
}

func NewSeedingService(db *sql.DB, seedCount int) *SeedingService {
	s := &SeedingService{
		Db:   db,
		Fake: faker.New(),
		Count: SeedCount{
			regionCount:      seedCount,
			cityCount:        seedCount,
			authCount:        seedCount,
			adminCount:       3,
			banCount:         10,
			categoryCount:    5,
			serviceCount:     15,
			bookingCount:     seedCount,
			promocodeCount:   20,
			transactionCount: seedCount * 2,
		},
	}
	s.Count.userCount = seedCount - s.Count.adminCount
	s.Count.clientCount = s.Count.userCount / 2
	s.Count.modelCount = s.Count.userCount - s.Count.clientCount

	return s
}

func (s *SeedingService) CheckTable(name string) bool {
	var exists bool
	query := `
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' AND table_name = $1
		)
	`
	err := s.Db.QueryRow(query, name).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

func (s *SeedingService) CheckIfFilled(name string) bool {
	var count int
	query := fmt.Sprintf("SELECT count(*) FROM %s", name)

	err := s.Db.QueryRow(query).Scan(&count)

	return err == nil && count > 0
}

type SeedCount struct {
	regionCount      int
	cityCount        int
	authCount        int
	adminCount       int
	userCount        int
	banCount         int
	clientCount      int
	modelCount       int
	categoryCount    int
	serviceCount     int
	bookingCount     int
	promocodeCount   int
	orderCount       int
	transactionCount int
}

func (s *SeedingService) GetRegionCount() int      { return s.Count.regionCount }
func (s *SeedingService) GetCityCount() int        { return s.Count.cityCount }
func (s *SeedingService) GetAuthCount() int        { return s.Count.authCount }
func (s *SeedingService) GetAdminCount() int       { return s.Count.adminCount }
func (s *SeedingService) GetUserCount() int        { return s.Count.userCount }
func (s *SeedingService) GetBanCount() int         { return s.Count.banCount }
func (s *SeedingService) GetClientCount() int      { return s.Count.clientCount }
func (s *SeedingService) GetModelCount() int       { return s.Count.modelCount }
func (s *SeedingService) GetCategoryCount() int    { return s.Count.categoryCount }
func (s *SeedingService) GetServiceCount() int     { return s.Count.serviceCount }
func (s *SeedingService) GetBookingCount() int     { return s.Count.bookingCount }
func (s *SeedingService) GetPromocodeCount() int   { return s.Count.promocodeCount }
func (s *SeedingService) GetOrderCount() int       { return s.Count.orderCount }
func (s *SeedingService) SetOrderCount(val int)    { s.Count.orderCount = val }
func (s *SeedingService) GetTransactionCount() int { return s.Count.transactionCount }
