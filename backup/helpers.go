package backup

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

const (
	manifestExt = "manifest"
	archExt     = "tar.gz"
	encExt      = "gpg"
	dateRegexp  = "\\d\\d\\d\\d-\\d\\d-\\d\\d_\\d\\d-\\d\\d-\\d\\d"
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
	bckName := lastPathElement(src)
	now := time.Now().Format("2006-01-02_15-04-05")
	bckName = fmt.Sprintf("%s_%s", bckName, now)
	snapshot := fmt.Sprintf("%s/%s.%s", dst, bckName, manifestExt)
	encArchName := fmt.Sprintf("%s/%s.%s.%s", dst, bckName, archExt, encExt)

	return Artifacts{snapshot, encArchName}
}

func getManifestRegexp(src string) string {
	last := lastPathElement(src)
	pattern := fmt.Sprintf("%s_%s.%s", last, dateRegexp, manifestExt)

	return pattern
}

func GetManifestForSrc(src string, lookup string) (*string, error) {
	files, err := os.ReadDir(lookup)

	if err != nil {
		return nil, err
	}

	re := getManifestRegexp(src)
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		name := file.Name()
		res, err := regexp.MatchString(re, name)

		if err != nil {
			return nil, err
		}

		if res {
			path := fmt.Sprintf("%s/%s", lookup, name)
			return &path, nil
		}
	}

	return nil, nil
}
