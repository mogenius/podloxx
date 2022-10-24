package kubernetes

import (
	"context"

	"github.com/mogenius/mo-go/logger"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Remove() {
	provider, err := NewKubeProviderLocal()
	if err != nil {
		panic(err)
	}

	removeDaemonset(provider)
	removeRbac(provider)
	removeRedis(provider)
	removeRedisService(provider)
}

func removeDaemonset(kubeProvider *KubeProvider) {
	daemonSetClient := kubeProvider.ClientSet.AppsV1().DaemonSets(NAMESPACE)

	// DELETE DaemonSet
	logger.Log.Info("Deleting podloxx daemonset ...")
	deletePolicy := metav1.DeletePropagationForeground
	err := daemonSetClient.Delete(context.TODO(), DAEMONSETNAME, metav1.DeleteOptions{PropagationPolicy: &deletePolicy})
	if err != nil {
		logger.Log.Error(err)
	}
	logger.Log.Info("Deleted podloxx daemonset.")
}

func removeRedis(kubeProvider *KubeProvider) {
	deploymentClient := kubeProvider.ClientSet.AppsV1().Deployments(NAMESPACE)

	// DELETE REDIS
	logger.Log.Info("Deleting podloxx redis ...")
	deletePolicy := metav1.DeletePropagationForeground
	err := deploymentClient.Delete(context.TODO(), REDISNAME, metav1.DeleteOptions{PropagationPolicy: &deletePolicy})
	if err != nil {
		logger.Log.Error(err)
	}
	logger.Log.Info("Deleted podloxx redis.")
}

func removeRedisService(kubeProvider *KubeProvider) {
	serviceClient := kubeProvider.ClientSet.CoreV1().Services(NAMESPACE)

	// DELETE REDIS
	logger.Log.Info("Deleting podloxx redis service ...")
	deletePolicy := metav1.DeletePropagationForeground
	err := serviceClient.Delete(context.TODO(), REDISSERVICENAME, metav1.DeleteOptions{PropagationPolicy: &deletePolicy})
	if err != nil {
		logger.Log.Error(err)
	}
	logger.Log.Info("Deleted podloxx redis service.")
}

func removeRbac(kubeProvider *KubeProvider) {
	// CREATE RBAC
	logger.Log.Info("Deleting podloxx RBAC ...")
	err := kubeProvider.ClientSet.CoreV1().ServiceAccounts(NAMESPACE).Delete(context.TODO(), SERVICEACCOUNTNAME, metav1.DeleteOptions{})
	if err != nil {
		logger.Log.Error(err)
	}
	err = kubeProvider.ClientSet.RbacV1().ClusterRoles().Delete(context.TODO(), CLUSTERROLENAME, metav1.DeleteOptions{})
	if err != nil {
		logger.Log.Error(err)
	}
	err = kubeProvider.ClientSet.RbacV1().ClusterRoleBindings().Delete(context.TODO(), CLUSTERROLEBINDINGNAME, metav1.DeleteOptions{})
	if err != nil {
		logger.Log.Error(err)
	}
	logger.Log.Info("Deleted podloxx RBAC.")
}
