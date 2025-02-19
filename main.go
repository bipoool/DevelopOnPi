package main

import (
	"context"
	"developOnPi/backendServer"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	backendServerStruct := backendServer.GetNewBackendServer()
	backendServerController := backendServer.GetNewBackendServerController(backendServerStruct)

	router := gin.Default()

	router.GET("/status", backendServerController.CheckStatus)
	router.GET("/health", backendServerController.IsAlive)
	router.GET("/available-cpu", backendServerController.CheckIfCpuIsAvailable)
	router.GET("/available-memory", backendServerController.CheckIfMemoryIsAvailable)
	router.POST("/spin-up-vm", backendServerController.SpinUpVM)
	router.GET("/check-vm-status", backendServerController.CheckVmStatus)
	router.PUT("/change-vm-status", backendServerController.ChangeStatus)

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Fatal(err)
		}
	}()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		fmt.Println("ListenAndServe() error" + err.Error())
	}

}
