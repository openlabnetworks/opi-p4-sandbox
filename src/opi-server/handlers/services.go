// Package handlers implements the OPI API service handlers for the sandbox.
// These are in-memory emulators that simulate the behavior of a real OPI-compliant
// infrastructure device (SmartNIC/DPU).
package handlers

import (
	"log"
	"sync"
	"time"
)

// --- Data Types (simplified OPI models) ---

// Port represents a virtual network port in the OPI model.
type Port struct {
	Name       string
	Id         string
	MacAddress string
	Mtu        int32
	AdminState string // UP, DOWN
	OperState  string // UP, DOWN, UNKNOWN
	Speed      string
	CreateTime time.Time
}

// Pipeline represents a P4 pipeline configuration.
type Pipeline struct {
	Name       string
	Id         string
	P4Program  string
	Status     string // LOADED, RUNNING, STOPPED, ERROR
	CreateTime time.Time
}

// --- Network Service ---

// NetworkService implements OPI network port management APIs.
type NetworkService struct {
	mu    sync.RWMutex
	ports map[string]*Port
}

// NewNetworkService creates a new NetworkService with a default port.
func NewNetworkService() *NetworkService {
	svc := &NetworkService{
		ports: make(map[string]*Port),
	}

	// Pre-populate with a default port so users see immediate results
	svc.ports["port-0"] = &Port{
		Name:       "ports/port-0",
		Id:         "port-0",
		MacAddress: "00:1a:2b:3c:4d:5e",
		Mtu:        1500,
		AdminState: "UP",
		OperState:  "UP",
		Speed:      "100G",
		CreateTime: time.Now(),
	}

	log.Printf("[NetworkService] Initialized with 1 default port")
	return svc
}

// --- Pipeline Service ---

// PipelineService implements OPI P4 pipeline management APIs.
type PipelineService struct {
	mu        sync.RWMutex
	pipelines map[string]*Pipeline
}

// NewPipelineService creates a new PipelineService.
func NewPipelineService() *PipelineService {
	svc := &PipelineService{
		pipelines: make(map[string]*Pipeline),
	}

	// Pre-populate with a default pipeline
	svc.pipelines["pipeline-default"] = &Pipeline{
		Name:       "pipelines/pipeline-default",
		Id:         "pipeline-default",
		P4Program:  "basic_forwarding.p4",
		Status:     "RUNNING",
		CreateTime: time.Now(),
	}

	log.Printf("[PipelineService] Initialized with 1 default pipeline")
	return svc
}
