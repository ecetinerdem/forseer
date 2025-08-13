package database

import (
	"fmt"
	"os"
)

func RunMigrations(db *DB) error {
	fmt.Println("📄 Running database migrations...")

	content, err := os.ReadFile("sql/schema.sql")

	if err != nil {
		return fmt.Errorf("could not read from migration file, %w", err)
	}

	_, err = db.Exec(string(content))

	if err != nil {
		return fmt.Errorf("could not run migration code, %w", err)
	}

	fmt.Println("✅ Migrations completed successfully!")
	return nil

}
