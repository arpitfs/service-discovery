package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"servicediscovery/model"
	"sync"
	"time"
)

type ServiceRegistry struct {
	mu       sync.RWMutex
	services map[string]model.Service
}

func NewServiceRegistry() *ServiceRegistry {
	return &ServiceRegistry{
		services: make(map[string]model.Service),
	}
}

func (s *ServiceRegistry) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var service model.Service
	if err := json.NewDecoder(r.Body).Decode(&service); err != nil {
		http.Error(w, "Invalid Input", http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.services[service.Name] = service
	w.WriteHeader(http.StatusOK)
	fmt.Println("Service Registered: ", service)
}

func (s *ServiceRegistry) DeRegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	serviceName := r.URL.Query().Get("name")
	if serviceName == "" {
		http.Error(w, "Missing service name", http.StatusBadRequest)
		return
	}
	var service model.Service
	if err := json.NewDecoder(r.Body).Decode(&service); err != nil {
		http.Error(w, "Invalid Input", http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.services, serviceName)
	fmt.Println("Service deregistered: ", serviceName)
}

func (s *ServiceRegistry) DiscoverHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Missing service name", http.StatusBadRequest)
		return
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	service, found := s.services[name]
	if !found {
		http.Error(w, "Service not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(service)
}

func (s *ServiceRegistry) DiscoverAllHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	var services []model.Service
	for _, service := range s.services {
		services = append(services, service)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services)
}

func (s *ServiceRegistry) HealthCheck() {
	for {
		s.mu.RLock()
		for _, service := range s.services {
			resp, err := http.Get(service.Url)
			if err != nil || resp.StatusCode != http.StatusOK {
				fmt.Printf("Service not healthy: %s\n", service.Name)
			} else {
				fmt.Printf("Service healthy: %s\n", service.Name)
			}
			if resp != nil {
				resp.Body.Close()
			}
		}
		s.mu.RUnlock()
		time.Sleep(10 * time.Second)
	}
}
