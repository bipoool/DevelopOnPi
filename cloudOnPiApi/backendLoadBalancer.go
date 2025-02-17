package cloudonpiapi

import (
	"bytes"
	"developOnPi/backendServer"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
)

type BackendStatusStruct struct {
	VmId             uint16 `json:"VmId"`
	FlagIsAlive      bool   `json:"flagIsAlive"`
	ActiveContainers int16  `json:"activeContainers"`
	CpuAvailable     int8   `json:"cpuAvailable"`
	MemoryAvailable  int16  `json:"memoryAvailable"`
}

type SpinUpVmResponseStruct struct {
	Message string
	VmId    uint16
	SshPort uint16
}

type BackendSetStruct struct {
	backendServerEndpoints []string
}

type BackendSetInterface interface {
	getBackendServerEndpoints() []string
	addBackendServerEndpoint(string)
	removeEndpoint(string)
	getBackendForDeployment(backendServer.VmConfig)
}

func (backendSet *BackendSetStruct) getBackendServerEndpoints() []string {
	return backendSet.backendServerEndpoints
}

func (backendSet *BackendSetStruct) addBackendServerEndpoint(endpoint string) {
	backendSet.backendServerEndpoints = append(backendSet.backendServerEndpoints, endpoint)
}

func (backendSet *BackendSetStruct) removeEndpoint(endpoint string) {
	for i, v := range backendSet.backendServerEndpoints {
		if v == endpoint {
			backendSet.backendServerEndpoints = append(backendSet.backendServerEndpoints[0:i], backendSet.backendServerEndpoints[i+1:]...)
		}
	}
}

func (backendSet *BackendSetStruct) getBackendForDeployment(vmConfig backendServer.VmConfig) string {
	cpu := vmConfig.Cpu
	memory := vmConfig.Memory
	mapOfHostEndpointToStatus := make(map[string]BackendStatusStruct)

	for _, v := range backendSet.backendServerEndpoints {
		endpoint := v
		req, _ := http.NewRequest("GET", endpoint+BackendStatusEndpoint, nil)
		req.Header.Set("Accept", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error making request:", err)
			continue
		}
		var backendStatus BackendStatusStruct
		if err := json.NewDecoder(resp.Body).Decode(&backendStatus); err != nil {
			fmt.Println("Error decoding response:", err)
			continue
		}
		resp.Body.Close()
		mapOfHostEndpointToStatus[endpoint] = backendStatus
	}

	var finalBackend string
	var leastContainerFound int16 = math.MaxInt16

	for endpoint, status := range mapOfHostEndpointToStatus {
		if status.CpuAvailable > cpu &&
			status.MemoryAvailable > memory &&
			status.FlagIsAlive &&
			status.ActiveContainers < leastContainerFound {
			leastContainerFound = status.ActiveContainers
			finalBackend = endpoint
		}
	}

	return finalBackend
}

func (backendSet *BackendSetStruct) spinUpVm(vmConfig backendServer.VmConfig, hostEndPoint string) string {

	requestBody, err := json.Marshal(vmConfig)

	req, _ := http.NewRequest("POST", hostEndPoint+BackendSpinUpVmEndpoint, bytes.NewBuffer(requestBody))
	req.Header.Set("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
	}
	var spinUpVmResponse SpinUpVmResponseStruct
	if err := json.NewDecoder(resp.Body).Decode(&spinUpVmResponse); err != nil {
		fmt.Println("Error decoding response:", err)
	}
	resp.Body.Close()
	sshPort := spinUpVmResponse.SshPort
	vmId := spinUpVmResponse.VmId
	subdomain := strconv.FormatUint(uint64(vmId), 10) +
		".cloudOnPi.com:" +
		strconv.FormatUint(uint64(sshPort), 10)
	return subdomain
}

func GetNewBackendSet() *BackendSetStruct {
	return &BackendSetStruct{
		backendServerEndpoints: []string{},
	}
}
