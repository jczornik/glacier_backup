package tasks

import (
	"database/sql"
	"log"
)

const createTableQuery = `
    CREATE TABLE jobs (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        workflow INTEGER NOT NULL,
        status TEXT NOT NULL,
        FOREIGN KEY (workflow) REFERENCES workflows(id)
    )`

const (
	PendingStatus        = "pending"
	RunningStatus        = "running"
	FinishedStatus       = "finished"
	FailedStatus         = "failed"
	RollbackStatus       = "rollback"
	RollbackedStatus     = "rollbacked"
	FailedRollbackStatus = "failed_rollback"
)

func CreateTable(db *sql.Tx) error {
	_, err := db.Exec(createTableQuery)
	if err != nil {
		log.Println("Error while creating jobs table")
	}

	return err
}

func Create(tx *sql.Tx, workflow int64, name string) (int64, error) {
	row, err := tx.Exec("INSERT INTO jobs (workflow, name, status) VALUES (?, ?, ?)", workflow, name, PendingStatus)
	if err != nil {
		log.Println("Error while creating job")
	}

	jobid, _ := row.LastInsertId()
	log.Printf("Created job number %d for for workflow %d\n", jobid, workflow)

	return jobid, err
}

func UpdateStatus(db *sql.DB, id int64, status string) error {
	_, err := db.Exec("UPDATE jobs SET status = ? WHERE id = ?", status, id)
	if err != nil {
		log.Printf("Error while updating status for job %d\n", id)
	}

	return err
}
