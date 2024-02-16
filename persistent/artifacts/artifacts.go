package artifacts

import (
	"database/sql"
	"log"
)

const createTableQuery = `
    CREATE TABLE workflow_artifacts
    (
	id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        archive_name TEXT NOT NULL,
        manifest_name TEXT NOT NULL,
        manifest_content BLOB NOT NULL,
	workflow INTEGER NOT NULL,
	FOREIGN KEY (workflow) REFERENCES workflows(id)
    )`

func CreateTable(db *sql.Tx) error {
	_, err := db.Exec(createTableQuery)
	if err != nil {
		log.Println("Error while creating workflow_artifacts table")
	}

	return err
}

func Create(db *sql.DB, archive string, manifest string, mcontent []byte, workflow int64) (int64, error) {
	var artifactid int64

	row, err := db.Exec("INSERT INTO workflow_artifacts (archive_name, manifest_name, manifest_content, workflow) VALUES (?, ?, ?, ?)", archive, manifest, mcontent, workflow)
	if err != nil {
		log.Println("Error while creating artifact")
		return artifactid, err
	}

	artifactid, _ = row.LastInsertId()
	log.Printf("Created artifact number %d for workflow %d\n", artifactid, workflow)

	return artifactid, err
}

func Delete(db *sql.DB, id int64) error {
	_, err := db.Exec("DELETE FROM workflow_artifact WHERE id = ?", id)
	return err
}
