package backendServer

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type BackendServerController struct {
	backendServer BackendServerInterface
}

func (backendServerController *BackendServerController) checkStatus(ctx *gin.Context) {
	status := backendServerController.backendServer.isAlive()
	activeConnections := backendServerController.backendServer.getActiveContainers()
	ctx.JSON(http.StatusOK, gin.H{"status": status, "activeConnections": activeConnections})
}

func (backendServerController *BackendServerController) checkIfCpuIsAvailable(ctx *gin.Context) {
	isCpuAvailable := backendServerController.backendServer.isCpuAvailable()
	ctx.JSON(http.StatusOK, gin.H{"isCpuAvailable": isCpuAvailable})
}

func (backendServerController *BackendServerController) checkIfMemoryIsAvailable(ctx *gin.Context) {
	isMemoryAvailable := backendServerController.backendServer.isMemoryAvailable()
	ctx.JSON(http.StatusOK, gin.H{"isMemoryAvailable": isMemoryAvailable})
}
