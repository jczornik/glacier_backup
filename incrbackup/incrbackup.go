package incrbackup

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/jczornik/glacier_backup/config"
	"github.com/jczornik/glacier_backup/tools"
)

func CreateBackup(src config.BackupSrc, dst config.BackupDst) error {
	srcLast := lastPathElement(src)
	snapshotFile := fmt.Sprintf("%s/%s", dst, srcLast)
	archive := fmt.Sprintf("%s/%s.tar.gz", dst, srcLast)

	return exec.Command(tools.Tar, "-czvg", snapshotFile, "-f", archive, src).Run()
}

func lastPathElement(src string) string {
	tmp := strings.Split(src, "/")
	for i := len(tmp); i > 0; i-- {
		if tmp[i - 1] != "" {
			return tmp[i - 1]
		}
	}
	return src
}
