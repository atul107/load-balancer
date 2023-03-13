package main

import (
	"net/http"
)

// LoadBalancer is a load balancer that selects the next server based on the selected algorithm
type LoadBalancer struct {
	serverList *ServerList
	algorithm  string
}

// NewLoadBalancer creates a new LoadBalancer instance
func NewLoadBalancer(serverList *ServerList, algorithm string) *LoadBalancer {
	return &LoadBalancer{serverList, algorithm}
}

// HandleRequest handles the incoming requests by selecting the next server based on the selected algorithm
func (lb *LoadBalancer) HandleRequest(w http.ResponseWriter, r *http.Request) {
	var serverURL string
	if lb.algorithm == "round-robin" {
		serverURL = lb.serverList.NextServerRoundRobin()
	} else if lb.algorithm == "intelligent" {
		serverURL = lb.serverList.NextServerIntelligent(r.RemoteAddr)
	} else {
		http.Error(w, "Invalid algorithm", http.StatusBadRequest)
		return
	}

	if serverURL == "" {
		http.Error(w, "No available servers", http.StatusServiceUnavailable)
		return
	}

	// Modify the request to forward it to the selected server
	r.URL.Scheme = "http"
	r.URL.Host = serverURL

	// Create a new request with the modified URL
	req, err := http.NewRequest(r.Method, r.URL.String(), r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header = r.Header

	// Forward the request to the selected server and copy the response to the original response writer
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	for k, v := range resp.Header {
		w.Header()[k] = v
	}

	w.WriteHeader(resp.StatusCode)
	resp.Write(w)
}
