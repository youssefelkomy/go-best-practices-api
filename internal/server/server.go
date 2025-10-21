package server

import (
	"net/http"
	"sync/atomic"
	"time"

	"my-go-server/pkg/handler"

	"github.com/gin-gonic/gin"
)

// Server uses a gin.Engine internally but exposes ServeHTTP so tests can use
// httptest with the server instance.
type Server struct {
	addr   string
	engine *gin.Engine
	start  time.Time
	reqCnt uint64
}

// New creates a Server listening on the default address ":8080".
func New() *Server {
	return NewWithAddr(":8080")
}

// NewWithAddr creates a Server that will listen on the provided address.
func NewWithAddr(addr string) *Server {
	g := gin.New()
	g.Use(gin.Logger())
	g.Use(gin.Recovery())

	s := &Server{
		addr:   addr,
		engine: g,
		start:  time.Now(),
	}

	// request counting middleware
	g.Use(func(c *gin.Context) {
		atomic.AddUint64(&s.reqCnt, 1)
		c.Next()
	})

	s.routes()
	return s
}

// NewServer is kept for compatibility with older code that expects
// a constructor named NewServer(addr string).
func NewServer(addr string) *Server {
	return NewWithAddr(addr)
}

// ServeHTTP makes Server implement http.Handler by delegating to the gin engine.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.engine.ServeHTTP(w, r)
}

// Start runs the HTTP server on the configured address.
func (s *Server) Start() error {
	return s.engine.Run(s.addr)
}

func (s *Server) routes() {
	// Wrap existing net/http handlers so we reuse pkg/handler implementations.
	// root: service info
	s.engine.GET("/", func(c *gin.Context) {
		uptime := time.Since(s.start).String()
		c.JSON(http.StatusOK, gin.H{
			"service":  "my-go-server",
			"version":  "0.1.0",
			"uptime":   uptime,
			"requests": atomic.LoadUint64(&s.reqCnt),
		})
	})

	// additional informative endpoints
	s.engine.GET("/info", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"name":    "my-go-server",
			"version": "0.1.0",
			"started": s.start.Format(time.RFC3339),
		})
	})

	s.engine.GET("/time", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"time": time.Now().Format(time.RFC3339Nano)})
	})

	s.engine.GET("/headers", func(c *gin.Context) {
		headers := map[string][]string{}
		for k, v := range c.Request.Header {
			headers[k] = v
		}
		c.JSON(http.StatusOK, gin.H{"headers": headers})
	})

	s.engine.Any("/echo", func(c *gin.Context) {
		// echo back method, headers, query and body (as raw)
		body, _ := c.GetRawData()
		c.JSON(http.StatusOK, gin.H{
			"method":  c.Request.Method,
			"query":   c.Request.URL.RawQuery,
			"headers": c.Request.Header,
			"body":    string(body),
		})
	})

	// metrics endpoint
	s.engine.GET("/metrics", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"uptime":   time.Since(s.start).String(),
			"requests": atomic.LoadUint64(&s.reqCnt),
		})
	})

	// keep compatibility routes
	s.engine.GET("/health", gin.WrapF(handler.HealthCheckHandler))
	s.engine.GET("/hello", gin.WrapF(handler.HelloHandler))
}
