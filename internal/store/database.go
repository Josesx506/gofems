package store

import (
	"database/sql"
	"fmt"
	"io/fs"
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib" // Import lib without using it
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
)

func Open() (*sql.DB, error) {

	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("Error loading .env file")
	}

	// dsn connection url or manual text
	// "host=localhost user=postgres password=postgres dbname=workoutDB port=5432 sslmode=disable"
	dsn := os.Getenv("DATABASE_URL")

	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is not set\n")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("db: open %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db: open %w", err)
	}

	// Additional security config
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	db.SetConnMaxIdleTime(time.Second * 60)

	fmt.Printf("Connected to database....\n")
	return db, nil
}

func MigrateFS(db *sql.DB, migrationsFS fs.FS, dir string) error {
	goose.SetBaseFS(migrationsFS)

	// Anonymous function run at the end that wipes global goose
	// state to default behavior
	defer func() {
		goose.SetBaseFS(nil)
	}()

	return Migrate(db, dir)
}

func Migrate(db *sql.DB, dir string) error {
	err := goose.SetDialect("postgres") // specify db type
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	err = goose.Up(db, dir)
	if err != nil {
		return fmt.Errorf("goose up: %w", err)
	}

	return nil
}

func SetupTestDB(t *testing.T, migrationDirectory string) *sql.DB {
	//  Define the host using the  service name in .devcontainer/docker-compose.yml
	db, err := sql.Open("pgx", "host=test_db user=postgres password=postgres dbname=postgres port=5432 sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Run the migrations for the test database
	err = Migrate(db, migrationDirectory) //e.g., "../../migrations/"
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	// Reset the database state before each test
	_, err = db.Exec(`Truncate workouts, workout_entries CASCADE`)
	if err != nil {
		t.Fatalf("Failed to truncate test database: %v", err)
	}

	return db
}
