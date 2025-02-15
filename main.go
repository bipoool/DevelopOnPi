package main

import (
	"context"
	"os/signal"

	"github.com/gin-gonic/gin"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background())
	defer stop()

	router := gin.Default()

	router.GET("/health", func(ctx *gin.Context) {
		ctx.Status(200)
	})

	router.GET("/status")
}
