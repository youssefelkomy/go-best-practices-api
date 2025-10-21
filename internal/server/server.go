package server

import (
	_ "embed"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	"my-go-server/pkg/handler"

	"github.com/gin-gonic/gin"
)

//go:embed openapi.yaml
var openapiYAML []byte

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

	// serve OpenAPI spec file and docs
	s.engine.GET("/openapi.yaml", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/x-yaml", openapiYAML)
	})

	// Redoc UI for interactive documentation
	s.engine.GET("/docs", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		// quote the YAML so it can be embedded safely into a JS string literal
		specLiteral := strconv.Quote(string(openapiYAML))

		html := `<!doctype html>
<html>
	<head>
		<title>API Docs</title>
		<meta charset="utf-8" />
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<script src="https://cdnjs.cloudflare.com/ajax/libs/js-yaml/4.1.0/js-yaml.min.js"></script>
		<script src="https://cdn.redoc.ly/redoc/latest/bundles/redoc.standalone.js"></script>
		<style>body{margin:0;padding:0}</style>
	</head>
	<body>
		<div id="redoc"></div>
		<script>
			(function(){
				try{
					const specYAML = ` + specLiteral + `;
					const spec = jsyaml.load(specYAML);
					Redoc.init(spec, {}, document.getElementById('redoc'))
				}catch(e){
					document.getElementById('redoc').innerText = 'Failed to render docs: '+e
				}
			})();
		</script>
		<noscript>
			<p>JavaScript required. You can view the raw spec <a href="/openapi.yaml">here</a>.</p>
		</noscript>
	</body>
</html>`

		c.String(http.StatusOK, html)
	})
}
