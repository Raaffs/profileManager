package main

import (
	"sync"

	"github.com/labstack/echo/v4"
)

type HealthStatus string

const(
	StatusHealthy   HealthStatus = "healthy"

	//Indicates Some services are unavailable but core functionality is
   	//working
	StatusDegraded HealthStatus = "degraded"

	//Indicates core functionality (ciphers and jwt tokens) is not working as intended
	StatusCritical HealthStatus = "critical"
	
	StatusDown HealthStatus = "down"
)

type HealthChecker struct{
	mu sync.RWMutex
	status HealthStatus
}

func (h *HealthChecker)SetStatus(s HealthStatus){
	h.mu.Lock()
	defer h.mu.Unlock()
	h.status = s
}

func (h *HealthChecker)GetStatus()HealthStatus{
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.status
}

func (h *HealthChecker) Handler(c echo.Context) error {
	status := h.GetStatus()
	var code int

	switch status {
	case StatusHealthy, StatusDegraded:
		code = 200
	case StatusCritical, StatusDown:
		code = 503
	default:
		code = 500
		status = "unknown"
	}
	return c.JSON(code, map[string]string{
		"status": string(status),
	})
}
