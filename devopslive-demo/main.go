package main

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
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
	clientset, namespace, _ := KubernetesClient()

	nodesClient := clientset.CoreV1().Nodes()
	podsClient := clientset.CoreV1().Pods(namespace)
	deploymentsClient := clientset.AppsV1().Deployments(namespace)

	// List Nodes
	fmt.Println("Nodes")
	nodes, _ := nodesClient.List(context.TODO(), metav1.ListOptions{})
	for _, node := range nodes.Items {
		fmt.Println(node.ObjectMeta.Name, node.Status.Addresses[2].Address)
	}

	// Create Deployment
	labels := map[string]string{
		"app": "go-example",
	}

	var replicas int32 = 5

	deployment, _ := deploymentsClient.Create(context.TODO(), &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "go-example",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "main",
							Image: "ondrejsika/go-hello-world:2",
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}, metav1.CreateOptions{})

	time.Sleep(time.Second)

	// List Pods
	fmt.Println("Pods")
	pods, _ := podsClient.List(context.TODO(), metav1.ListOptions{})
	for _, pod := range pods.Items {
		fmt.Println(pod.ObjectMeta.Name, pod.Spec.Containers[0].Image)
	}

	fmt.Println("Apache Pods")
	pods, _ = podsClient.List(context.TODO(), metav1.ListOptions{
		LabelSelector: "app=apache",
	})
	for _, pod := range pods.Items {
		fmt.Println(pod.ObjectMeta.Name, pod.ObjectMeta.Labels)
	}

	// Delete Deployment
	deploymentsClient.Delete(context.TODO(), deployment.ObjectMeta.Name, metav1.DeleteOptions{})

}
