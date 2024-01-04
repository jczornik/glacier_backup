package tools

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/jczornik/glacier_backup/config"
)

func CreateBackup(src config.BackupSrc, dst config.BackupDst) error {
	log.Println(fmt.Sprintf("Creating backup for %s", src))

	srcLast := lastPathElement(src)
	snapshotFile := fmt.Sprintf("%s/%s.snapshot", dst, srcLast)
	time := time.Now().Format("2006-01-02_15:04:05")
	archive := fmt.Sprintf("%s/%s-%s.tar.gz", dst, srcLast, time)

	if err := exec.Command(Tar, "-czvg", snapshotFile, "-f", archive, src).Run(); err != nil {
		log.Println("Error while creating backup file")
		return err
	}

	log.Println(fmt.Sprintf("Backup file '%s' successfully created", archive))
	return nil
}

func lastPathElement(src string) string {
	tmp := strings.Split(src, "/")
	for i := len(tmp); i > 0; i-- {
		if tmp[i-1] != "" {
			return tmp[i-1]
		}
	}
	return src
}
