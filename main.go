package main

import (
	"context"
	"embed"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/kubectl-logz/kubectl-logz/internal"
	"k8s.io/client-go/util/homedir"

	"github.com/gin-gonic/gin"
	"github.com/pkg/browser"
)

//go:generate npm install
//go:generate npm run build

//go:embed build
var fs embed.FS

func main() {
	var openBrowser bool
	var kubeconfig string
	flag.BoolVar(&openBrowser, "b", false, "open browser")
	flag.StringVar(&kubeconfig, "k", filepath.Join(homedir.HomeDir(), ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	flag.Parse()
	log.Printf(" openBrowser=%v\n", openBrowser)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go internal.Run(ctx, kubeconfig)

	r := gin.Default()

	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		log.Printf("GET %q\n", path)
		c.FileFromFS("build/"+path, http.FS(fs))
	})
	go func() {
		time.Sleep(time.Second)
		if openBrowser {
			if err := browser.OpenURL("http://localhost:5649"); err != nil {
				log.Fatal(err)
			}
		}
	}()
	srv := &http.Server{Addr: "localhost:5649", Handler: r}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
	<-ctx.Done()
}
