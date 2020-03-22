package ace

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleServeFileFallback(t *testing.T) {
	ass := assert.New(t)

	a := Default()
	a.RouteNotFound(HandleServeFileFallback("./", "./README.md"))

	r, _ := http.NewRequest("GET", "/notExisting", nil)
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)
	ass.Equal(StatusOK, w.Code)
	ass.Equal("ACE", w.Body.String()[:3])

	r, _ = http.NewRequest("GET", "/LICENSE", nil)
	w = httptest.NewRecorder()
	a.ServeHTTP(w, r)
	ass.Equal(StatusOK, w.Code)
	ass.Equal("Apache", w.Body.String()[:6])

	// Ensure we get 500 when the fallback file doesn't exist
	a.RouteNotFound(HandleServeFileFallback("./", "./notExistingFile"))
	r, _ = http.NewRequest("GET", "/notExisting", nil)
	w = httptest.NewRecorder()
	a.ServeHTTP(w, r)
	ass.Equal(StatusInternalServerError, w.Code)
}
