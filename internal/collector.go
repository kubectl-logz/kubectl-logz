package internal

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/kubectl-logz/kubectl-logz/internal/types"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
)

func init() {
	_ = os.Mkdir("logs", 0700)
}

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
			go func(container string) {
				err := collectContainerLogs(ctx, clientset, pod.Namespace, pod.Name, container)
				if err != nil {
					log.Printf("err=%q\n", err)
				}
			}(container.Name)
		}
	}
	return nil
}

func collectContainerLogs(ctx context.Context, clientset *kubernetes.Clientset, namespace, pod, container string) error {
	log.Printf("pod=%s, container=%s\n", pod, container)
	logs, err := clientset.CoreV1().Pods(namespace).GetLogs(pod, &corev1.PodLogOptions{
		Container: container,
		Follow:    true,
	}).Stream(ctx)
	if err != nil {
		return err
	}
	defer logs.Close()

	scanner := bufio.NewScanner(logs)

	c := types.Ctx{Host: fmt.Sprintf("%s.%s.%s", namespace, pod, container)}

	f, err := os.Create(filepath.Join("logs", c.Host+".log"))
	if err != nil {
		return err
	}
	defer f.Close()

	for scanner.Scan() {
		r := parse(scanner.Bytes())
		if _, err := f.WriteString(r.String() + "\n"); err != nil {
			return err
		}
	}

	return nil
}
