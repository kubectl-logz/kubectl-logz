package main

import (
	"embed"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/browser"
)

//go:generate npm install
//go:generate npm run build

//go:embed build
var fs embed.FS

func main() {
	var openBrowser bool
	flag.BoolVar(&openBrowser, "b", true, "open browser")
	flag.Parse()
	log.Printf(" openBrowser=%v\n", openBrowser)

	r := gin.Default()

	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		log.Printf("GET %q\n", path)
		c.FileFromFS("build/"+path, http.FS(fs))
	})
	go func() {
		time.Sleep(time.Second)
		if err := browser.OpenURL("http://localhost:8080"); err != nil {
			log.Fatal(err)
		}
	}()
	if err := r.Run("localhost:8080"); err != nil {
		log.Fatal(err)
	}
}
