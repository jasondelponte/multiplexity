package main

import (
	"flag"
	"github.com/jasondelponte/multiplexity"
	"log"
	"runtime"
)

var (
	cfgFileName = flag.String("config", "config.json", "Configuration file")
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	config, err := multiplexity.LoadConfig(*cfgFileName)
	if err != nil {
		log.Fatalf("Failed to open load configuration %s, error %s", *cfgFileName, err)
		return
	}

	proxy := multiplexity.NewProxy(config)
	proxy.Start()
}
