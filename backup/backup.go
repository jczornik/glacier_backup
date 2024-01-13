package backup

import (
	"fmt"
	"os"
	"time"

	"github.com/jczornik/glacier_backup/config"
	"github.com/jczornik/glacier_backup/tools"
)

type Artifacts struct {
	snapshot string
	archive  string
}

func CreateEncryptedBackup(src config.BackupSrc, dst config.BackupDst, pass string) (Artifacts, error) {
	bckName := lastPathElement(src)
	now := time.Now().Format("2006-01-02_15-04-05")
	bckName = fmt.Sprintf("%s_%s", bckName, now)
	snapshot := fmt.Sprintf("%s/%s.manifest", dst, bckName)
	encArchName := fmt.Sprintf("%s/%s.tar.gz.gpg", dst, bckName)

	artifacts := Artifacts{snapshot, encArchName}

	bckCmd := tools.NewBackupCmd(src, snapshot)
	encCmd := tools.NewEncryptArmoredStdOutCmd(pass)

	full, err := tools.Pipe(bckCmd, encCmd)
	if err != nil {
		return artifacts, err
	}

	f, err := os.Create(encArchName)
	if err != nil {
		return artifacts, err
	}
	defer f.Close()

	if err := full.SetStdout(f); err != nil {
		return artifacts, err
	}

	return artifacts, full.Run()
}
