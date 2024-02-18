package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/ada0l/tb/tb"
)

func main() {
	port := 3200
	addr := fmt.Sprintf("localhost:%d", port)
	urls := make([]*url.URL, 0)

	for i := 1; i < 4; i++ {
		backendUrl, err := url.Parse(fmt.Sprintf("http://localhost:%d", 3200+i))
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		urls = append(urls, backendUrl)
	}

	roundServerPool := &tb.RoundServerPool{}
	loadBalance := tb.GetLoadBalanceFunction(roundServerPool)

	for _, backendUrl := range urls {
		proxy := tb.GetProxy(roundServerPool, loadBalance, backendUrl)
		backend := &tb.Backend{
			URL:          backendUrl,
			Alive:        true,
			ReverseProxy: proxy,
		}
		roundServerPool.AddBackend(backend)
	}

	server := http.Server{
		Addr:    addr,
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
