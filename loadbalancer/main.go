package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func createBackend(rawURL string) *Backend {
	u, _ := url.Parse(rawURL)
	return &Backend{
		URL:          rawURL,
		Alive:        true,
		ReverseProxy: httputil.NewSingleHostReverseProxy(u),
	}
}

func main() {
	InitMetrics() 	

	backends := []*Backend{
		createBackend("http://localhost:8081"),
		createBackend("http://localhost:8082"),
		createBackend("http://localhost:8083"),
	}

	lb := &LoadBalancer{Backends: backends}

	// Start health checker
	go func() {
		for {
			aliveCount := 0
			for _, b := range lb.Backends {
				isAlive := healthCheck(b.URL)
				b.SetAlive(isAlive)
				if isAlive {
					aliveCount++
				}
			}
			ActiveBackends.Set(float64(aliveCount))
			time.Sleep(5 * time.Second)
		}
	}()

	rl := NewRateLimiter(10, 1*time.Second) // 10 req/sec per IP

	http.Handle("/metrics", promhttp.Handler())
	http.Handle("/", rl.Limit(http.HandlerFunc(lbHandler(lb))))

	log.Println("Load Balancer running on :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
