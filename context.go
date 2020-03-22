package ace

import (
	"encoding/json"
	"fmt"
	"github.com/contentway/ace/sessions"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type C struct {
	Request      *http.Request
	writercache  responseWriter
	Writer       ResponseWriter
	handlerIndex int8
	handlers     []HandlerFunc
	params       httprouter.Params
	data         map[string]interface{}
	sessions     *sessions.Sessions
	render       Renderer
}

func (a *Ace) createContext(w http.ResponseWriter, r *http.Request) *C {
	c := a.contextPool.Get().(*C)
	c.writercache.reset(w)
	c.Request = r
	c.handlerIndex = -1
	c.data = nil
	c.render = a.render

	return c
}

// Next runs the next HandlerFunc in the stack (ie. the next middleware)
func (c *C) Next() {
	c.handlerIndex++
	s := int8(len(c.handlers))
	if c.handlerIndex < s {
		c.handlers[c.handlerIndex](c)
	}
}

// Abort stops the middleware chain
func (c *C) Abort() {
	c.handlerIndex = AbortMiddlewareIndex
}

// AbortWithStatus stops the middleware chain and return the specified HTTP status code
func (c *C) AbortWithStatus(status int) {
	c.Writer.WriteHeader(status)
	c.Abort()
}

// AddHeader adds an header to the response
func (c *C) AddHeader(key string, value string) {
	c.Writer.Header().Add(key, value)
}

// SetHeader replaces and sets an header to the response
func (c *C) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// String response with text/html; charset=UTF-8 Content type
func (c *C) String(status int, format string, val ...interface{}) {
	c.Writer.Header().Set(HeaderContentType, "text/html; charset=UTF-8")
	c.Writer.WriteHeader(status)

	buf := bufPool.Get()
	defer bufPool.Put(buf)

	if len(val) == 0 {
		buf.WriteString(format)
	} else {
		buf.WriteString(fmt.Sprintf(format, val...))
	}

	c.Writer.Write(buf.Bytes())
}

// JSON response with application/json; charset=UTF-8 Content type
func (c *C) JSON(status int, v interface{}) {
	c.Writer.Header().Set(HeaderContentType, "application/json; charset=UTF-8")
	c.Writer.WriteHeader(status)
	if v == nil {
		return
	}

	buf := bufPool.Get()
	defer bufPool.Put(buf)

	if err := json.NewEncoder(buf).Encode(v); err != nil {
		panic(err)
	}

	c.Writer.Write(buf.Bytes())
}

// File responds a file with text/html content type from a specified file path
func (c *C) File(status int, filePath string) {
	bytes, err := ioutil.ReadFile(filePath)

	if err != nil {
		log.Printf("Could not read file: %s", filePath)
		c.AbortWithStatus(500)
		return
	}

	// SVG are incorrectly detected as text/xml (well technically they are), but browsers require image/svg+xml
	// to properly display SVGs. We'll assume that ".svg" files systematically of type image/svg+xml. Same happens
	// with JavaScript files, they should be transferred as text/javascript
	if strings.HasSuffix(filePath, ".svg") {
		c.Writer.Header().Set(HeaderContentType, "image/svg+xml")
	} else if strings.HasSuffix(filePath, ".js") {
		c.Writer.Header().Set(HeaderContentType, "text/javascript")
	} else if strings.HasSuffix(filePath, ".css") {
		c.Writer.Header().Set(HeaderContentType, "text/css")
	} else {
		c.Writer.Header().Set(HeaderContentType, http.DetectContentType(bytes))
	}
	c.Writer.WriteHeader(status)
	c.Writer.Write(bytes)
}

// Param returns a parameter value from the route
func (c *C) Param(name string) string {
	return c.params.ByName(name)
}

// ParseJSON decodes json to interface{}
func (c *C) ParseJSON(v interface{}) error {
	defer c.Request.Body.Close()
	return json.NewDecoder(c.Request.Body).Decode(v)
}

// Download response with application/octet-stream; charset=UTF-8 Content type
func (c *C) Download(status int, filename string, v []byte) {
	c.Writer.Header().Set(HeaderContentType, "application/octet-stream; charset=UTF-8")
	c.Writer.Header().Set(HeaderContentDisposition, fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Writer.WriteHeader(status)
	c.Writer.Write(v)
}

// DownloadStream response with application/octet-stream, but reading from a io.Reader
func (c *C) DownloadStream(status int, filename string, v io.Reader) {
	c.Writer.Header().Set(HeaderContentType, "application/octet-stream; charset=UTF-8")
	c.Writer.Header().Set(HeaderContentDisposition, fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Writer.WriteHeader(status)
	io.Copy(c.Writer, v)
}

// ClientRemoteAddress returns the remote IP address and port
func (c *C) ClientRemoteAddress() string {
	return c.Request.RemoteAddr
}

// ClientIP returns the remote IP address, without port. See ClientRemoteAddress for source IP and port.
func (c *C) ClientIP() string {
	ra := c.Request.RemoteAddr
	return ra[:strings.LastIndex(ra, ":")]
}

// Set stores arbitrary data inside this context to re-use in handlers
func (c *C) Set(key string, v interface{}) {
	if c.data == nil {
		c.data = make(map[string]interface{})
	}
	c.data[key] = v
}

// SetAll stores and erases the existing data inside this context. See C.Set().
func (c *C) SetAll(data map[string]interface{}) {
	c.data = data
}

// Get retrieves data previously stored in this context using C.Set() or C.SetAll().
func (c *C) Get(key string) interface{} {
	return c.data[key]
}

// GetAll returns all the data previously stored in this context using C.Set() or C.SetAll().
func (c *C) GetAll() map[string]interface{} {
	return c.data
}

// Sessions returns the sessions with the provided name
func (c *C) Sessions(name string) *sessions.Session {
	return c.sessions.Get(name)
}

// QueryString returns a param value from the URL GET parameters, where "key" is the parameter key, and "d" is
// the default value when the key isn't set in the current request.
func (c *C) QueryString(key, d string) string {
	val := c.Request.URL.Query().Get(key)
	if val == "" {
		return d
	}

	return val
}

// QueryString returns a param array from the URL GET parameters, where "key" is the parameter key, and "d" is
// the default value when the key isn't set in the current request.
func (c *C) QueryStringArray(key string, d []string) []string {
	val := c.Request.URL.Query()[key]
	if val == nil {
		return d
	}

	return val
}

// QueryStringInteger returns a param value from the URL GET parameters, where "key" is the parameter key, and "d" is
// the default value when the key isn't set in the current request, or if the value is not a valid integer..
func (c *C) QueryStringInteger(key string, d int64) int64 {
	val := c.Request.URL.Query().Get(key)
	if val == "" {
		return d
	}

	valInt, err := strconv.ParseInt(val, 10, 64)

	if err != nil {
		log.Printf("Failed to parse QueryString key '%s' value '%s' to int: %s", key, val, err)
		return d
	}

	return valInt
}

func (c *C) FormData() (url.Values, error) {
	var err error
	if c.Request.PostForm == nil {
		err = c.Request.ParseForm()
	}

	return c.Request.PostForm, err
}

func (c *C) MultipartFormData() (*multipart.Form, error) {
	var err error
	if c.Request.MultipartForm == nil {
		err = c.Request.ParseMultipartForm(51200000)
	}

	return c.Request.MultipartForm, err
}

func (c *C) FormFile(key string) (multipart.File, *multipart.FileHeader, error) {
	return c.Request.FormFile(key)
}
