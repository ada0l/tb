package main

import (
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/ada0l/tb/tb"
)

var (
	configPath string
)

func main() {
	flag.StringVar(&configPath, "config", "/etc/tb.yml", "path to config file")
	flag.Parse()

	config, err := tb.LoadConfigFromFile(configPath)

	if err != nil {
		log.Fatalf("Failed to load config file: %v", err)
		os.Exit(1)
	}
	log.Println("Config loaded...")

	roundServerPool := &tb.RoundServerPool{}
	loadBalance := tb.GetLoadBalanceFunction(roundServerPool)

	for _, backend := range config.Backends {
		backendUrl, err := url.Parse(backend)
		if err != nil {
			log.Fatalf("Failed to parse backend url: %v", err)
			os.Exit(1)
		}
		proxy := tb.GetProxy(roundServerPool, loadBalance, backendUrl)
		backend := &tb.Backend{
			URL:          backendUrl,
			Alive:        true,
			ReverseProxy: proxy,
		}
		roundServerPool.AddBackend(backend)
	}

	server := http.Server{
		Addr:    config.Host,
		Handler: http.HandlerFunc(loadBalance),
	}

	go func() {
		t := time.NewTicker(time.Second * 10)
		for {
			<-t.C
			log.Println("Starting health check...")
			roundServerPool.ChechHealth()
			log.Println("Health check completed")
		}
	}()

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
