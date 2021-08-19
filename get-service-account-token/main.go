package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func KubernetesClient() (*kubernetes.Clientset, string, error) {
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)
	namespace, _, _ := kubeConfig.Namespace()
	restconfig, _ := kubeConfig.ClientConfig()
	clientset, _ := kubernetes.NewForConfig(restconfig)
	return clientset, namespace, nil
}

func main() {
	clientset, defaultNamespace, _ := KubernetesClient()

	serviceAccountName := flag.String("service-account", "default", "")
	namespace := flag.String("namespace", defaultNamespace, "")

	flag.Parse()

	saClient := clientset.CoreV1().ServiceAccounts(*namespace)
	secretClient := clientset.CoreV1().Secrets(*namespace)

	sa, err := saClient.Get(context.TODO(), *serviceAccountName, metav1.GetOptions{})
	if err != nil {
		log.Fatal(err)
	}
	secret, err := secretClient.Get(context.TODO(), sa.Secrets[0].Name, metav1.GetOptions{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(secret.Data["token"]))
}
