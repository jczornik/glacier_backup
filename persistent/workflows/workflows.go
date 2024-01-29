package workflows

import (
	"database/sql"
	"log"
)

const createTableQuery = `
    CREATE TABLE workflows
    (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        status TEXT NOT NULL
    )`

const (
	PendingStatus        = "pending"
	RunningStatus        = "running"
	FinishedStatus       = "finished"
	RollbackedStatus     = "rollbacked"
	FailedRollbackStatus = "failed_rollback"
)

func CreateTable(db *sql.Tx) error {
	_, err := db.Exec(createTableQuery)
	if err != nil {
		log.Println("Error while creating workflows table")
	}
	return err
}

func Create(tx *sql.Tx) (int64, error) {
	row, err := tx.Exec("INSERT INTO workflows (status) VALUES (?)", PendingStatus)
	if err != nil {
		log.Println("Error while creating workflow")
	}

	workflowid, _ := row.LastInsertId()
	log.Printf("Created workflow number %d\n", workflowid)

	return workflowid, err
}

func UpdateStatus(db *sql.DB, workflow int64, status string) error {
	_, err := db.Exec("UPDATE workflows SET status = ? WHERE id = ?", status, workflow)
	if err != nil {
		log.Printf("Error while updating status for workflow %d\n", workflow)
	}
	return err
}
