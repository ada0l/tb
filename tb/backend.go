package tb

import (
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

type Backend struct {
	URL          *url.URL
	Alive        bool
	mux          sync.RWMutex
	ReverseProxy *httputil.ReverseProxy
}

func NewBackend(backendURL *url.URL) *Backend {
	proxy := httputil.NewSingleHostReverseProxy(backendURL)
	return &Backend{URL: backendURL, Alive: true, ReverseProxy: proxy}
}

func (b *Backend) SetProxyErrorHandler(errorHandler func(writer http.ResponseWriter, request *http.Request, err error)) {
	b.ReverseProxy.ErrorHandler = errorHandler
}

func (b *Backend) SetAlive(alive bool) {
	b.mux.Lock()
	b.Alive = alive
	b.mux.Unlock()
}

func (b *Backend) IsAlive() (alive bool) {
	b.mux.RLock()
	alive = b.Alive
	b.mux.RUnlock()
	return
}

func (b *Backend) CheckHealth() (alive bool, err error) {
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", b.URL.Host, timeout)
	if err != nil {
		log.Println("Site unreachable, error: ", err)
		return false, err
	}
	_ = conn.Close()
	return true, nil
}
