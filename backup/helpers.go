package backup

import "strings"

func lastPathElement(src string) string {
	tmp := strings.Split(src, "/")
	for i := len(tmp); i > 0; i-- {
		if tmp[i-1] != "" {
			return tmp[i-1]
		}
	}
	return src
}
