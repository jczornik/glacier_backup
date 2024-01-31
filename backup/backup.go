package backup

import (
	"fmt"
	"os"
	"strings"

	"github.com/jczornik/glacier_backup/config"
	"github.com/jczornik/glacier_backup/tools"
)

type Artifacts struct {
	Snapshot string
	Archive  string
}

type EncryptedBackupError struct {
	err    string
	stderr strings.Builder
}

func (e EncryptedBackupError) Error() string {
	return fmt.Sprintf("Error while creating encrypted backup: %s. Stderr: %s\n", e.err, e.stderr.String())
}

func CreateEncryptedBackup(src config.BackupSrc, artifacst Artifacts, pass string) error {
	bckCmd := tools.NewBackupCmd(src, artifacst.Snapshot)
	encCmd := tools.NewEncryptArmoredStdOutCmd(pass)

	full, err := tools.Pipe(bckCmd, encCmd)
	if err != nil {
		return err
	}

	f, err := os.Create(artifacst.Archive)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := full.SetStdout(f); err != nil {
		return err
	}

	backupErr := EncryptedBackupError{err: "", stderr: strings.Builder{}}
	if err := full.SetStderr(&backupErr.stderr); err != nil {
		return err
	}

	if err := full.Run(); err != nil {
		backupErr.err = err.Error()
		return backupErr
	}

	return nil
}
