package backup

import (
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	manifestExt = "manifest"
	archExt     = "tar.gz"
	encExt      = "gpg"
)

func lastPathElement(src string) string {
	tmp := strings.Split(src, "/")
	for i := len(tmp); i > 0; i-- {
		if tmp[i-1] != "" {
			return tmp[i-1]
		}
	}
	return src
}

func NewArtifactNames(src string, dst string) Artifacts {
	last := lastPathElement(src)
	now := time.Now().Format("2006-01-02_15-04-05")
	bckName := fmt.Sprintf("%s_%s", last, now)
	manifest := MakeManifestName(src, dst)
	encArchName := fmt.Sprintf("%s/%s.%s.%s", dst, bckName, archExt, encExt)

	return Artifacts{manifest, encArchName}
}

func MakeManifestName(src string, dst string) string {
	last := lastPathElement(src)
	return fmt.Sprintf("%s/%s.%s", dst, last, manifestExt)
}

func GetManifestForSrc(src string, lookup string) (*string, error) {
	bckName := lastPathElement(src)
	expected := fmt.Sprintf("%s/%s.%s", lookup, bckName, manifestExt)

	if _, err := os.Stat(expected); err != nil && !os.IsNotExist(err) {
		return nil, err
	} else if os.IsNotExist(err) {
		return nil, nil
	} else {
		return &expected, nil
	}
}
