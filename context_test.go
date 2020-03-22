package ace

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJSONResp(t *testing.T) {
	ass := assert.New(t)

	data := map[string]interface{}{
		"s": "test",
		"n": 123,
		"b": true,
	}

	a := New()
	a.GET("/", func(c *C) {
		c.JSON(200, data)
	})

	buf := &bytes.Buffer{}
	json.NewEncoder(buf).Encode(data)

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)
	ass.Equal(200, w.Code)
	ass.Equal(buf.String(), w.Body.String())
	ass.Equal("application/json; charset=UTF-8", w.Header().Get("Content-Type"))
}

func TestStringResp(t *testing.T) {
	ass := assert.New(t)
	a := New()
	a.GET("/", func(c *C) {
		c.String(200, "123")
	})

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)
	ass.Equal(200, w.Code)
	ass.Equal("123", w.Body.String())
	ass.Equal("text/html; charset=UTF-8", w.Header().Get("Content-Type"))
}

func TestDownloadResp(t *testing.T) {
	ass := assert.New(t)
	a := New()
	a.GET("/", func(c *C) {
		c.Download(200, "test.txt", []byte("123"))
	})

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)
	ass.Equal(200, w.Code)
	ass.Equal("123", w.Body.String())
	ass.Equal("application/octet-stream; charset=UTF-8", w.Header().Get(HeaderContentType))
	ass.Equal("attachment; filename=\"test.txt\"", w.Header().Get(HeaderContentDisposition))
}

func TestCData(t *testing.T) {
	ass := assert.New(t)
	a := New()

	a.Use(func(c *C) {
		c.Set("test", "123")
		c.Next()
	})

	a.GET("/", func(c *C) {
		c.GetAll()
		c.String(200, c.Get("test").(string))
	})

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)
	ass.Equal(200, w.Code)
	ass.Equal("123", w.Body.String())
}
