package main

import (
	"flag"
	"github.com/zouchunxu/deployment/internal/svr"
	"github.com/zouchunxu/gof"
	"log"
)

func main() {
	var path *string
	path = flag.String("path", "./app.yaml", "config")
	flag.Parse()
	if path == nil {
		log.Fatal("fail")
	}
	app := gof.New(*path)
	svr.Init(app)
	app.Log.Info("start hhhhhhhhhh \n")
	if err := app.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
