package tasks

import (
	"log"
	"os"

	"github.com/jczornik/glacier_backup/persistent"
	"github.com/jczornik/glacier_backup/persistent/artifacts"
)

const saveArtifactsName = "SaveArtifacts"

type SaveArtifactsTask struct {
	workflowId      int64
	manifest        string
	archive         string
	db              persistent.DBClient
	savedArtifactId *int64
}

func NewSaveArtifactsTask(manifest string, archive string, dbClient persistent.DBClient, workflowId int64) *SaveArtifactsTask {
	return &SaveArtifactsTask{workflowId, manifest, archive, dbClient, nil}
}

func (t *SaveArtifactsTask) Name() string {
	return saveArtifactsName
}

func (t *SaveArtifactsTask) Exec() error {
	log.Printf("Saving artifacts for workflow %d\n", t.workflowId)

	if _, err := os.Stat(t.archive); os.IsNotExist(err) {
		log.Printf("Archive %s does not exist\n", t.archive)
		return err
	}

	mf, err := os.Open(t.manifest)
	if err != nil {
		log.Printf("Error while opening manifest %s\n", t.manifest)
		return err
	}
	defer mf.Close()

	mstats, err := mf.Stat()
	if err != nil {
		log.Printf("Error while getting stats for manifest %s\n", t.manifest)
		return err
	}

	mcontent := make([]byte, mstats.Size())
	_, err = mf.Read(mcontent)

	if err != nil {
		log.Printf("Error while reading manifest %s\n", t.manifest)
		return err
	}

	db, err := t.db.OpenDB()
	if err != nil {
		log.Printf("Error while opening db for saving artifacts for workflow %d\n", t.workflowId)
		return err
	}

	id, err := artifacts.Create(db, t.archive, t.manifest, mcontent, t.workflowId)
	if err != nil {
		log.Printf("Error while saving artifact into db for workflow %d\n", t.workflowId)
		return err
	}

	t.savedArtifactId = &id

	log.Printf("Successfully saved artifacts for workflow %d\n", t.workflowId)
	return nil
}

func (t *SaveArtifactsTask) Rollback() error {
	if t.savedArtifactId == nil {
		log.Println("No artifact saved - nothing to rollback")
		return nil
	}

	log.Println("Rollback for saving artifacts")
	db, err := t.db.OpenDB()
	if err != nil {
		log.Println("Error while opening db for rollback")
	}

	_, err = db.Exec("DELETE FROM workflow_artifacts WHERE id = ?", *t.savedArtifactId)
	if err != nil {
		log.Println("Error while removing artifact from db")
		return err
	}

	return nil
}
