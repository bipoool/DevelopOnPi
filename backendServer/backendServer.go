package backendServer

import (
	"sync"
)

type BackendServerStruct struct {
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
	isCpuAvailable() int8
	isMemoryAvailable() int16
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

func (backendServer *BackendServerStruct) isCpuAvailable() int8 {
	return backendServer.cpuCores - backendServer.coresOccupied
}

func (backendServer *BackendServerStruct) isMemoryAvailable() int16 {
	return backendServer.memory - backendServer.memoryOccupied
}
