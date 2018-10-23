package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/antihax/evedata/internal/redigohelper"
	"github.com/antihax/evedata/services/zkillboard"
)

func main() {
	redis := redigohelper.ConnectRedisPool(
		[]string{"sentinel1:26379", "sentinel2:26379", "sentinel3:26379"},
		os.Getenv("REDIS_PASSWORD"),
		"evedata",
		true,
	)

	// Make a new service and send it into the background.
	zkill := zkillboard.NewZKillboard(redis)
	go zkill.Run()

	// Handle SIGINT and SIGTERM.
	ch := make(chan os.Signal, 2)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	// Stop the service gracefully.
	zkill.Close()
}
