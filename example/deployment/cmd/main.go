package main

import (
	"flag"
	"github.com/zouchunxu/deployment/internal/svr"
	"github.com/zouchunxu/gof/server"
	"log"
)

func main() {
	var path *string
	path = flag.String("path", "./app.yaml", "config")
	flag.Parse()
	if path == nil {
		log.Fatal("fail")
	}
	app := server.New(*path)
	svr.Init(app)

	if err := app.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
