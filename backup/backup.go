package backup

import (
	"os"

	"github.com/jczornik/glacier_backup/config"
	"github.com/jczornik/glacier_backup/tools"
)

type Artifacts struct {
	Snapshot string
	Archive  string
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

	return full.Run()
}
