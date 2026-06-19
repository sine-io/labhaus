package integration

import (
	"context"
	"fmt"
	"testing"

	"github.com/labhaus/backend/internal/infrastructure/persistence"
	"gorm.io/gorm"
)

// TestDB holds test database configuration
type TestDB struct {
	DB *gorm.DB
}

// SetupTestDB creates a test database connection
// Note: Requires PostgreSQL running locally or via Docker
func SetupTestDB(t *testing.T) *TestDB {
	config := persistence.DBConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		DBName:   "labhaus_test",
		SSLMode:  "disable",
	}

	db, err := persistence.NewDB(config)
	if err != nil {
		t.Skipf("Skipping integration test: database not available: %v", err)
	}

	// Run migrations
	if err := persistence.AutoMigrate(db); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	return &TestDB{DB: db}
}

// Cleanup removes all test data
func (tdb *TestDB) Cleanup(t *testing.T) {
	// Truncate all tables
	tables := []string{"workflows", "users", "styles"}
	for _, table := range tables {
		if err := tdb.DB.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)).Error; err != nil {
			t.Logf("Warning: failed to truncate table %s: %v", table, err)
		}
	}
}

// Close closes the database connection
func (tdb *TestDB) Close() {
	sqlDB, _ := tdb.DB.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}
}

// WithTransaction runs a function in a transaction and rolls back
func (tdb *TestDB) WithTransaction(t *testing.T, fn func(*gorm.DB)) {
	tx := tdb.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			t.Fatalf("Transaction panicked: %v", r)
		}
	}()

	fn(tx)
	tx.Rollback()
}

// ClearContext returns a background context
func ClearContext() context.Context {
	return context.Background()
}
