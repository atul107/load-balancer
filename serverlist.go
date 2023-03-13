package main

import (
	"math/rand"
	"net"
	"net/url"
	"strings"
	"sync"
	"time"
)

// ServerList is a list of servers
type ServerList struct {
	servers      []string
	currentIndex int
	mux          sync.Mutex
}

// NewServerList creates a new ServerList instance
func NewServerList() *ServerList {
	return &ServerList{servers: []string{}}
}

// AddServer adds a new server to the list
func (sl *ServerList) AddServer(serverURL string) {
	sl.servers = append(sl.servers, serverURL)
}

// NextServerRoundRobin selects the next server in the list using round-robin algorithm
func (sl *ServerList) NextServerRoundRobin() string {
	sl.mux.Lock()
	defer sl.mux.Unlock()

	if len(sl.servers) == 0 {
		return ""
	}

	serverURL := sl.servers[sl.currentIndex]
	sl.currentIndex = (sl.currentIndex + 1) % len(sl.servers)

	return serverURL
}

// NextServerIntelligent selects the next server in the list using intelligent load balancing algorithm
func (sl *ServerList) NextServerIntelligent(clientAddr string) string {
	sl.mux.Lock()
	defer sl.mux.Unlock()

	if len(sl.servers) == 0 {
		return ""
	}

	// Get the IP address of the client
	ip, _, err := net.SplitHostPort(clientAddr)
	if err != nil {
		return sl.NextServerRoundRobin()
	}

	// Find the server with the closest IP address to the client
	var closestServer string
	closestDistance := -1

	for _, serverURL := range sl.servers {
		u, err := url.Parse(serverURL)
		if err != nil {
			continue
		}

		host, _, err := net.SplitHostPort(u.Host)
		if err != nil {
			continue
		}

		distance := ipDistance(ip, host)
		if closestDistance == -1 || distance < closestDistance {
			closestDistance = distance
			closestServer = serverURL
		}
	}

	return closestServer
}

// ipDistance calculates the distance between two IP addresses
func ipDistance(ip1 string, ip2 string) int {
	ip1Parts := strings.Split(ip1, ".")
	ip2Parts := strings.Split(ip2, ".")

	distance := 0
	for i := 0; i < 4; i++ {
		n1 := parseInt(ip1Parts[i])
		n2 := parseInt(ip2Parts[i])

		if n1 == n2 {
			continue
		}

		distance += (n1 - n2) * (n1 - n2)
	}

	return distance
}

// parseInt converts a string to an integer
func parseInt(s string) int {
	n := 0
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return -1
		}
		n = n*10 + int(s[i]-'0')
	}
	return n
}

// Shuffle randomizes the order of the servers in the list
func (sl *ServerList) Shuffle() {
	sl.mux.Lock()
	defer sl.mux.Unlock()

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(sl.servers), func(i, j int) {
		sl.servers[i], sl.servers[j] = sl.servers[j], sl.servers[i]
	})
}
