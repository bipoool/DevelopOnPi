package backendServer

import (
	"sync"
)

type BackendServerStruct struct {
	vmId             uint16
	sshPort          uint16
	activeContainers int16
	cpuCores         int8
	coresOccupied    int8
	memory           int16
	memoryOccupied   int16
	alive            bool
	mux              sync.RWMutex
}

type BackendServerInterface interface {
	isAlive() bool
	setAlive(bool)
	getActiveContainers() int16
	cpuAvailable() int8
	memoryAvailable() int16
	useCpu(int8)
	useMemory(int16)
	incrementActiveContainers()
	rotateSshPort() uint16
	rotateVmId() uint16
}

func (backendServer *BackendServerStruct) isAlive() bool {
	return backendServer.alive
}

func (backendServer *BackendServerStruct) setAlive(isAlive bool) {
	backendServer.alive = isAlive
}

func (backendServer *BackendServerStruct) getActiveContainers() int16 {
	return backendServer.activeContainers
}

func (backendServer *BackendServerStruct) cpuAvailable() int8 {
	return backendServer.cpuCores - backendServer.coresOccupied
}

func (backendServer *BackendServerStruct) useCpu(cpuCores int8) {
	backendServer.coresOccupied += cpuCores
}

func (backendServer *BackendServerStruct) memoryAvailable() int16 {
	return backendServer.memory - backendServer.memoryOccupied
}

func (backendServer *BackendServerStruct) useMemory(memory int16) {
	backendServer.memoryOccupied += memory
}

func (backendServer *BackendServerStruct) incrementActiveContainers() {
	backendServer.activeContainers++
}

func (backendServer *BackendServerStruct) rotateSshPort() uint16 {
	backendServer.sshPort++
	return backendServer.sshPort
}

func (backendServer *BackendServerStruct) rotateVmId() uint16 {
	backendServer.vmId++
	return backendServer.vmId
}

func GetNewBackendServer() *BackendServerStruct {
	return &BackendServerStruct{
		vmId:             0,
		sshPort:          60000,
		activeContainers: 0,
		cpuCores:         4,
		coresOccupied:    0,
		memory:           8,
		memoryOccupied:   0,
		alive:            true,
	}
}
