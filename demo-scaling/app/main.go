package main

import (
	"fmt"
	"log"
	"os"
	"sync/atomic"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

var (
	requestCount uint64
	instanceID   string
	startTime    time.Time
)

func init() {
	// Get hostname as instance ID (Docker container ID)
	hostname, err := os.Hostname()
	if err != nil {
		instanceID = "unknown"
	} else {
		instanceID = hostname
	}
	startTime = time.Now()
}

func main() {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	// Logger middleware
	app.Use(logger.New(logger.Config{
		Format: fmt.Sprintf("[${time}] ${status} ${method} ${path} - ${latency} (instance: %s)\n", instanceID[:min(12, len(instanceID))]),
	}))

	// Health check endpoint
	app.Get("/health", healthHandler)

	// CPU-intensive work endpoint
	app.Get("/api/work", workHandler)

	// Metrics endpoint
	app.Get("/api/metrics", metricsHandler)

	// Root endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"service":  "scaling-demo",
			"instance": instanceID,
			"message":  "JIF USK x Twibbonize Workshop",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Instance %s starting on port %s", instanceID[:min(12, len(instanceID))], port)
	log.Fatal(app.Listen(":" + port))
}

func healthHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":   "ok",
		"instance": instanceID,
		"uptime":   time.Since(startTime).String(),
	})
}

func workHandler(c *fiber.Ctx) error {
	atomic.AddUint64(&requestCount, 1)

	// Simulate CPU-intensive work (Fibonacci calculation)
	start := time.Now()
	result := fibonacci(35) // ~50-100ms of CPU work
	duration := time.Since(start)

	return c.JSON(fiber.Map{
		"instance":    instanceID,
		"result":      result,
		"duration_ms": duration.Milliseconds(),
		"request_num": atomic.LoadUint64(&requestCount),
	})
}

func metricsHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"instance":      instanceID,
		"request_count": atomic.LoadUint64(&requestCount),
		"uptime":        time.Since(startTime).String(),
		"uptime_secs":   time.Since(startTime).Seconds(),
	})
}

// fibonacci calculates nth fibonacci number (CPU intensive)
func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
