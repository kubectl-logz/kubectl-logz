package internal

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/kubectl-logz/kubectl-logz/internal/parser"
	"github.com/kubectl-logz/kubectl-logz/internal/types"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
)

func init() {
	_ = os.Mkdir("logs", 0700)
}

type Collector struct {
	clientset *kubernetes.Clientset
	namespace string
}

func NewCollector(kubeconfig string) (*Collector, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	loadingRules.DefaultClientConfig = &clientcmd.DefaultClientConfig
	loadingRules.ExplicitPath = kubeconfig
	clientConfig := clientcmd.NewInteractiveDeferredLoadingClientConfig(loadingRules, &clientcmd.ConfigOverrides{}, os.Stdin)
	namespace, _, err := clientConfig.Namespace()
	if err != nil {
		return nil, err
	}
	config, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &Collector{
		clientset: clientset,
		namespace: namespace,
	}, nil
}

func (c *Collector) Run(ctx context.Context) {
	for {
		err := c.collectPods(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (c *Collector) collectPods(ctx context.Context) error {
	pods, err := c.clientset.CoreV1().Pods(c.namespace).Watch(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}
	defer pods.Stop()
	for event := range pods.ResultChan() {
		pod, ok := event.Object.(*corev1.Pod)
		if !ok {
			return fmt.Errorf("%v", event.Object)
		}
		if event.Type != watch.Added {
			continue
		}
		for _, container := range pod.Spec.Containers {
			go func(namespace, pod, container string) {
				defer runtime.HandleCrash()
				err := c.collectContainerLogs(ctx, pod, container)
				if err != nil {
					log.Printf("err=%q\n", err)
				}
			}(pod.Namespace, pod.Name, container.Name)
		}
	}
	return nil
}

func (c *Collector) collectContainerLogs(ctx context.Context, pod, container string) error {
	log.Printf("pod=%s, container=%s\n", pod, container)
	logs, err := c.clientset.CoreV1().Pods(c.namespace).GetLogs(pod, &corev1.PodLogOptions{Container: container}).Stream(ctx)
	if err != nil {
		return err
	}
	defer logs.Close()

	lc := types.Ctx{Hostname: fmt.Sprintf("%s.%s.%s", c.namespace, pod, container)}
	log.Println(lc.String())

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
			f, err := os.Create(filepath.Join("logs", lc.Hostname+".log"))
			if err != nil {
				return err
			}
			defer f.Close()
			var offset int64 = 0
			for entry := range entries {
				if n, err := f.WriteString(entry.String() + "\n"); err != nil {
					return err
				} else {
					offset = offset + int64(n)
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
