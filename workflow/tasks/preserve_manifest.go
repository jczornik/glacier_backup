package tasks

import "os"

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
	return moveManifest(t.original, t.copy)
}

func (t PreserveTask) Rollback() error {
	return moveManifest(t.copy, t.original)
}
