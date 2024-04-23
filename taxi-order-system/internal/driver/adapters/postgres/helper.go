package repository

import (
	"errors"
	"fmt"
	"os"
	"qwe/internal/driver/config"

	"github.com/jmoiron/sqlx"
)

func openConnection(cfg *config.DB) (*sqlx.DB, error) {
	source := fmt.Sprintf("host=%v user=%v password=%v dbname=postgres sslmode=disable", cfg.HOST, cfg.USER, cfg.PWD)

	db, err := sqlx.Connect("postgres", source)

	return db, err
}

func readMigrations(cfg *config.DB) ([]string, error) {
	path, _ := os.Getwd()

	files, err := os.ReadDir(path + cfg.MigrationPath)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, errors.New("no migration files")
	}

	migrations := make([]string, 0, len(files))

	for _, file := range files {
		bts, err := os.ReadFile(path + cfg.MigrationPath + file.Name())
		if err != nil {
			return nil, err
		}
		migrations = append(migrations, string(bts))
	}

	return migrations, nil
}

func migrate(db *sqlx.DB, migrations []string) error {

	tx, _ := db.Begin()

	defer tx.Rollback()

	for _, v := range migrations {
		_, err := tx.Exec(v)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
