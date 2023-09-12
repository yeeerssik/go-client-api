package main

import (
	"flag"
	"fmt"
	"go_client_service/config"
	"go_client_service/core/handlers"
	"go_client_service/core/helpers"
	"go_client_service/core/middlewares"
	"go_client_service/models"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Server diagnostics
type about struct {
	Version    string
	MinVersion string
	BuildTime  string
	StartedAt  time.Time
	Uptime     string
}

var (
	Version    string
	MinVersion string
	BuildTime  string
)

var (
	host = flag.String("host", "0.0.0.0", "Host ip")
	port = flag.String("port", "8080", "Host port")
)

var serverDetails about

func main() {
	flag.Parse()
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.BodyLimit("1024KB"))
	e.Use(middleware.Secure())
	e.Use(middleware.RemoveTrailingSlash())
	e.Use(middlewares.RequestID)
	e.Use(middlewares.Method)
	e.Use(CORSMiddlewareWrapper)
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `${host} ${remote_ip} &{time_rfc3339_nano} ${id} ${method} ${uri} ${status} "${client_agent}" ${latency} ${bytes_in} ${bytes_out}` + "\n",
	}))

	if err := config.Init("."); err != nil {
		log.Fatal("Got error while initializing config file", err)
	}

	if err := helpers.Init(); err != nil {
		log.Fatal("Got error while initializing helpers", err)
	}
	if err := models.Init(); err != nil {
		log.Fatal("Got error while initializing models", err)
	}

	v := helpers.CustomValidator{Validator: validator.New()}
	v.Init()
	e.Validator = &v

	if err := handlers.Init(); err != nil {
		log.Fatal("Got error while initializing handlers", err)
	}

	e.GET("/", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, Heartbeat())
	})

	e.Any("/api/client", handlers.ClientHandler{}.Any)
	e.Any("/api/client/:client_id", handlers.ClientHandler{}.Any)

	if err := e.Start(fmt.Sprintf("%s:%s", *host, *port)); err != nil {
		fmt.Println("Failed to start server!", err)
		os.Exit(1)
	}
	return
}

func init() {
	serverDetails = about{Version: Version, MinVersion: MinVersion, BuildTime: BuildTime, StartedAt: time.Now()}
}

func Heartbeat() interface{} {
	uptime := time.Since(serverDetails.StartedAt)
	serverDetails.Uptime = fmt.Sprintf("%d days %s", uptime/(time.Hour*24), time.Time{}.Add(uptime).Format("15:04:05"))
	return serverDetails
}

func CORSMiddlewareWrapper(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		req := ctx.Request()
		dynamicCORSConfig := middleware.CORSConfig{
			AllowCredentials: true,
			AllowOrigins:     []string{req.Header.Get("Origin")},
			AllowHeaders:     []string{"Accept", "Cache-Control", "Content-Type", "X-Requested-with"},
		}
		CORSMiddleware := middleware.CORSWithConfig(dynamicCORSConfig)
		CORSHandler := CORSMiddleware(next)
		return CORSHandler(ctx)
	}
}
