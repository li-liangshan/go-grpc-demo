package web

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/labstack/echo"
)

type (
	Stats struct {
		Uptime       time.Time      `json:"uptime"`
		RequestCount int64          `json:"requestCount"`
		Statuses     map[string]int `json:"statuses"`
		mutex        sync.RWMutex
	}
)

func NewStats() *Stats {
	return &Stats{
		Uptime:       time.Now(),
		RequestCount: 0,
		Statuses:     make(map[string]int),
	}
}

// Process is the middleware function.
func (stats *Stats) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println("url:", c.Request().URL.RequestURI())
		if err := next(c); err != nil {
			c.Error(err)
		}
		uri := c.Request().URL.RequestURI()
		if uri == "/stats" {
			return nil
		}
		stats.mutex.Lock()
		defer stats.mutex.Unlock()

		stats.RequestCount++
		status := strconv.Itoa(c.Response().Status)
		stats.Statuses[status]++
		return nil
	}
}

// Handle is the endpoint to get stats.
func (stats *Stats) Handle(c echo.Context) error {
	stats.mutex.Lock()
	defer stats.mutex.Unlock()
	return c.JSON(http.StatusOK, stats)
}

func ServerHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderServer, "Echo/3.0")
		return next(c)
	}
}

func StatsRun() {
	e := echo.New()

	e.Debug = true

	stats := NewStats()
	e.Use(stats.Process)
	e.GET("/stats", stats.Handle)

	e.Use(ServerHeader)
	// Handler
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
