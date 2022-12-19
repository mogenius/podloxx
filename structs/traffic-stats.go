package structs

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"podloxx/logger"
	"strconv"
	"strings"
	"unsafe"

	"github.com/google/gopacket"
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

type InterfaceStatsNumbers struct {
	PacketsSum         uint64 `json:"packetsSum"`
	TransmitBytes      uint64 `json:"transmitBytes"`
	ReceivedBytes      uint64 `json:"receivedBytes"`
	UnknownBytes       uint64 `json:"unknownBytes"`
	LocalTransmitBytes uint64 `json:"localTransmitBytes"`
	LocalReceivedBytes uint64 `json:"localReceivedBytes"`
}

type Overview struct {
	PacketsSum              uint64 `json:"packetsSum"`
	TransmitBytes           uint64 `json:"transmitBytes"`
	ReceivedBytes           uint64 `json:"receivedBytes"`
	UnknownBytes            uint64 `json:"unknownBytes"`
	LocalTransmitBytes      uint64 `json:"localTransmitBytes"`
	LocalReceivedBytes      uint64 `json:"localReceivedBytes"`
	TotalNodes              int    `json:"totalNodes"`
	TotalNamespaces         int    `json:"totalNamespaces"`
	TotalPods               uint64 `json:"totalPods"`
	Uptime                  string `json:"uptime"`
	LastUpdate              int64  `json:"lastUpdate"`
	ExternalBandwidthPerSec string `json:"externalBandwidthPerSec"`
	InternalBandwidthPerSec string `json:"internalBandwidthPerSec"`
	PacketsPerSec           int    `json:"packetsPerSec"`
}

func (i InterfaceStats) MarshalBinary() (data []byte, err error) {
	bytes, err := json.Marshal(i)
	return bytes, err
}

func PrettyPrint(i interface{}) {
	iJson, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Printf("%s\n", string(iJson))
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

func CopyInterface(src InterfaceStats) InterfaceStats {
	dst := InterfaceStats{}
	dst.Name = src.Name
	dst.Ip = src.Ip
	dst.PodName = src.PodName
	dst.Namespace = src.Namespace
	dst.ContainerId = src.ContainerId
	dst.PacketsSum = src.PacketsSum
	dst.TransmitBytes = src.TransmitBytes
	dst.ReceivedBytes = src.ReceivedBytes
	dst.UnknownBytes = src.UnknownBytes
	dst.LocalTransmitBytes = src.LocalTransmitBytes
	dst.LocalReceivedBytes = src.LocalReceivedBytes
	dst.TransmitStartBytes = src.TransmitStartBytes
	dst.ReceivedStartBytes = src.ReceivedStartBytes
	dst.StartTime = src.StartTime
	dst.Connections = src.Connections
	return dst
}

func Minify(src InterfaceStats) InterfaceStatsNumbers {
	dst := InterfaceStatsNumbers{}
	dst.PacketsSum = src.PacketsSum
	dst.TransmitBytes = src.TransmitBytes
	dst.ReceivedBytes = src.ReceivedBytes
	dst.UnknownBytes = src.UnknownBytes
	dst.LocalTransmitBytes = src.LocalTransmitBytes
	dst.LocalReceivedBytes = src.LocalReceivedBytes
	return dst
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
