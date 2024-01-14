package tasks

import (
	"log"
	"os"
)

type PreserveTask struct {
	original string
	copy     string
}

func NewPreserveTask(manifest string) PreserveTask {
	copy := manifest + ".old"
	return PreserveTask{manifest, copy}
}

func moveManifest(src string, dst string) error {
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return nil
	}

	return os.Rename(src, dst)
}

func (t PreserveTask) Exec() error {
	log.Printf("Starting presering old manifest %s\n", t.original)
	err := moveManifest(t.original, t.copy)

	if err != nil {
		log.Printf("Error while preserving manifest %s", t.original)
	} else {
		log.Printf("Successfully preserved manifest %s.\nSaved to %s\n", t.original, t.copy)
	}

	return err
}

func (t PreserveTask) Rollback() error {
	log.Printf("Rollback for preserve %s\n", t.original)
	err := moveManifest(t.copy, t.original)

	if err != nil {
		log.Printf("Error while rollbacking preserve for %s\n", t.original)
	} else {
		log.Printf("Successfully rollbacked preserve for %s\n", t.original)
	}

	return err
}
