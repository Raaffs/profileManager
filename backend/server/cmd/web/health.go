package main

import (
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
)

type HealthStatus int



const(
	StatusHealthy   HealthStatus = iota

	//Indicates Some services are unavailable but core functionality is
   	//working
	StatusDegraded 

	//Indicates core functionality (ciphers and jwt tokens) is not working as intended
	StatusCritical  

	StatusDown 
)

func (s HealthStatus) String() string {
	switch s {
	case StatusHealthy:
		return "healthy"
	case StatusDegraded:
		return "degraded"
	case StatusCritical:
		return "critical"
	case StatusDown:
		return "down"
	default:
		return "unknown"
	}
}


type HealthChecker struct{
	mu sync.RWMutex
	status HealthStatus
}

func (h *HealthChecker)SetStatus(s HealthStatus){
	h.mu.Lock()
	defer h.mu.Unlock()
	if s<h.status{
		return
	}
	h.status = s
}

func (h *HealthChecker)GetStatus()HealthStatus{
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.status
}

// ResetStatus forces the status back to Healthy, bypassing the escalation check.
func (h *HealthChecker) ResetStatus() {
    h.mu.Lock()
    defer h.mu.Unlock()
    h.status = StatusHealthy
}
func (h *HealthChecker) Handler(c echo.Context) error {
	status := h.GetStatus()
	var code int

	// Logic for HTTP Codes
	switch status {
	case StatusHealthy, StatusDegraded:
		code = http.StatusOK
	case StatusCritical, StatusDown:
		code = http.StatusServiceUnavailable
	default:
		code = http.StatusInternalServerError
	}

	return c.JSON(code, map[string]string{
		"status": status.String(), // Uses the method we created above
	})
}