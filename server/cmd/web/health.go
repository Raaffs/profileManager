package main

import (
	"sync"

	"github.com/labstack/echo/v4"
)

type HealthStatus string

const(
	StatusHealthy   HealthStatus = "healthy"
	StatusDegraded HealthStatus = "degraded"
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

func (h *HealthChecker)Handler(c echo.Context)error{
	switch h.GetStatus(){
	case StatusHealthy:
		return c.String(200,string(StatusHealthy))
	case StatusDegraded:
		return c.String(200,string(StatusDegraded))
	case StatusCritical:
		return c.String(503,string(StatusCritical))
	case StatusDown:
		return c.String(503,string(StatusDown))
	default:
		return c.String(500,"unknown")
	}	
}