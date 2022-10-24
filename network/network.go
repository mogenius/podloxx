package network

import (
	"context"
	"fmt"
	"net"
	"os"
	"podloxx-collector/kubernetes"
	"podloxx-collector/structs"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/mogenius/mo-go/logger"
	"github.com/mogenius/mo-go/utils"
	"golang.org/x/exp/maps"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/jedib0t/go-pretty/v6/table"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

const (
	DEFAULTSNAPLEN                     = 65536   // Max Size of TCP Packets
	BYTES_CHANGE_SEND_TRESHHOLD uint64 = 1048576 // 1048576 wait until X bytes are gathered until we send an update to the API server
)

var TrafficData = make(map[string]*structs.InterfaceStats) // KEY: interfaceName e.g. veth1234abc
var lastTrafficDataBytesSum = make(map[string]uint64)      // KEY: interfaceName e.g. veth1234abc
var containerIds = make(map[string]v1.Pod)                 // KEY: containerId e.g. dad2f775d748b7fabdf333279219962a68af4f8bbf0e11933614bcba1d018de6
var containerPids = make(map[uint32]v1.Pod)                // KEY: HostProcessId e.g. 27123
var handles = make(map[string]*pcap.Handle)                // KEY: interfaceName e.g. veth1234abc

var appStartedAt = time.Now()
var ingressIps []net.IP

var mutex = &sync.Mutex{}

var eventCount uint64 = 0
var httpRequestCount uint64 = 0

var APIHOST string
var APIPORT string
var APIKEY string
var INTERFACEPREFIX string

var ReceiverChannel = make(chan structs.InterfaceStats)

var redisClient *redis.Client

func Init() {
	APIHOST = os.Getenv("API_HOST")
	APIPORT = os.Getenv("API_PORT")
	INTERFACEPREFIX = os.Getenv("INTERFACE_PREFIX")

	redisConnectionStr := fmt.Sprintf("%s:%s", os.Getenv("REDIS_SERVICE_NAME"), os.Getenv("REDIS_PORT"))
	logger.Log.Infof("REDIS: Connecting to: %s", redisConnectionStr)

	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisConnectionStr,
		Password: "",
		DB:       0,
	})
}

func MonitorAll(useLocalKubeConfig bool, overwriteInterfacePrefix string) {
	if overwriteInterfacePrefix != "" {
		INTERFACEPREFIX = overwriteInterfacePrefix
	}

	ingressIps = kubernetes.GetIngressControllerIps(false)

	go mainLoopAndWait(30)

	for {
		podWatch(&eventCount)
		logger.Log.Error("Watcher ended. Restarting ...")
		time.Sleep(3 * time.Second)
	}
}

// Periodically do all the work
func mainLoopAndWait(seconds time.Duration) {
	for {
		mainLoop()
		time.Sleep(seconds * time.Second)
	}
}

func mainLoop() {
	checkTaps()
	loadContainerPids()
	printEntriesTable()
}

func writeDataToRedis(data *structs.InterfaceStats) {
	// dataBytes := new(bytes.Buffer)
	// json.NewEncoder(dataBytes).Encode(data)

	err := redisClient.Set(data.PodName, data, 0).Err()
	if err != nil {
		logger.Log.Error(err)
	}
	val, err := redisClient.Get(data.PodName).Result()
	if err != nil {
		logger.Log.Error(err)
	}
	fmt.Println(val)
}

func MonitorLocal() {
	devices := GetAllDevices(false)
	var interesstingDevices []pcap.Interface
	for _, device := range devices {
		var ip string = ""
		for _, address := range device.Addresses {
			if len(address.IP) == 4 {
				ip = address.IP.String()
				break
			} else {
				ip = address.IP.String()
			}
		}
		if ip == "" {
			continue
		}
		interesstingDevices = append(interesstingDevices, device)
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"interface", "addresses"})
	var interfaceCount int = 0
	var addressCount int = 0
	for _, device := range interesstingDevices {
		if len(device.Addresses) > 0 {
			var addresses []string
			for _, address := range device.Addresses {
				addresses = append(addresses, address.IP.String())
				addressCount++
				if len(address.IP) == 4 {
					var ip = address.IP.String()
					go monitorInterface(ip, "LOCAL", device.Name, ip, time.Now().Format(time.RFC3339), "LOCAL")
				}
			}
			t.AppendRow(
				table.Row{device.Name, strings.Join(addresses[:], ", ")},
			)
			interfaceCount++
		}
	}
	t.AppendSeparator()
	t.AppendFooter(table.Row{"Count", "Count"})
	t.AppendFooter(table.Row{interfaceCount, addressCount})
	t.Render()

	for {
		reportData()
		printEntriesTable()
		time.Sleep(10 * time.Second)
	}
}

// Check if new interfaces need to be tapped
func checkTaps() {
	for containerId, container := range containerIds {
		var isTapped bool = false
		for _, trafficData := range TrafficData {
			if container.Name == trafficData.PodName {
				isTapped = true
			}
		}
		if isTapped == false {
			tapInterface(containerId, container)
		}
	}
}

// Check if new ProcessIds are available within the host
func loadContainerPids() {
	var newContainerPids, err = findContainerPids(containerIds)
	if err != nil {
		logger.Log.Error(err)
	}
	mutex.Lock()
	containerPids = newContainerPids
	mutex.Unlock()
}

// Listen to PodEvents to get notified for ADDED and DELETED pods
func podWatch(eventCount *uint64) error {
	provider, err := kubernetes.NewKubeProviderInCluster()
	var ownNodeName = os.Getenv("OWN_NODE_NAME")
	fieldSelector := ""
	if ownNodeName != "" {
		fieldSelector = fmt.Sprintf("metadata.namespace!=kube-system,spec.nodeName=%s", ownNodeName)
	} else {
		fieldSelector = "metadata.namespace!=kube-system"
	}
	logger.Log.Infof("Start watching for pods on node '%s' ...", ownNodeName)
	podsWatcher, err := provider.ClientSet.CoreV1().Pods("").Watch(context.TODO(), metav1.ListOptions{Watch: true, FieldSelector: fieldSelector})
	if err != nil {
		return err
	}
	podsChan := podsWatcher.ResultChan()
	for event := range podsChan {
		pod := event.Object.(*v1.Pod)
		*eventCount++

		switch event.Type {
		case watch.Added, watch.Modified:
			var mapToAdd = buildContainerIdsMap(*pod)
			for containerId, podToAdd := range mapToAdd {
				mutex.Lock()
				containerIds[containerId] = podToAdd
				mutex.Unlock()
				logger.Log.Info(*eventCount, "ADDED", pod.Name, containerId, pod.Status.Phase)
			}
		case watch.Deleted:
			var mapToDelete = buildContainerIdsMap(*pod)
			for containerId := range mapToDelete {
				logger.Log.Info(*eventCount, "DELETED", pod.Name, containerId, pod.Status.Phase)
				stopMonitoring(pod.Name)
			}
		case watch.Bookmark, watch.Error:
			// we do care yet
		}
	}
	return fmt.Errorf("podWatcher closed.")
}

// Monitor an host pod interface with pcap
func monitorInterface(podName string, namespace string, interfaceName string, ip string, startTime string, containerId string) {
	logger.Log.Noticef("Start monitoring: %s - %s (%s)", podName, interfaceName, ip)
	containerIp := net.ParseIP(ip)

	handle, err := pcap.OpenLive(interfaceName, DEFAULTSNAPLEN, true, pcap.BlockForever)

	var runsInCluster = true
	if containerId == "LOCAL" && namespace == "LOCAL" {
		runsInCluster = false
	}

	entry := structs.InitializeInterface(interfaceName, ip, podName, namespace, startTime, containerId, runsInCluster)
	mutex.Lock()
	TrafficData[interfaceName] = &entry
	handles[interfaceName] = handle
	mutex.Unlock()
	if err != nil {
		logger.Log.Errorf("ERROR (%s/%s): %s", podName, interfaceName, err)
		return
	}
	defer stopMonitoring(podName)

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	packets := packetSource.Packets()

	for pkt := range packets {
		if pkt.Layer(layers.LayerTypeTCP) != nil || pkt.Layer(layers.LayerTypeUDP) != nil {
			if ipFromGoPacketEndpoint(pkt.NetworkLayer().NetworkFlow().Src()).Equal(containerIp) {
				entry.TransmitBytes += uint64(pkt.Metadata().Length)

				if ipFromGoPacketEndpoint(pkt.NetworkLayer().NetworkFlow().Dst()).IsPrivate() {
					// check if not ingress-controller
					if !ipIsContainedInList(ipFromGoPacketEndpoint(pkt.NetworkLayer().NetworkFlow().Dst()), ingressIps) {
						entry.LocalTransmitBytes += uint64(pkt.Metadata().Length)
					}
				}
			} else {
				entry.ReceivedBytes += uint64(pkt.Metadata().Length)

				if ipFromGoPacketEndpoint(pkt.NetworkLayer().NetworkFlow().Src()).IsPrivate() && ipFromGoPacketEndpoint(pkt.NetworkLayer().NetworkFlow().Dst()).IsPrivate() {
					// check if not ingress-controller
					if !ipIsContainedInList(ipFromGoPacketEndpoint(pkt.NetworkLayer().NetworkFlow().Src()), ingressIps) {
						entry.LocalReceivedBytes += uint64(pkt.Metadata().Length)
					}
				}
			}
			entry.AddConnection(pkt.NetworkLayer().NetworkFlow().Src(), pkt.NetworkLayer().NetworkFlow().Dst())
		} else {
			entry.UnknownBytes += uint64(pkt.Metadata().Length)
		}
		entry.PacketsSum++
	}
}

// Stop Monitoring and cleanup everything
func stopMonitoring(podname string) {
	logger.Log.Warning("Stoping monitoring: ", podname)
	for interfaceName, entry := range TrafficData {
		if entry.PodName == podname {
			mutex.Lock()
			delete(TrafficData, interfaceName)
			delete(containerIds, entry.ContainerId)
			delete(lastTrafficDataBytesSum, interfaceName)
			handle, isOk := handles[interfaceName]
			if isOk {
				handle.Close()
				delete(handles, interfaceName)
			}
			mutex.Unlock()
			logger.Log.Warning("Stopped monitoring: ", podname)
			return
		}
	}
}

// Periodically print Information to stdout (statistics/debug/general information)
func printEntriesTable() {
	var totalPackets uint64 = 0
	var totalData uint64 = 0
	var totalTransmit uint64 = 0
	var totalReceived uint64 = 0
	var totalLocalTransmit uint64 = 0
	var totalLocalReceiced uint64 = 0
	var totalUnknown uint64 = 0
	var totalStartRx uint64 = 0
	var totalStartTx uint64 = 0
	for _, data := range TrafficData {
		totalPackets += data.PacketsSum
		totalData += (data.ReceivedBytes + data.TransmitBytes)
		totalLocalTransmit += data.LocalTransmitBytes
		totalLocalReceiced += data.LocalReceivedBytes
		totalTransmit += data.TransmitBytes
		totalReceived += data.ReceivedBytes
		totalUnknown += data.UnknownBytes
		totalStartRx += data.ReceivedStartBytes
		totalStartTx += data.TransmitStartBytes
	}
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{fmt.Sprintf("PODS (since %s)", utils.HumanDuration(time.Since(appStartedAt))), "Packets", "Transmit", "Received", "Total", "Unknown", "LocalRX", "LocalTX", "startTX", "startRX", "Age"})
	withTrafficCount := 0
	trafficArray := maps.Values(TrafficData)
	sort.Slice(trafficArray, func(i, j int) bool {
		return trafficArray[i].PacketsSum < trafficArray[j].PacketsSum
	})
	for _, entry := range trafficArray {
		if entry.PacketsSum > 10 {
			withTrafficCount++
			t.AppendRow(
				table.Row{entry.PodName, entry.PacketsSum, utils.BytesToHumanReadable(entry.TransmitBytes), utils.BytesToHumanReadable(entry.ReceivedBytes), utils.BytesToHumanReadable(entry.UnknownBytes + entry.TransmitBytes + entry.ReceivedBytes), utils.BytesToHumanReadable(entry.UnknownBytes), utils.BytesToHumanReadable(entry.LocalReceivedBytes), utils.BytesToHumanReadable(entry.LocalTransmitBytes), utils.BytesToHumanReadable(entry.TransmitStartBytes), utils.BytesToHumanReadable(entry.ReceivedStartBytes), utils.JsonStringToHumanDuration(entry.StartTime)},
			)
		}
	}
	t.AppendSeparator()
	t.AppendFooter(table.Row{"", "Packets", "Transmit", "Received", "Total", "Unknown", "LocalRX", "LocalTX", "startTX", "startRX"})
	t.AppendFooter(table.Row{fmt.Sprintf("%d Pods (%d with real traffic)", len(TrafficData), withTrafficCount), utils.NumberToHumanReadable(totalPackets), utils.BytesToHumanReadable(totalTransmit), utils.BytesToHumanReadable(totalReceived), utils.BytesToHumanReadable(totalData), utils.BytesToHumanReadable(totalUnknown), utils.BytesToHumanReadable(totalLocalReceiced), utils.BytesToHumanReadable(totalLocalTransmit), utils.BytesToHumanReadable(totalStartTx), utils.BytesToHumanReadable(totalStartRx), ""})
	t.Render()

	debugTable := table.NewWriter()
	debugTable.SetOutputMirror(os.Stdout)
	debugTable.AppendHeader(table.Row{"since", "ProcessedPodEvents", "TrafficData", "containerIds", "containerPids", "handles", "httpRequestCount", "ingressIps"})
	debugTable.AppendSeparator()
	debugTable.AppendRow(
		table.Row{utils.HumanDuration(time.Since(appStartedAt)), eventCount, len(TrafficData), len(containerIds), len(containerPids), len(handles), httpRequestCount, ingressIps},
	)
	debugTable.Render()
}

// Report Data to API server
func reportData() {
	for id, entry := range TrafficData {
		lastPacketSum, exists := lastTrafficDataBytesSum[id]
		if exists == false || (entry.TransmitBytes+entry.ReceivedBytes) >= lastPacketSum+BYTES_CHANGE_SEND_TRESHHOLD {
			// SEND DATA TO REDIS
			writeDataToRedis(entry)

			ReceiverChannel <- *entry
			httpRequestCount++
			lastTrafficDataBytesSum[id] = entry.TransmitBytes + entry.ReceivedBytes
		}
	}
}

// Create the connection between Pod and VETH to start monitoring
func tapInterface(containerId string, pod v1.Pod) error {
	var pid uint32 = 0

	for aPid, aPod := range containerPids {
		if aPod.Name == pod.Name {
			pid = aPid
			break
		}
	}

	if pid != 0 {
		index, err := getVethIndex(pid)
		if err != nil {
			if strings.Contains(err.Error(), "Permission denied") ||
				strings.Contains(err.Error(), "No such file or directory") ||
				strings.Contains(err.Error(), "Empty vethIndex.") {
				cleanBecauseOfErrors(pid, pod)
			}
			return fmt.Errorf("GetVethIndex (%s): %s", pod.Name, err.Error())
		}
		vethName, err := getVethInterfaceForIndex(index, INTERFACEPREFIX)
		if err != nil {
			if strings.Contains(err.Error(), "exit status 2") {
				cleanBecauseOfErrors(pid, pod)
			}
			return fmt.Errorf("GetVethInterfaceForIndex (%s): %s", pod.Name, err.Error())
		}
		logger.Log.Info(pid, index, vethName, pod.Name)
		if _, exists := TrafficData[vethName]; !exists {
			go monitorInterface(pod.Name, pod.Namespace, vethName, pod.Status.PodIP, pod.Status.StartTime.Format(time.RFC3339), containerId)
		}
	}
	return nil
}

// Cleanup stuff from tapInterface if something goes wrong
func cleanBecauseOfErrors(pid uint32, pod v1.Pod) {
	logger.Log.Warningf("REMOVING '%s' from containerPids and containerIds.", pod.Name)
	mutex.Lock()
	delete(containerPids, pid)
	mutex.Unlock()
	logger.Log.Warningf("REMOVED '%d' from containerPids.", pid)
	for id := range buildContainerIdsMap(pod) {
		mutex.Lock()
		delete(containerIds, id)
		mutex.Unlock()
		logger.Log.Warningf("REMOVED '%s' from containerIds.", id)
	}
}

// Search /hostproc for folders (numeric only like 27122) because the will be pids which contain the cGroup which lets us connect containerid to virtual interface in the host
// FOLLOWING CODE HAS BEEN COPIED FROM https://github.com/up9inc/mizu/tree/main Thanks for the great work @UP9 Inc
func findContainerPids(containerIds map[string]v1.Pod) (map[uint32]v1.Pod, error) {
	result := make(map[uint32]v1.Pod)

	pidFiles, err := os.ReadDir(kubernetes.PROCFSMOUNTPATH)

	if err != nil {
		return result, err
	}

	for _, pid := range pidFiles {
		if !pid.IsDir() {
			continue
		}

		if !regexp.MustCompile("[0-9]+").MatchString(pid.Name()) {
			continue
		}

		pidNumber, errAtoi := strconv.Atoi(pid.Name())
		if errAtoi != nil {
			continue
		}

		cgroup, err := getProcessCgroup(pid.Name())
		if err != nil {
			mutex.Lock()
			delete(containerPids, uint32(pidNumber))
			mutex.Unlock()
			logger.Log.Warningf("ProcessCgroup Error (%s): %s", pid.Name(), err)
			continue
		}

		mutex.Lock()
		pod, ok := containerIds[cgroup]
		mutex.Unlock()
		if !ok {
			continue
		}

		result[uint32(pidNumber)] = pod
	}

	return result, nil
}
