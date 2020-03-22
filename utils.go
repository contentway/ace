package ace

import (
	"os"
	"strings"
)

func concat(s ...string) string {
	return strings.Join(s, "")
}

// File exists returns whether or not the filePath exists and is a file
func FileExists(filePath string) bool {
	if fi, err := os.Stat(filePath); err == nil && !fi.IsDir() {
		return true
	} else {
		return false
	}
}
