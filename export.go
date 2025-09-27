package gojs

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/raoqu/gojs/internal/api"
	"github.com/raoqu/gojs/internal/server"
	"github.com/raoqu/gojs/internal/util"
	"github.com/raoqu/gojs/internal/util/mysql"
)

type GoJSInstance struct {
	pool          *api.ScriptPool
	ScriptManager *server.ScriptManager
	Config        *Config
	Server        *http.Server
}

func CreateInstance(config *Config) *GoJSInstance {
	return &GoJSInstance{
		Config: config,
	}
}

func (js *GoJSInstance) Run() (*server.ScriptManager, error) {
	config := js.Config
	// 是否依赖 Redis
	if config.Redis.Addr != "" {
		if err := util.InitRedisClient(config.Redis.Addr, config.Redis.DB, config.Redis.DBConfig, config.Redis.Password); err != nil {
			log.Printf("[gojs] Failed to initialize Redis: %v", err)
			return nil, err
		}
	}
	// 是否依赖 MySQL
	if config.MySQL.Host != "" {
		if err := mysql.InitializeMySQL(config.MySQL.Host, config.MySQL.Port, config.MySQL.DB, config.MySQL.User, config.MySQL.Password, config.MySQL.Timeout); err != nil {
			log.Printf("[gojs] Failed to initialize MySQL: %v", err)
			return nil, err
		}
	}
	pool := server.InitScriptPool(config.GOJS.RedisKey)
	js.ScriptManager = server.NewScriptManager(pool, config.GOJS.Endpoint)
	js.pool = pool
	if config.GOJS.Server {
		js.startServer(js.ScriptManager)
	}
	return js.ScriptManager, nil
}

func (js *GoJSInstance) startServer(manager *server.ScriptManager) {
	port := js.Config.GOJS.Port
	webDir := js.Config.GOJS.WebDir
	endpoint := js.Config.GOJS.Endpoint
	title := js.Config.GOJS.Title

	log.Printf("[gojs] Starting web server on port %d...", port)

	// Create Gin router
	router := gin.Default()
	// Setup HTML template rendering
	router.LoadHTMLGlob(webDir + "/*.html")

	router.Static("/static", webDir)
	router.GET("/favicon.ico", func(c *gin.Context) {
		c.File(webDir + "/favicon.ico")
	})

	if endpoint != "" {
		server.SetupScriptsRoutes(router, manager)
	}

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"ScriptEndpoint": endpoint,
			"AppTitle":       title,
		})
	})

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	// Start HTTP server in a goroutine
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("[gojs] HTTP server error: %v", err)
		}
	}()

	log.Printf("GoJS web server started: http://127.0.0.1:%d", port)

	js.Server = httpServer
}

func (js *GoJSInstance) StopServer() {
	if js.Server != nil {
		if err := js.Server.Shutdown(context.Background()); err != nil {
			log.Printf("[gojs] HTTP server shutdown error: %v", err)
		}
	}
}

func (js *GoJSInstance) Wait() {
	if js.Server == nil {
		return
	}

	// Wait for interrupt/terminate signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Printf("[gojs] Shutting down web server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := js.Server.Shutdown(ctx); err != nil {
		log.Printf("[gojs] HTTP server shutdown error: %v", err)
	}
}

func (js *GoJSInstance) Update(scriptName string, script string) {
	js.ScriptManager.Store(scriptName, script)
}

func (js *GoJSInstance) Execute(scriptName string, params map[string]interface{}) (interface{}, error) {
	return js.ScriptManager.Call(scriptName, params)
}
