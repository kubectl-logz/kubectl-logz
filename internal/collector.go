package internal

import (
	"bufio"
	"context"
	"log"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
)

func Run(ctx context.Context, kubeconfig string) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	for {
		err := collectPods(ctx, clientset)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func collectPods(ctx context.Context, clientset *kubernetes.Clientset) error {
	pods, err := clientset.CoreV1().Pods("").Watch(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}
	defer pods.Stop()
	for event := range pods.ResultChan() {
		pod := event.Object.(*corev1.Pod)
		log.Printf("type=%v, pod=%s\n", event.Type, pod.Name)
		for _, container := range pod.Spec.Containers {
			err := collectContainerLogs(ctx, clientset, pod.Namespace, pod.Name, container.Name)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func collectContainerLogs(ctx context.Context, clientset *kubernetes.Clientset, namespace, pod, container string) error {
	log.Printf("pod=%s, container=%s\n", pod, container)
	logs, err := clientset.CoreV1().Pods(namespace).GetLogs(pod, &corev1.PodLogOptions{
		Container:  container,
		Follow:     true,
		Timestamps: true,
	}).Stream(ctx)
	if err != nil {
		return err
	}
	defer logs.Close()

	scanner := bufio.NewScanner(logs)

	for scanner.Scan() {
		r := parse(scanner.Bytes())
		if r.Valid() {
			log.Println(r.String())
		}
	}

	return nil
}
