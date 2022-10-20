package kubernetes

import (
	"context"

	"github.com/mogenius/mo-go/logger"

	apiv1 "k8s.io/api/core/v1"
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
}

func removeDaemonset(kubeProvider *KubeProvider) {
	daemonSetClient := kubeProvider.ClientSet.AppsV1().DaemonSets(apiv1.NamespaceDefault)

	// DELETE DaemonSet
	logger.Log.Info("Deleting podloxx daemonset ...")
	deletePolicy := metav1.DeletePropagationForeground
	err := daemonSetClient.Delete(context.TODO(), DAEMONSETNAME, metav1.DeleteOptions{PropagationPolicy: &deletePolicy})
	if err != nil {
		panic(err)
	}
	logger.Log.Info("Deleted podloxx daemonset.")
}

func removeRedis(kubeProvider *KubeProvider) {
	deploymentClient := kubeProvider.ClientSet.AppsV1().Deployments(apiv1.NamespaceDefault)

	// DELETE REDIS
	logger.Log.Info("Deleting podloxx redis ...")
	deletePolicy := metav1.DeletePropagationForeground
	err := deploymentClient.Delete(context.TODO(), REDISNAME, metav1.DeleteOptions{PropagationPolicy: &deletePolicy})
	if err != nil {
		panic(err)
	}
	logger.Log.Info("Deleted podloxx redis.")
}

func removeRbac(kubeProvider *KubeProvider) {
	// CREATE RBAC
	logger.Log.Info("Deleting podloxx RBAC ...")
	err := kubeProvider.ClientSet.CoreV1().ServiceAccounts(NAMESPACE).Delete(context.TODO(), SERVICEACCOUNTNAME, metav1.DeleteOptions{})
	if err != nil {
		panic(err)
	}
	err = kubeProvider.ClientSet.RbacV1().ClusterRoles().Delete(context.TODO(), CLUSTERROLENAME, metav1.DeleteOptions{})
	if err != nil {
		panic(err)
	}
	err = kubeProvider.ClientSet.RbacV1().ClusterRoleBindings().Delete(context.TODO(), CLUSTERROLEBINDINGNAME, metav1.DeleteOptions{})
	if err != nil {
		panic(err)
	}
	logger.Log.Info("Deleted podloxx RBAC.")
}
