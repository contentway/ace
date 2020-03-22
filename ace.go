package ace

import (
	"github.com/contentway/ace/pool"
	"net/http"
	"sync"

	"github.com/julienschmidt/httprouter"
)

var bufPool = pool.NewBufferPool(100)

type Ace struct {
	*Router
	httpRouter   *httprouter.Router
	contextPool  sync.Pool
	render       Renderer
	panicFunc    PanicHandler
	notfoundFunc HandlerFunc
}

type PanicHandler func(c *C, rcv interface{})
type HandlerFunc func(c *C)

func New() *Ace {
	a := &Ace{}
	a.Router = &Router{
		handlers: nil,
		prefix:   "/",
		ace:      a,
	}
	a.panicFunc = defaultPanic
	a.notfoundFunc = defaultNotfound
	a.httpRouter = httprouter.New()
	a.httpRouter.HandleMethodNotAllowed = false

	a.contextPool.New = func() interface{} {
		c := &C{}
		c.handlerIndex = -1
		c.Writer = &c.writercache
		return c
	}

	a.httpRouter.PanicHandler = func(w http.ResponseWriter, req *http.Request, rcv interface{}) {
		c := a.createContext(w, req)
		a.panicFunc(c, rcv)
		a.contextPool.Put(c)
	}

	a.httpRouter.NotFound = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		c := a.createContext(w, req)
		c.handlers = a.Router.handlers
		c.Next()
		a.notfoundFunc(c)
		a.contextPool.Put(c)
	})

	return a
}

// Default server with recovery and logger middleware. Used for Unit Testing.
func Default() *Ace {
	a := New()
	a.Use(Logger())
	return a
}

// ServeHTTP serves a request locally as if it was handled by the HTTP ace
func (a *Ace) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	a.httpRouter.ServeHTTP(w, req)
}

// Run ace with specific address and port
func (a *Ace) Run(addr string) {
	if err := http.ListenAndServe(addr, a); err != nil {
		panic(err)
	}
}

// RunTLS ace with specific address, port and TLS certificates
func (a *Ace) RunTLS(addr string, cert string, key string) {
	if err := http.ListenAndServeTLS(addr, cert, key, a); err != nil {
		panic(err)
	}
}

// SetPoolSize defines the number of write buffers we should keep
func (a *Ace) SetPoolSize(poolSize int) {
	bufPool = pool.NewBufferPool(poolSize)
}
