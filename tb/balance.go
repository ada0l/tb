package tb

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type LoadBalanceFunc http.HandlerFunc

func GetLoadBalanceFunction(s ServerPool) LoadBalanceFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		attemps := GetRetryFromContext(request)
		if attemps > 3 {
			log.Printf("%s(%s) Max attempts reached, terminating\n", request.RemoteAddr, request.URL.Path)
			http.Error(writer, "Service not available", http.StatusServiceUnavailable)
			return
		}

		peer := s.GetNextPeer()
		if peer != nil {
			peer.ReverseProxy.ServeHTTP(writer, request)
			return
		}
		http.Error(writer, "Server not available", http.StatusServiceUnavailable)
	}
}

func GetProxy(s ServerPool, loadBalanceFunction LoadBalanceFunc, backendURL *url.URL) *httputil.ReverseProxy {
	proxy := httputil.NewSingleHostReverseProxy(backendURL)
	proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, err error) {
		log.Printf("[%s] %s\n", request.URL.Host, err.Error())
		retries := GetRetryFromContext(request)

		if retries < 3 {
			<-time.After(10 * time.Millisecond)
			ctx := SetRetryForContext(request, retries+1)
			proxy.ServeHTTP(writer, request.WithContext(ctx))
			return
		}

		s.MarkBackendStatus(backendURL, false)

		attempts := GetAttemptsFromContext(request)
		log.Printf("%s(%s) Attempting retry %d\n", request.RemoteAddr, request.URL.Path, attempts)
		ctx := SetAttemptsForContext(request, attempts+1)
		loadBalanceFunction(writer, request.WithContext(ctx))
	}

	return proxy
}
