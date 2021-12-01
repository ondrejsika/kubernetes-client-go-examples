package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	batchv1 "k8s.io/api/batch/v1"
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

	jobsClient := clientset.BatchV1().Jobs(namespace)
	podsClient := clientset.CoreV1().Pods(namespace)

	jobName := "go-job"

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: jobName,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: "Never",
					Containers: []corev1.Container{
						{
							Name:  "main",
							Image: "ondrejsika/cowsay",
							Args:  []string{"hello world"},
						},
					},
				},
			},
		},
	}

	jobsClient.Create(context.TODO(), job, metav1.CreateOptions{})

	var pod corev1.Pod

	for {
		pods, _ := podsClient.List(context.TODO(), metav1.ListOptions{
			LabelSelector: "job-name=" + jobName,
		})
		if len(pods.Items) == 0 {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		pod = pods.Items[0]
		break
	}

	fmt.Println(pod.ObjectMeta.Name)

	for {
		pod, _ := podsClient.Get(context.TODO(), pod.ObjectMeta.Name, metav1.GetOptions{})
		if pod.Status.Phase == "Succeeded" {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	logs, _ := podsClient.GetLogs(pod.ObjectMeta.Name, &corev1.PodLogOptions{}).Stream(context.TODO())
	io.Copy(os.Stdout, logs)

	jobsClient.Delete(context.TODO(), jobName, metav1.DeleteOptions{})
	podsClient.DeleteCollection(context.TODO(), metav1.DeleteOptions{}, metav1.ListOptions{
		LabelSelector: "job-name=" + jobName,
	})
}
