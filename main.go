package main

import (
	"flag"
	"log"

	"github.com/cristaloleg/bc/app"
)

var configFile = *flag.String("config", "config.json", "file with a configuration")

func init() {
	log.SetFlags(log.Lshortfile)
	flag.Parse()
}

func main() {
	a := app.New()
	defer a.Stop()
	a.Init(configFile)
	a.Run()
}
