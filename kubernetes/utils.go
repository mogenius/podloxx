package kubernetes

import (
	"flag"
	"path/filepath"
	"podloxx-collector/version"

	"github.com/mogenius/mo-go/logger"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var (
	NAMESPACE        = "default"
	DAEMONSETNAME    = "podloxx"
	DAEMONSETIMAGE   = "ghcr.io/mogenius/podloxx-collector:" + version.Ver
	PROCFSVOLUMENAME = "proc"
	PROCFSMOUNTPATH  = "/hostproc"
	SYSFSVOLUMENAME  = "sys"
	SYSFSMOUNTPATH   = "/sys"

	SERVICEACCOUNTNAME     = "podloxx-service-account-app"
	CLUSTERROLENAME        = "podloxx-cluster-role-app"
	CLUSTERROLEBINDINGNAME = "podloxx-cluster-role-binding-app"
	RBACRESOURCES          = []string{"pods", "services", "endpoints"}
)

type KubeProvider struct {
	ClientSet *kubernetes.Clientset
	//kubernetesConfig clientcmd.ClientConfig
	ClientConfig rest.Config
}

func NewKubeProviderLocal() (*KubeProvider, error) {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	restConfig, errConfig := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if errConfig != nil {
		panic(errConfig.Error())
	}

	clientSet, errClientSet := kubernetes.NewForConfig(restConfig)
	if errClientSet != nil {
		panic(errClientSet.Error())
	}

	logger.Log.Debugf("K8s client config (init with .kube/config), host: %s", restConfig.Host)

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

	logger.Log.Debugf("K8s client config (init InCluster), host: %s", config.Host)

	return &KubeProvider{
		ClientSet:    clientset,
		ClientConfig: *config,
	}, nil
}

func int32Ptr(i int32) *int32 { return &i }
