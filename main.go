package main

import (
	"bufio"
	"context"
	"embed"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"time"

	"github.com/kubectl-logz/kubectl-logz/internal/parser/logfmt"
	"github.com/kubectl-logz/kubectl-logz/internal/types"

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

	collector, err := internal.NewCollector(kubeconfig)
	if err != nil {
		log.Fatal(err)
	}
	go collector.Run(ctx)

	r := gin.Default()

	r.GET("/api/v1/logs", func(c *gin.Context) {
		dir, err := os.ReadDir("logs")
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		var files []string
		for _, entry := range dir {
			files = append(files, entry.Name())
		}
		c.JSON(http.StatusOK, gin.H{
			"files": files,
		})
	})
	r.GET("/api/v1/logs/:file", func(c *gin.Context) {
		f, err := os.Open(filepath.Join("logs", c.Param("file")))
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		defer f.Close()
		var entries []gin.H
		scanner := bufio.NewScanner(f)
		level := types.Level(c.Query("level"))
		page, _ := strconv.Atoi(c.Query("page"))
		if page < 0 {
			page = 0
		}
		limit, _ := strconv.Atoi(c.Query("limit"))
		if limit <= 0 {
			limit = 100
		}
		count := -1
		for scanner.Scan() {
			entry := types.Entry{}
			err := logfmt.Unmarshal(scanner.Bytes(), &entry)
			if err != nil {
				c.JSON(http.StatusInternalServerError, err.Error())
				return
			}
			if entry.Level.Less(level) {
				continue
			}
			count++
			if count < page*limit {
				continue
			}
			if len(entries) <= limit {
				h := gin.H{
					"level": entry.Level,
					"msg":   entry.Msg,
				}
				if !entry.Time.IsZero() {
					h["time"] = entry.Time
				}
				entries = append(entries, h)
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"entries": entries,
			"metadata": gin.H{
				"count": count,
			},
		})
	})

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
	srv := &http.Server{Addr: ":5649", Handler: r}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
	<-ctx.Done()
}
