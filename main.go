package main

import (
	"fmt"
	"net/http"
)

func main() {
	// Create a server list and add some servers
	serverList := NewServerList()
	serverList.AddServer("http://localhost:8001")
	serverList.AddServer("http://localhost:8002")
	serverList.AddServer("http://localhost:8003")

	// Create a load balancer with round-robin algorithm
	loadBalancerRR := NewLoadBalancer(serverList, "round-robin")

	// Create a load balancer with intelligent algorithm
	loadBalancerIntelligent := NewLoadBalancer(serverList, "intelligent")

	// Register the load balancers with the HTTP server
	http.HandleFunc("/rr", loadBalancerRR.HandleRequest)
	http.HandleFunc("/intelligent", loadBalancerIntelligent.HandleRequest)

	// Start the HTTP server
	fmt.Println("Load balancers listening on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
