package main

import (
	"context"
	"os/signal"
	"sync"

	"github.com/gin-gonic/gin"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background())
	defer stop()

	backendServerStruct := Back{
		activeContainers: 0,
		cpuCores:         4,
		coresOccupied:    0,
		memory:           16,
		memoryOccupied:   0,
		alive:            true,
		mux:              sync.RWMutex{},
	}
	router := gin.Default()

	router.GET("/health", func(ctx *gin.Context) {
		ctx.Status(200)
	})

	router.GET("/available-cpu", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"availableCpu": 4})
	})
}
