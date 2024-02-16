package persistent

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jczornik/glacier_backup/persistent/artifacts"
	"github.com/jczornik/glacier_backup/persistent/tasks"
	"github.com/jczornik/glacier_backup/persistent/workflows"
)

const (
	schemaVersion   = 1
	setVersionQuery = "PRAGMA user_version = %d" // Cannot set version using `?` semantics
	getVersionQuery = "PRAGMA user_version"
)

func CheckAndUpdateSchema(c DBClient) error {
	db, err := c.OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	version, err := getVersion(db)
	if err != nil {
		return err
	}

	if version != schemaVersion {
		log.Printf("Current schema version is '%d' but expected '%d' - starting migretion\n", version, schemaVersion)
		if err := createSchema(db); err != nil {
			log.Println("Error while runing migration")
			return err
		}

		log.Println("Migration finished")
	} else {
		log.Println("Schema version is up to date")
	}

	return nil
}

func getVersion(db *sql.DB) (int, error) {
	var version int

	row := db.QueryRow(getVersionQuery)
	if err := row.Scan(&version); err != nil {
		return version, err
	}

	return version, nil
}

func setVersion(db *sql.Tx) error {
	_, err := db.Exec(fmt.Sprintf(setVersionQuery, schemaVersion))
	return err
}

func createOrRollback(tx *sql.Tx, fn func(*sql.Tx) error) error {
	if err := fn(tx); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			log.Println("Error while rolling back transaction")
			return rerr
		}

		return err
	}

	return nil
}

func createSchema(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	toCreate := []func(*sql.Tx) error{workflows.CreateTable, tasks.CreateTable, artifacts.CreateTable, setVersion}

	for _, fn := range toCreate {
		if err := createOrRollback(tx, fn); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		log.Println("Error while commiting changes")
		if rerr := tx.Rollback(); rerr != nil {
			log.Println("Error while rolling back transaction")
			return rerr
		}

		return err
	}

	return nil
}
