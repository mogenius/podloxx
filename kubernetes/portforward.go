// Taken from https://github.com/gianarb/kube-port-forward
// Thanks for the wonderfull work @gianarb and the great blog entry

package kubernetes

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/mogenius/mo-go/logger"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
)

type PortForwardAPodRequest struct {
	Pod       v1.Pod
	LocalPort int
	PodPort   int
	// Steams configures where to write or read input from
	Streams genericclioptions.IOStreams
	// StopCh is the channel used to manage the port forward lifecycle
	StopCh <-chan struct{}
	// ReadyCh communicates when the tunnel is ready to receive traffic
	ReadyCh chan struct{}
}

func StartPortForward(useLocalKubeConfig bool) {
	pods := make([]v1.Pod, 0)

	logger.Log.Info("PORT_FORWARD: Waiting for pods to become available ...")

	for {
		var newPods = getPodloxxPodnames(useLocalKubeConfig)
		if len(newPods) > 0 {
			pods = append(pods, newPods...)
			break
		}
		time.Sleep(1 * time.Second)
		logger.Log.Infof(".")
	}

	for _, pod := range pods {
		logger.Log.Info("Starting PortForward for %s ...", pod.Name)

		var wg sync.WaitGroup
		wg.Add(1)

		// stopCh control the port forwarding lifecycle. When it gets closed the
		// port forward will terminate
		stopCh := make(chan struct{}, 1)
		// readyCh communicate when the port forward is ready to get traffic
		readyCh := make(chan struct{})
		// stream is used to tell the port forwarder where to place its output or
		// where to expect input if needed. For the port forwarding we just need
		// the output eventually
		stream := genericclioptions.IOStreams{
			In:     os.Stdin,
			Out:    os.Stdout,
			ErrOut: os.Stderr,
		}

		// managing termination signal from the terminal. As you can see the stopCh
		// gets closed to gracefully handle its termination.
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-sigs
			fmt.Println("Bye...")
			close(stopCh)
			wg.Done()
		}()

		go func() {
			// PortForward the pod specified from its port 9090 to the local port
			// 8080
			err := portForwardAPod(useLocalKubeConfig, PortForwardAPodRequest{
				Pod:       pod,
				LocalPort: 8080,
				PodPort:   1337,
				Streams:   stream,
				StopCh:    stopCh,
				ReadyCh:   readyCh,
			})
			if err != nil {
				panic(err)
			}
		}()

		select {
		case <-readyCh:
			break
		}
		println("Port forwarding is ready to get traffic. have fun!")

		wg.Wait()
	}
}

func portForwardAPod(useLocalKubeConfig bool, req PortForwardAPodRequest) error {
	const localPort = 8080
	const PodPort = 1337

	var kubeProvider *KubeProvider
	var err error
	if useLocalKubeConfig == true {
		kubeProvider, err = NewKubeProviderLocal()
	} else {
		kubeProvider, err = NewKubeProviderInCluster()
	}
	if err != nil {
		panic(err)
	}

	path := fmt.Sprintf("/api/v1/namespaces/%s/pods/%s/portforward", req.Pod.Namespace, req.Pod.Name)
	hostIP := strings.TrimLeft(kubeProvider.ClientConfig.Host, "htps:/")

	transport, upgrader, err := spdy.RoundTripperFor(&kubeProvider.ClientConfig)
	if err != nil {
		return err
	}

	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, http.MethodPost, &url.URL{Scheme: "https", Path: path, Host: hostIP})
	fw, err := portforward.New(dialer, []string{fmt.Sprintf("%d:%d", req.LocalPort, req.PodPort)}, req.StopCh, req.ReadyCh, req.Streams.Out, req.Streams.ErrOut)
	if err != nil {
		return err
	}
	return fw.ForwardPorts()
}

func getPodloxxPodnames(useLocalKubeConfig bool) []v1.Pod {
	var result []v1.Pod
	var kubeProvider *KubeProvider
	var err error
	if useLocalKubeConfig == true {
		kubeProvider, err = NewKubeProviderLocal()
	} else {
		kubeProvider, err = NewKubeProviderInCluster()
	}
	if err != nil {
		panic(err)
	}

	labelSelector := fmt.Sprintf("app=podloxx")
	pods, err := kubeProvider.ClientSet.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{LabelSelector: labelSelector})

	for _, pod := range pods.Items {
		result = append(result, pod)
	}

	if err != nil {
		fmt.Println("Error:", err)
		return result
	}
	return result
}
