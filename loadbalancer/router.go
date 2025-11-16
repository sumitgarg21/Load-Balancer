package main

import (
	"net/http"
)

func lbHandler(lb *LoadBalancer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// extract client IP
		clientIP := r.RemoteAddr

		backend := lb.GetBackendSticky(clientIP)

		// set cookie so next time browser will send same backend ID
		http.SetCookie(w, &http.Cookie{
			Name:  "X-Backend",
			Value: backend.URL,
			Path:  "/",
		})

		backend.ReverseProxy.ServeHTTP(w, r)
	}
}
