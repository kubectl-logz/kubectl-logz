package internal

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"k8s.io/apimachinery/pkg/watch"

	"github.com/kubectl-logz/kubectl-logz/internal/parser"
	"github.com/kubectl-logz/kubectl-logz/internal/types"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
)

func init() {
	_ = os.Mkdir("logs", 0700)
}

func Run(ctx context.Context, kubeconfig string) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	loadingRules.DefaultClientConfig = &clientcmd.DefaultClientConfig
	loadingRules.ExplicitPath = kubeconfig
	clientConfig := clientcmd.NewInteractiveDeferredLoadingClientConfig(loadingRules, &clientcmd.ConfigOverrides{}, os.Stdin)
	namespace, _, err := clientConfig.Namespace()
	if err != nil {
		log.Fatal(err)
	}
	config, err := clientConfig.ClientConfig()
	if err != nil {
		log.Fatal(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	for {
		err := collectPods(ctx, namespace, clientset)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func collectPods(ctx context.Context, namespace string, clientset *kubernetes.Clientset) error {
	pods, err := clientset.CoreV1().Pods(namespace).Watch(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}
	defer pods.Stop()
	for event := range pods.ResultChan() {
		pod := event.Object.(*corev1.Pod)
		log.Printf("type=%v, pod=%s\n", event.Type, pod.Name)
		if event.Type != watch.Added {
			continue
		}
		for _, container := range pod.Spec.Containers {
			go func(namespace, pod, container string) {
				defer runtime.HandleCrash()
				err := collectContainerLogs(ctx, clientset, namespace, pod, container)
				if err != nil {
					log.Printf("err=%q\n", err)
				}
			}(pod.Namespace, pod.Name, container.Name)
		}
	}
	return nil
}

func collectContainerLogs(ctx context.Context, clientset *kubernetes.Clientset, namespace, pod, container string) error {
	log.Printf("pod=%s, container=%s\n", pod, container)
	logs, err := clientset.CoreV1().Pods(namespace).GetLogs(pod, &corev1.PodLogOptions{Container: container}).Stream(ctx)
	if err != nil {
		return err
	}
	defer logs.Close()

	c := types.Ctx{Host: fmt.Sprintf("%s.%s.%s", namespace, pod, container)}

	lines := make(chan []byte, 100)
	defer close(lines)

	entries := make(chan types.Entry, 10)
	errors := make(chan error, 1)
	defer close(errors)

	go func() {
		defer runtime.HandleCrash()
		defer close(entries)
		parser.Parse(lines, entries)
	}()
	go func() {
		defer runtime.HandleCrash()
		err := func() error {
			f, err := os.Create(filepath.Join("logs", c.Host+".log"))
			if err != nil {
				return err
			}
			defer f.Close()
			for entry := range entries {
				if _, err := f.WriteString(entry.String() + "\n"); err != nil {
					return err
				}
			}
			return nil
		}()
		if err != nil {
			errors <- err
		}
	}()

	scanner := bufio.NewScanner(logs)
	for scanner.Scan() {
		select {
		case err := <-errors:
			return err
		default:
			// make a copy because the scanner re-uses it's array
			i := scanner.Bytes()
			bytes := make([]byte, len(i))
			copy(bytes, i)
			lines <- bytes
		}
	}

	return nil
}
