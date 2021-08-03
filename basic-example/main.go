package main

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)
	namespace, _, _ := kubeConfig.Namespace()
	restconfig, _ := kubeConfig.ClientConfig()
	clientset, _ := kubernetes.NewForConfig(restconfig)

	podsClient := clientset.CoreV1().Pods(namespace)
	pods, _ := podsClient.List(context.TODO(), metav1.ListOptions{})
	for _, pod := range pods.Items {
		fmt.Printf("name=%s image=%s\n", pod.Name, pod.Spec.Containers[0].Image)
	}
}
