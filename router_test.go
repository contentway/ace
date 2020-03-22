package ace

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var testHandler = func(c *C) { c.Next() }

func TestHTTPMethod(t *testing.T) {
	ass := assert.New(t)

	a := Default()
	a.GET("/test", func(c *C) {
		c.String(200, "Test")
	})

	a.POST("/test", func(c *C) {
		c.String(200, c.Request.FormValue("test"))
	})

	a.PUT("/", func(c *C) {
		c.String(200, c.Request.FormValue("test"))
	})

	a.PATCH("/", func(c *C) {
		c.String(200, c.Request.FormValue("test"))
	})

	a.DELETE("/", func(c *C) {
		c.String(200, "deleted")
	})

	a.OPTIONS("/", func(c *C) {
		c.String(200, "options")
	})

	a.HEAD("/test", func(c *C) {
		c.String(200, "head")
	})

	r, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)
	ass.Equal(200, w.Code)
	ass.Equal("Test", w.Body.String())

	r, _ = http.NewRequest("POST", "/test", nil)
	r.ParseForm()
	r.Form.Add("test", "hello")
	w = httptest.NewRecorder()
	a.ServeHTTP(w, r)
	ass.Equal(200, w.Code)
	ass.Equal("hello", w.Body.String())

	r, _ = http.NewRequest("PUT", "/", nil)
	r.ParseForm()
	r.Form.Add("test", "hello")
	w = httptest.NewRecorder()
	a.ServeHTTP(w, r)
	ass.Equal(200, w.Code)
	ass.Equal("hello", w.Body.String())

	r, _ = http.NewRequest("PATCH", "/", nil)
	r.ParseForm()
	r.Form.Add("test", "hello")
	w = httptest.NewRecorder()
	a.ServeHTTP(w, r)
	ass.Equal(200, w.Code)
	ass.Equal("hello", w.Body.String())

	r, _ = http.NewRequest("DELETE", "/", nil)
	w = httptest.NewRecorder()
	a.ServeHTTP(w, r)
	ass.Equal(200, w.Code)
	ass.Equal("deleted", w.Body.String())

	r, _ = http.NewRequest("OPTIONS", "/", nil)
	w = httptest.NewRecorder()
	a.ServeHTTP(w, r)
	ass.Equal(200, w.Code)
	ass.Equal("options", w.Body.String())

	r, _ = http.NewRequest("HEAD", "/test", nil)
	w = httptest.NewRecorder()
	a.ServeHTTP(w, r)
	ass.Equal(200, w.Code)
	ass.Equal("head", w.Body.String())

	// trailing slash
	r, _ = http.NewRequest("GET", "/test/", nil)
	w = httptest.NewRecorder()
	a.ServeHTTP(w, r)
	ass.Equal(301, w.Code)

	r, _ = http.NewRequest("POST", "/test/", nil)
	w = httptest.NewRecorder()
	a.ServeHTTP(w, r)
	ass.Equal(307, w.Code)
}

func TestNestedGroupRoute(t *testing.T) {
	ass := assert.New(t)

	a := Default()
	g1 := a.Group("/g1", testHandler)
	g2 := g1.Group("/g2", testHandler)
	g3 := g2.Group("/g3", testHandler)

	g3.GET("/", func(c *C) {
		c.String(200, "g3")
	})

	g3.GET("/test", func(c *C) {
		c.String(200, "g3/test")
	})

	r, _ := http.NewRequest("GET", "/g1/g2/g3/", nil)
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)
	ass.Equal(200, w.Code)
	ass.Equal("g3", w.Body.String())

	r, _ = http.NewRequest("GET", "/g1/g2/g3/test", nil)
	w = httptest.NewRecorder()
	a.ServeHTTP(w, r)
	ass.Equal(200, w.Code)
	ass.Equal("g3/test", w.Body.String())
}

func TestGroupRoute(t *testing.T) {
	ass := assert.New(t)

	a := Default()
	g1 := a.Group("/g1", testHandler)
	g2 := a.Group("/g2", testHandler)

	g1.GET("/", func(c *C) {
		c.String(200, "g1")
	})

	g1.GET("/test", func(c *C) {
		c.String(200, "g1/test")
	})

	g2.POST("/", func(c *C) {
		c.String(200, "g2")
	})

	g2.POST("/test", func(c *C) {
		c.String(200, "g2/test")
	})

	r, _ := http.NewRequest("GET", "/g1/", nil)
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)
	ass.Equal(200, w.Code)
	ass.Equal("g1", w.Body.String())

	r, _ = http.NewRequest("GET", "/g1/test", nil)
	w = httptest.NewRecorder()
	a.ServeHTTP(w, r)
	ass.Equal(200, w.Code)
	ass.Equal("g1/test", w.Body.String())

	r, _ = http.NewRequest("POST", "/g2/", nil)
	w = httptest.NewRecorder()
	a.ServeHTTP(w, r)
	ass.Equal(200, w.Code)
	ass.Equal("g2", w.Body.String())

	r, _ = http.NewRequest("POST", "/g2/test", nil)
	w = httptest.NewRecorder()
	a.ServeHTTP(w, r)
	ass.Equal(200, w.Code)
	ass.Equal("g2/test", w.Body.String())
}

func TestServeStatic(t *testing.T) {
	ass := assert.New(t)

	os.MkdirAll("/tmp/dims-autotest", 0755)
	os.Create("/tmp/dims-autotest/segv_output.ZBG0OQ")

	a := Default()
	a.Static("/assets", "/tmp/dims-autotest", testHandler)

	r, _ := http.NewRequest("GET", "/assets/segv_output.ZBG0OQ", nil)
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)
	ass.Equal(200, w.Code)

	r, _ = http.NewRequest("GET", "/assets/test.text", nil)
	w = httptest.NewRecorder()
	a.ServeHTTP(w, r)
	ass.Equal(404, w.Code)
}

func TestConvertHandlerFunc(t *testing.T) {
	ass := assert.New(t)

	a := Default()
	a.GET("/", a.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("test"))
	}))

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)
	ass.Equal(200, w.Code)
	ass.Equal("test", w.Body.String())
}

func TestRouteNotFound(t *testing.T) {
	ass := assert.New(t)

	a := Default()
	a.RouteNotFound(func(c *C) {
		c.String(404, "test not found")
	})

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)
	ass.Equal(404, w.Code)
	ass.Equal("test not found", w.Body.String())
}

func TestStaticPath(t *testing.T) {
	ass := assert.New(t)

	a := New()
	path := a.Router.staticPath("/")
	ass.Equal("/*filepath", path)

	path = a.Router.staticPath("/public")
	ass.Equal("/public/*filepath", path)
}
