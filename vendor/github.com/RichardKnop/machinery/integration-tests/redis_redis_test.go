package integrationtests

import (
	"fmt"
	"os"
	"testing"

	"github.com/RichardKnop/machinery/v1/config"
)

func TestRedisRedis(t *testing.T) {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		return
	}

	// Redis broker, Redis result backend
	server := setup(&config.Config{
		Broker:        fmt.Sprintf("redis://%v", redisURL),
		DefaultQueue:  "test_queue",
		ResultBackend: fmt.Sprintf("redis://%v", redisURL),
	})
	worker := server.NewWorker("test_worker")
	go worker.Launch()
	testAll(server, t)
	worker.Quit()
}
