package backup

import (
	"fmt"
	"os"
	"time"

	"github.com/jczornik/glacier_backup/config"
	"github.com/jczornik/glacier_backup/tools"
)

func CreateEncryptedBackup(src config.BackupSrc, dst config.BackupDst, pass string) error {
	bckName := lastPathElement(src)
	now := time.Now().Format("2006-01-02_15-04-05")
	bckName = fmt.Sprintf("%s_%s", bckName, now)
	snapshot := fmt.Sprintf("%s/%s.manifest", dst, bckName)
	encArchName := fmt.Sprintf("%s/%s.tar.gz.gpg", dst, bckName)

	bckCmd := tools.NewBackupCmd(src, snapshot)
	encCmd := tools.NewEncryptArmoredStdOutCmd(pass)

	full, err := tools.Pipe(bckCmd, encCmd)
	if err != nil {
		return err
	}

	f, err := os.Create(encArchName)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := full.SetStdout(f); err != nil {
		return err
	}

	return full.Run()
}
