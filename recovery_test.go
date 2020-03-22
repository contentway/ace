package ace

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestPanicInHandler assert that panic has been recovered.
func TestPanicInHandler(t *testing.T) {
	// SETUP
	r := New()
	r.GET("/recovery", func(_ *C) {
		panic("Oupps, Houston, we have a problem")
	})

	// RUN
	req, _ := http.NewRequest("GET", "/recovery", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 500 {
		t.Errorf("Response code should be Internal Ace Error, was: %d", w.Code)
	}
}

// TestPanicWithAbort assert that panic has been recovered even if context.Abort was used.
func TestPanicWithAbort(t *testing.T) {
	// SETUP
	r := New()
	r.GET("/recovery", func(c *C) {
		c.AbortWithStatus(500)
		panic("Oupps, Houston, we have a problem")
	})

	// RUN
	req, _ := http.NewRequest("GET", "/recovery", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// TEST
	if w.Code != 500 {
		t.Errorf("Response code should be Bad request, was: %v", w.Code)
	}
}
