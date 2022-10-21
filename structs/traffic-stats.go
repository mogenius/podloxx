package structs

import (
	"os"
	"strconv"
	"strings"
	"unsafe"

	"github.com/google/gopacket"
	"github.com/mogenius/mo-go/logger"
)

const MAXCONNECTIONSSIZE int = 10

type InterfaceStats struct {
	Name               string                         `json:"fileName"`
	Ip                 string                         `json:"ip"`
	PodName            string                         `json:"podName"`
	Namespace          string                         `json:"namespace"`
	Node               string                         `json:"node"`
	ContainerId        string                         `json:"containerId"`
	PacketsSum         uint64                         `json:"packetsSum"`
	TransmitBytes      uint64                         `json:"transmitBytes"`
	ReceivedBytes      uint64                         `json:"receivedBytes"`
	UnknownBytes       uint64                         `json:"unknownBytes"`
	LocalTransmitBytes uint64                         `json:"localTransmitBytes"`
	LocalReceivedBytes uint64                         `json:"localReceivedBytes"`
	TransmitStartBytes uint64                         `json:"transmitStartBytes"`
	ReceivedStartBytes uint64                         `json:"receivedStartBytes"`
	StartTime          string                         `json:"startTime"`
	Connections        map[uint64]InterfaceConnection `json:"connections"`
}

type InterfaceConnection struct {
	Ip1       string `json:"ip1"`
	Ip2       string `json:"ip2"`
	PacketSum uint64 `json:"packetSum"`
}

func (is InterfaceStats) AddConnection(ip1 gopacket.Endpoint, ip2 gopacket.Endpoint) {
	var hash uint64 = 0
	if ip1.LessThan(ip2) {
		hash = *(*uint64)(unsafe.Pointer(&append(ip1.Raw(), ip2.Raw()...)[0]))
	} else {
		hash = *(*uint64)(unsafe.Pointer(&append(ip2.Raw(), ip1.Raw()...)[0]))
	}
	entry, doesExist := is.Connections[hash]
	if doesExist {
		entry.PacketSum += 1
		is.Connections[hash] = entry
	} else {
		if len(is.Connections) < MAXCONNECTIONSSIZE {
			is.Connections[hash] = InterfaceConnection{Ip1: ip1.String(), Ip2: ip2.String(), PacketSum: 1}
		}
	}
}

func InitializeInterface(name string, ip string, podName string, namespace string, startTime string, containerId string, runsInCluster bool) InterfaceStats {
	entry := InterfaceStats{}
	entry.Name = name
	entry.Ip = ip
	entry.PodName = podName
	entry.ContainerId = containerId
	entry.Namespace = namespace
	entry.Node = os.Getenv("OWN_NODE_NAME")
	if runsInCluster {
		entry.ReceivedStartBytes = loadUint64FromFile("/sys/class/net/" + name + "/statistics/rx_bytes")
		entry.TransmitStartBytes = loadUint64FromFile("/sys/class/net/" + name + "/statistics/tx_bytes")
	}
	entry.StartTime = startTime
	entry.Connections = make(map[uint64]InterfaceConnection)
	return entry
}

func loadUint64FromFile(filePath string) uint64 {
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		logger.Log.Error(err)
		return uint64(0)
	}

	var stringData = strings.TrimSuffix(string(fileContent), "\n")
	number, err := strconv.ParseUint(stringData, 10, 64)
	if err != nil {
		logger.Log.Error(err)
		return uint64(0)
	}

	return number
}
