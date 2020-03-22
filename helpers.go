package ace

import (
	"fmt"
	"log"
	"net/http"
)

// HandleServeFileFallback is a helper method that generates a HTTP handler func. The handler will try to serve
// the requested URL, starting from basePath, and will otherwise serve the file located at fallbackFilePath if
// the requested file wasn't found in the first place. This effectively acts like nginx try_files, or Apache's
// RewriteCond XX -f.
func HandleServeFileFallback(basePath string, fallbackFilePath string) func(c *C) {
	return func(c *C) {
		if FileExists(fmt.Sprintf("%s/%s", basePath, c.Request.URL.Path)) {
			c.File(StatusOK, fmt.Sprintf("%s/%s", basePath, c.Request.URL.Path))
		} else if FileExists(fallbackFilePath) {
			c.File(StatusOK, fallbackFilePath)
		} else {
			log.Printf("Failed to serve fallback file path %s: file does not exist", fallbackFilePath)
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	}
}
