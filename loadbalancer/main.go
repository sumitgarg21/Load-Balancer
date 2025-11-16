package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"os"
	"strings"
	"github.com/joho/godotenv"
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

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	backendEnv := os.Getenv("BACKENDS")
	port := os.Getenv("PORT")

	if backendEnv == "" {
		log.Fatal("ERROR: BACKENDS not set in .env file")
	}
	if port == "" {
		port = "8080" // default
	}

	// Split comma separated list
	backendURLs := strings.Split(backendEnv, ",")

	// Create backend objects
	backends := []*Backend{}
	for _, url := range backendURLs {
		backends = append(backends, createBackend(url))
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

	log.Println("Load Balancer running on", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
