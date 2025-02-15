package backendServer

import (
	"fmt"
	"net/http"
	"os/exec"
	"strconv"

	"github.com/gin-gonic/gin"
)

type VmConfig struct {
	Cpu       int8   `json:"cpu"`
	Memory    int16  `json:"memory"`
	PublicKey string `json:"publicKey"`
	VmId      uint16 `json:"vmId"`
	Status    uint8  `json:"status"` // 0: In Progress (stop/start), 1: Running, 2: Stopped, 3: Error, 4: Terminated
}

type BackendServerController struct {
	backendServer BackendServerInterface
	vmStatusMap   map[uint16]*VmConfig
}

func (backendServerController *BackendServerController) CheckStatus(ctx *gin.Context) {
	status := backendServerController.backendServer.isAlive()
	activeContainers := backendServerController.backendServer.getActiveContainers()
	cpuAvailable := backendServerController.backendServer.cpuAvailable()
	memoryAvailable := backendServerController.backendServer.memoryAvailable()
	ctx.JSON(http.StatusOK, gin.H{"status": status, "activeContainers": activeContainers,
		"cpuAvailable": cpuAvailable, "memoryAvailable": memoryAvailable})
}

func (backendServerController *BackendServerController) CheckIfCpuIsAvailable(ctx *gin.Context) {
	isCpuAvailable := backendServerController.backendServer.cpuAvailable()
	ctx.JSON(http.StatusOK, gin.H{"isCpuAvailable": isCpuAvailable})
}

func (backendServerController *BackendServerController) CheckIfMemoryIsAvailable(ctx *gin.Context) {
	isMemoryAvailable := backendServerController.backendServer.memoryAvailable()
	ctx.JSON(http.StatusOK, gin.H{"isMemoryAvailable": isMemoryAvailable})
}

func (backendServerController *BackendServerController) IsAlive(ctx *gin.Context) {
	flagIsAlive := backendServerController.backendServer.isAlive()
	var httpCode int
	if flagIsAlive {
		httpCode = http.StatusOK
	} else {
		httpCode = http.StatusServiceUnavailable
	}
	ctx.JSON(httpCode, gin.H{"isAlive": flagIsAlive})
}

func (backendServerController *BackendServerController) SpinUpVM(ctx *gin.Context) {

	var vmConfig VmConfig

	if err := ctx.ShouldBindJSON(&vmConfig); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if CPU and RAM is available
	if backendServerController.backendServer.cpuAvailable() < vmConfig.Cpu || vmConfig.Cpu <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "CPUs not available"})
		return
	}

	if backendServerController.backendServer.memoryAvailable() < vmConfig.Memory || vmConfig.Memory <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Memory not available"})
		return
	}

	vmConfig.Status = 0
	backendServerController.backendServer.incrementActiveContainers()
	sshPort := backendServerController.backendServer.rotateSshPort()
	vmId := backendServerController.backendServer.rotateVmId()
	vmConfig.VmId = vmId
	backendServerController.vmStatusMap[vmId] = &vmConfig

	cmd := exec.Command("/home/vipulgupta/workspace/developWithPi/scripts/spinUpUbuntuContainer.sh",
		strconv.FormatUint(uint64(sshPort), 10),
		strconv.FormatInt(int64(vmConfig.Cpu), 10),
		strconv.FormatInt(int64(vmConfig.Memory), 10),
		vmConfig.PublicKey,
		strconv.FormatUint(uint64(vmId), 10),
	)

	go func() {
		_, err := cmd.CombinedOutput()
		if err != nil {
			vmConfig.Status = 3
			fmt.Println("Error spinning up VM:", err)
			return
		} else {
			vmConfig.Status = 1
			backendServerController.backendServer.useCpu(vmConfig.Cpu)
			backendServerController.backendServer.useMemory(vmConfig.Memory)
		}
	}()

	ctx.JSON(http.StatusOK, gin.H{"message": "VM spinning up", "vmId": vmId, "sshPort": sshPort})
}

func (backendServerController *BackendServerController) CheckVmStatus(ctx *gin.Context) {
	vmId, err := strconv.ParseUint(ctx.Query("vmId"), 10, 16)
	if err != nil || backendServerController.vmStatusMap[uint16(vmId)] == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vmId"})
		return
	}

	vmConfig := backendServerController.vmStatusMap[uint16(vmId)]
	ctx.JSON(http.StatusOK, gin.H{"status": vmConfig.Status})
}

func (backendServerController *BackendServerController) ChangeStatus(ctx *gin.Context) {
	var vmConfig VmConfig

	if err := ctx.ShouldBindJSON(&vmConfig); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if backendServerController.vmStatusMap[vmConfig.VmId] == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vmId"})
		return
	}

	status := vmConfig.Status
	currentConfig := backendServerController.vmStatusMap[vmConfig.VmId]

	if currentConfig.Status == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "VM is in progress"})
		return
	}

	if status == 1 {
		if currentConfig.Status != 2 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "VM is already running/terminated"})
			return
		}
		cmd := exec.Command("/home/vipulgupta/workspace/developWithPi/scripts/startUbuntuContainer.sh",
			strconv.FormatUint(uint64(currentConfig.VmId), 10),
		)
		go func() {
			_, err := cmd.CombinedOutput()
			if err != nil {
				currentConfig.Status = 3
				fmt.Println("Error starting up VM:", err)
				return
			} else {
				currentConfig.Status = 1
			}
		}()
		ctx.JSON(http.StatusOK, gin.H{"message": "VM starting up"})
	} else if status == 2 {
		if currentConfig.Status != 1 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "VM is already stopped/terminated"})
			return
		}
		cmd := exec.Command("/home/vipulgupta/workspace/developWithPi/scripts/stopUbuntuContainer.sh",
			strconv.FormatUint(uint64(currentConfig.VmId), 10),
		)
		fmt.Println(cmd.String())
		go func() {
			_, err := cmd.CombinedOutput()
			if err != nil {
				currentConfig.Status = 3
				fmt.Println("Error Stopping up VM:", err)
				return
			} else {
				currentConfig.Status = 2
			}
		}()
		ctx.JSON(http.StatusOK, gin.H{"message": "VM stopping"})
	} else if status == 4 {
		if currentConfig.Status != 2 && currentConfig.Status != 1 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "VM is already terminated"})
			return
		}
		cmd := exec.Command("/home/vipulgupta/workspace/developWithPi/scripts/terminateUbuntuContainer.sh",
			strconv.FormatUint(uint64(currentConfig.VmId), 10),
		)
		go func() {
			_, err := cmd.CombinedOutput()
			if err != nil {
				currentConfig.Status = 3
				fmt.Println("Error Terminating up VM:", err)
				return
			} else {
				currentConfig.Status = 4
				backendServerController.backendServer.useCpu(-currentConfig.Cpu)
				backendServerController.backendServer.useMemory(-currentConfig.Memory)
			}
		}()
		ctx.JSON(http.StatusOK, gin.H{"message": "VM terminating"})
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
	}
	backendServerController.vmStatusMap[vmConfig.VmId] = currentConfig
}

func GetNewBackendServerController(backendServer BackendServerInterface) *BackendServerController {
	return &BackendServerController{
		backendServer: backendServer,
		vmStatusMap:   make(map[uint16]*VmConfig),
	}
}
