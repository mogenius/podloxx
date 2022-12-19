package network

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"podloxx/kubernetes"
	"podloxx/logger"
	"podloxx/structs"
	"runtime"
	"strings"

	"podloxx/utils"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
)

func GetAllDevices(print bool) []pcap.Interface {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		logger.Log.Fatal(err)
	}

	if print {
		for index, device := range devices {
			if len(device.Addresses) > 0 {
				fmt.Println("Interface: ", device.Name)
				for _, address := range device.Addresses {
					var ipType string
					if ipType = "IPv4"; len(address.IP) == 4 {
						ipType = "😎 IPv4"
					} else {
						ipType = "IPv6"
					}
					fmt.Println("- IP address:  ", address.IP, " ", ipType)
					fmt.Println("- Subnet mask: ", address.Netmask)
				}
				fmt.Println("")
			} else {
				removeIndex(devices, index)
			}
		}
	}
	return devices
}

func removeIndex(s []pcap.Interface, index int) []pcap.Interface {
	return append(s[:index], s[index+1:]...)
}

// Parse IP bytes from gopacket Endpoint into golang net ip (this is way faster than parse-ip)
// we use the functions a million times per minute so it needs to be as fast as possible
func ipFromGoPacketEndpoint(endpoint gopacket.Endpoint) net.IP {
	return net.IPv4(endpoint.Raw()[0], endpoint.Raw()[1], endpoint.Raw()[2], endpoint.Raw()[3])
}

func ipIsContainedInList(ip net.IP, list []net.IP) bool {
	for i := 0; i < len(list); i++ {
		if ip.Equal(list[i]) {
			return true
		}
	}
	return false
}

// Print infos to stdout in a human readable format (mostly debug and eye-candy)
func printLog(entry structs.InterfaceStats) {
	logger.Log.Info(entry.PodName, " => ", entry.PacketsSum, ": ", utils.BytesToHumanReadable(entry.TransmitBytes), "/", utils.BytesToHumanReadable(entry.ReceivedBytes), " (", utils.BytesToHumanReadable(entry.UnknownBytes), ") => ", utils.BytesToHumanReadable(entry.UnknownBytes+entry.TransmitBytes+entry.ReceivedBytes), " LOCAL: ", utils.BytesToHumanReadable(entry.LocalReceivedBytes+entry.LocalTransmitBytes))
}

// Get the index of the virtual namespace in the containers network namespace
func getVethIndex(pid uint32) (string, error) {
	switch runtime.GOOS {
	case "darwin", "linux":
		var cmdStr string = fmt.Sprintf(`nsenter --net=/hostproc/%d/ns/net -n ip link | sed -n -e 's/.*eth0@if\([0-9]*\):.*/\1/p'`, pid)
		cmd := exec.Command("bash", "-c", cmdStr)
		out, err := cmd.CombinedOutput()
		if err != nil {
			logger.Log.Info(cmdStr)
			logger.Log.Error(err)
			return "", err
		}

		vethIndex := strings.Replace(string(out), "\n", "", -1)
		if vethIndex == "" {
			return "", fmt.Errorf("Empty vethIndex.")
		}
		return vethIndex, nil
	default:
		return "", fmt.Errorf("Unknown OS: %s", runtime.GOOS)
	}
}

// Get the virtual host interface name from the index provided in the containers network namespace
func getVethInterfaceForIndex(index string) (string, error) {
	switch runtime.GOOS {
	case "darwin", "linux":
		var cmdStr string = fmt.Sprintf(`ip -o link | grep ^%s: | sed -n -e 's/.* \([[:alnum:]]*\)@.*/\1/p'`, index)
		cmd := exec.Command("bash", "-c", cmdStr)
		out, err := cmd.CombinedOutput()
		if err != nil {
			logger.Log.Info(cmdStr)
			logger.Log.Error(err)
			return "", err
		}

		veth := strings.Replace(string(out), "\n", "", -1)
		if veth == "" {
			return "", fmt.Errorf("Empty veth. Skipping.")
		}
		return veth, nil
	default:
		return "", fmt.Errorf("Unknown OS: %s", runtime.GOOS)
	}
}

// FOLLOWING CODE HAS BEEN COPIED FROM https://github.com/up9inc/mizu/tree/main Thanks for the great work @UP9 Inc
func buildContainerIdsMap(pod v1.Pod) map[string]v1.Pod {
	result := make(map[string]v1.Pod)

	for _, container := range pod.Status.ContainerStatuses {
		parsedUrl, err := url.Parse(container.ContainerID)

		if err != nil {
			logger.Log.Warningf("Expecting URL like container ID %v", container.ContainerID)
			continue
		}

		result[parsedUrl.Host] = pod
	}

	return result
}

// FOLLOWING CODE HAS BEEN COPIED FROM https://github.com/up9inc/mizu/tree/main Thanks for the great work @UP9 Inc
func getProcessCgroup(pid string) (string, error) {
	filePath := fmt.Sprintf("%s/%s/cgroup", kubernetes.PROCFSMOUNTPATH, pid)

	bytes, err := os.ReadFile(filePath)

	if err != nil {
		return "", fmt.Errorf("Error reading cgroup file %s - %v", filePath, err)
	}

	lines := strings.Split(string(bytes), "\n")
	cgrouppath := extractCgroup(lines)

	if cgrouppath == "" {
		return "", errors.Errorf("Cgroup path not found for %s, %s", pid, lines)
	}

	return normalizeCgroup(cgrouppath), nil
}

// FOLLOWING CODE HAS BEEN COPIED FROM https://github.com/up9inc/mizu/tree/main Thanks for the great work @UP9 Inc
func extractCgroup(lines []string) string {
	if len(lines) == 1 {
		parts := strings.Split(lines[0], ":")
		return parts[len(parts)-1]
	} else {
		for _, line := range lines {
			if strings.Contains(line, ":pids:") {
				parts := strings.Split(line, ":")
				return parts[len(parts)-1]
			}
		}
	}

	return ""
}

// FOLLOWING CODE HAS BEEN COPIED FROM https://github.com/up9inc/mizu/tree/main Thanks for the great work @UP9 Inc
func normalizeCgroup(cgrouppath string) string {
	basename := strings.TrimSpace(path.Base(cgrouppath))

	if strings.Contains(basename, "-") {
		basename = basename[strings.Index(basename, "-")+1:]
	}

	if strings.Contains(basename, ".") {
		return strings.TrimSuffix(basename, filepath.Ext(basename))
	} else {
		return basename
	}
}
