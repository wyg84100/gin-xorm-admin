package main

import (
	"flag"

	"github.com/gin-gonic/gin"

	_ "./db"
	"./router"
)

func main() {
	var port, mode string
	flag.StringVar(&port, "port", "3000", "service listening at, default 3000")
	flag.StringVar(&mode, "mode", "debug", "service running mode, default debug mode")

	flag.Parse()

	gin.SetMode(mode)
	router.Init(port)
}
