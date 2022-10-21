package kubernetes

import (
	"path/filepath"
	"podloxx-collector/version"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var (
	NAMESPACE              = "default"
	DAEMONSETNAME          = "podloxx"
	DAEMONSETIMAGE         = "ghcr.io/mogenius/podloxx-collector:" + version.Ver
	REDISNAME              = "podloxx-redis"
	REDISSERVICENAME       = "podloxx-redis-service"
	REDISIMAGE             = "redis:latest"
	REDISPORT        int32 = 6379
	REDISTARGETPORT        = "redis"
	PROCFSVOLUMENAME       = "proc"
	PROCFSMOUNTPATH        = "/hostproc"
	SYSFSVOLUMENAME        = "sys"
	SYSFSMOUNTPATH         = "/sys"

	SERVICEACCOUNTNAME     = "podloxx-service-account-app"
	CLUSTERROLENAME        = "podloxx-cluster-role-app"
	CLUSTERROLEBINDINGNAME = "podloxx-cluster-role-binding-app"
	RBACRESOURCES          = []string{"pods", "services", "endpoints"}
)

type KubeProvider struct {
	ClientSet    *kubernetes.Clientset
	ClientConfig rest.Config
}

func NewKubeProviderLocal() (*KubeProvider, error) {
	var kubeconfig string = ""
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	restConfig, errConfig := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if errConfig != nil {
		panic(errConfig.Error())
	}

	clientSet, errClientSet := kubernetes.NewForConfig(restConfig)
	if errClientSet != nil {
		panic(errClientSet.Error())
	}

	return &KubeProvider{
		ClientSet:    clientSet,
		ClientConfig: *restConfig,
	}, nil
}

func NewKubeProviderInCluster() (*KubeProvider, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return &KubeProvider{
		ClientSet:    clientset,
		ClientConfig: *config,
	}, nil
}

func int32Ptr(i int32) *int32 { return &i }
