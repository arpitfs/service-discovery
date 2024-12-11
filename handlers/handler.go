package handler

import (
	"net/http"
	"servicediscovery/clients"
)

func RegisterHandlers() {
	registry := NewServiceRegistry()

	http.HandleFunc("/register", registry.RegisterHandler)
	http.HandleFunc("/deregister", registry.DeRegisterHandler)
	http.HandleFunc("/discover", registry.DiscoverHandler)
	http.HandleFunc("/discoverall", registry.DiscoverAllHandler)
	http.HandleFunc("/client1-health", clients.ClientOneHealthCheck)
	http.HandleFunc("/client2-health", clients.ClientTwoHealthCheck)

	go registry.HealthCheck()
}
