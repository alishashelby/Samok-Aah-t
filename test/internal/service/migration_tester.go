package service

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type MigrationTester struct {
	db           *sql.DB
	migrationDir string
	url          string
}

func NewMigrationTester(db *sql.DB, migrationDir string, url string) *MigrationTester {
	return &MigrationTester{
		db:           db,
		migrationDir: migrationDir,
		url:          url,
	}
}

func (m *MigrationTester) Initialize() error {
	_, err := m.db.Exec(`CREATE TABLE IF NOT EXISTS goose_db_version(
    	id SERIAL PRIMARY KEY,
    	version_id BIGINT NOT NULL,
		is_applied BOOLEAN NOT NULL,
		tstamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating migration table: %v", err)
	}

	_, err = m.db.Exec(`INSERT INTO goose_db_version(version_id, is_applied)
		SELECT 0, true
		WHERE NOT EXISTS (SELECT 1 FROM goose_db_version)
	`)
	if err != nil {
		return fmt.Errorf("error initializing migration table: %v", err)
	}

	return nil
}

func (m *MigrationTester) GetMigrations() ([]string, error) {
	files, err := os.ReadDir(m.migrationDir)
	if err != nil {
		return nil, fmt.Errorf("error reading migration dir: %v", err)
	}

	var migrations []string
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".sql") {
			migrations = append(migrations, f.Name())
		}
	}

	sort.Strings(migrations)
	return migrations, nil
}

func (m *MigrationTester) TestMigration(migration string) error {
	log.Printf("testing migration: %s\n", migration)

	if err := m.runMigrationUpByOne(); err != nil {
		return fmt.Errorf("error running migration up-by-one: %v", err)
	}

	snapshot1, err := m.createSnapshot()
	if err != nil {
		return fmt.Errorf("error creating snapshot: %v", err)
	}
	defer os.Remove(*snapshot1)

	if err := m.runMigrationDown(); err != nil {
		return fmt.Errorf("error running migration down: %v", err)
	}

	if err := m.runMigrationUpByOne(); err != nil {
		return fmt.Errorf("error running migration up-by-one again: %v", err)
	}

	snapshot2, err := m.createSnapshot()
	if err != nil {
		return fmt.Errorf("error creating snapshot: %v", err)
	}
	defer os.Remove(*snapshot2)

	equal, err := m.compareSnapshots(*snapshot1, *snapshot2)
	if err != nil {
		return fmt.Errorf("error comparing snapshots: %v", err)
	}
	if !equal {
		return fmt.Errorf("FAILURE! Migration %s is not idempotent", migration)
	}

	log.Printf("SUCCESS! Migration %s is idempotent\n", migration)
	return nil
}

func (m *MigrationTester) runMigrationUpByOne() error {
	return m.runGooseMigarion("up-by-one")
}

func (m *MigrationTester) runMigrationDown() error {
	return m.runGooseMigarion("down")
}

func (m *MigrationTester) runGooseMigarion(command string) error {
	cmd := exec.Command("goose",
		"-dir", m.migrationDir,
		"postgres", m.url, command)
	cmd.Stdout = os.Stdout

	return cmd.Run()
}

func (m *MigrationTester) createSnapshot() (*string, error) {
	file, err := os.CreateTemp("", "snapshot_*.sql")
	if err != nil {
		return nil, err
	}

	path := file.Name()
	defer file.Close()

	if err := m.dumpSchema(path); err != nil {
		log.Printf("error dumping schema: %v\n", err)
		return nil, err
	}

	return &path, nil
}

func (m *MigrationTester) dumpSchema(filename string) error {
	cmd := exec.Command("pg_dump", "-s",
		"-h", os.Getenv("POSTGRES_HOST"),
		"-p", os.Getenv("POSTGRES_PORT"),
		"-U", os.Getenv("POSTGRES_USER"),
		"-d", os.Getenv("POSTGRES_DB"),
	)
	cmd.Env = append(os.Environ(), "PGPASSWORD="+os.Getenv("POSTGRES_PASSWORD"))

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	cmd.Stdout = file
	return cmd.Run()
}

func (m *MigrationTester) compareSnapshots(file1, file2 string) (bool, error) {
	snapshot1, err := os.ReadFile(file1)
	if err != nil {
		return false, err
	}

	snapshot2, err := os.ReadFile(file2)
	if err != nil {
		return false, err
	}

	return bytes.Equal(snapshot1, snapshot2), nil
}
