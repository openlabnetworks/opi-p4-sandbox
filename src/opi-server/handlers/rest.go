// Package handlers provides HTTP/JSON REST handlers that emulate OPI APIs.
// This file implements the REST API server that runs alongside or instead of
// raw gRPC, making the sandbox accessible via curl without needing grpcurl.
package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// RegisterRESTRoutes registers all OPI-style REST API routes on the given mux.
func RegisterRESTRoutes(mux *http.ServeMux, netSvc *NetworkService, pipeSvc *PipelineService) {
	// API info
	mux.HandleFunc("/", handleAPIInfo)
	mux.HandleFunc("/healthz", handleHealthz)

	// Network service
	mux.HandleFunc("/v1/ports", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleListPorts(w, r, netSvc)
		case http.MethodPost:
			handleCreatePort(w, r, netSvc)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/v1/ports/", func(w http.ResponseWriter, r *http.Request) {
		portId := strings.TrimPrefix(r.URL.Path, "/v1/ports/")
		if portId == "" {
			http.Error(w, "port_id required", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			handleGetPort(w, r, netSvc, portId)
		case http.MethodDelete:
			handleDeletePort(w, r, netSvc, portId)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Pipeline service
	mux.HandleFunc("/v1/pipelines", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleListPipelines(w, r, pipeSvc)
		case http.MethodPost:
			handleCreatePipeline(w, r, pipeSvc)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/v1/pipelines/", func(w http.ResponseWriter, r *http.Request) {
		pipelineId := strings.TrimPrefix(r.URL.Path, "/v1/pipelines/")
		if pipelineId == "" {
			http.Error(w, "pipeline_id required", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			handleGetPipeline(w, r, pipeSvc, pipelineId)
		case http.MethodDelete:
			handleDeletePipeline(w, r, pipeSvc, pipelineId)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

// --- API Info ---

func handleAPIInfo(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	writeJSON(w, http.StatusOK, APIInfoResponse{
		Name:    "OPI P4 Sandbox API Server",
		Version: "0.1.0",
		Services: []string{
			"opi.v1.NetworkService",
			"opi.v1.PipelineService",
		},
		DocsURL: "https://github.com/opi-p4-sandbox/opi-p4-sandbox",
	})
}

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, HealthResponse{
		Status: "SERVING",
		Services: map[string]string{
			"opi.v1.NetworkService":  "SERVING",
			"opi.v1.PipelineService": "SERVING",
		},
	})
}

// --- Network Service Handlers ---

func handleListPorts(w http.ResponseWriter, r *http.Request, svc *NetworkService) {
	svc.mu.RLock()
	defer svc.mu.RUnlock()

	ports := make([]*PortInfo, 0, len(svc.ports))
	for _, p := range svc.ports {
		ports = append(ports, &PortInfo{
			Name:       p.Name,
			Id:         p.Id,
			MacAddress: p.MacAddress,
			Mtu:        p.Mtu,
			AdminState: p.AdminState,
			OperState:  p.OperState,
			Speed:      p.Speed,
			CreateTime: timestamppb.New(p.CreateTime),
		})
	}

	writeJSON(w, http.StatusOK, PortListResponse{Ports: ports})
}

func handleGetPort(w http.ResponseWriter, r *http.Request, svc *NetworkService, portId string) {
	svc.mu.RLock()
	defer svc.mu.RUnlock()

	p, ok := svc.ports[portId]
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": fmt.Sprintf("port %q not found", portId)})
		return
	}

	writeJSON(w, http.StatusOK, PortInfo{
		Name:       p.Name,
		Id:         p.Id,
		MacAddress: p.MacAddress,
		Mtu:        p.Mtu,
		AdminState: p.AdminState,
		OperState:  p.OperState,
		Speed:      p.Speed,
		CreateTime: timestamppb.New(p.CreateTime),
	})
}

func handleCreatePort(w http.ResponseWriter, r *http.Request, svc *NetworkService) {
	var req CreatePortRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Allow empty body — use defaults
		req = CreatePortRequest{}
	}

	svc.mu.Lock()
	defer svc.mu.Unlock()

	id := fmt.Sprintf("port-%d", len(svc.ports))
	name := fmt.Sprintf("ports/%s", id)

	if req.MacAddress == "" {
		req.MacAddress = fmt.Sprintf("00:1a:2b:3c:%02x:%02x", len(svc.ports)/256, len(svc.ports)%256)
	}
	if req.Mtu == 0 {
		req.Mtu = 1500
	}

	port := &Port{
		Name:       name,
		Id:         id,
		MacAddress: req.MacAddress,
		Mtu:        req.Mtu,
		AdminState: "UP",
		OperState:  "UP",
		Speed:      "25G",
		CreateTime: time.Now(),
	}

	svc.ports[id] = port
	log.Printf("[REST] Created port: %s (MAC: %s)", id, port.MacAddress)

	writeJSON(w, http.StatusCreated, PortInfo{
		Name:       port.Name,
		Id:         port.Id,
		MacAddress: port.MacAddress,
		Mtu:        port.Mtu,
		AdminState: port.AdminState,
		OperState:  port.OperState,
		Speed:      port.Speed,
		CreateTime: timestamppb.New(port.CreateTime),
	})
}

func handleDeletePort(w http.ResponseWriter, r *http.Request, svc *NetworkService, portId string) {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	if _, ok := svc.ports[portId]; !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": fmt.Sprintf("port %q not found", portId)})
		return
	}

	delete(svc.ports, portId)
	log.Printf("[REST] Deleted port: %s", portId)

	w.WriteHeader(http.StatusNoContent)
}

// --- Pipeline Service Handlers ---

func handleListPipelines(w http.ResponseWriter, r *http.Request, svc *PipelineService) {
	svc.mu.RLock()
	defer svc.mu.RUnlock()

	pipelines := make([]*PipelineInfo, 0, len(svc.pipelines))
	for _, p := range svc.pipelines {
		pipelines = append(pipelines, &PipelineInfo{
			Name:       p.Name,
			Id:         p.Id,
			P4Program:  p.P4Program,
			Status:     p.Status,
			CreateTime: timestamppb.New(p.CreateTime),
		})
	}

	writeJSON(w, http.StatusOK, PipelineListResponse{Pipelines: pipelines})
}

func handleGetPipeline(w http.ResponseWriter, r *http.Request, svc *PipelineService, pipelineId string) {
	svc.mu.RLock()
	defer svc.mu.RUnlock()

	p, ok := svc.pipelines[pipelineId]
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": fmt.Sprintf("pipeline %q not found", pipelineId)})
		return
	}

	writeJSON(w, http.StatusOK, PipelineInfo{
		Name:       p.Name,
		Id:         p.Id,
		P4Program:  p.P4Program,
		Status:     p.Status,
		CreateTime: timestamppb.New(p.CreateTime),
	})
}

func handleCreatePipeline(w http.ResponseWriter, r *http.Request, svc *PipelineService) {
	var req CreatePipelineRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		req = CreatePipelineRequest{}
	}

	svc.mu.Lock()
	defer svc.mu.Unlock()

	id := fmt.Sprintf("pipeline-%d", len(svc.pipelines))
	name := fmt.Sprintf("pipelines/%s", id)

	if req.P4Program == "" {
		req.P4Program = "basic_forwarding.p4"
	}

	pipeline := &Pipeline{
		Name:       name,
		Id:         id,
		P4Program:  req.P4Program,
		Status:     "LOADED",
		CreateTime: time.Now(),
	}

	svc.pipelines[id] = pipeline
	log.Printf("[REST] Created pipeline: %s (P4: %s)", id, pipeline.P4Program)

	writeJSON(w, http.StatusCreated, PipelineInfo{
		Name:       pipeline.Name,
		Id:         pipeline.Id,
		P4Program:  pipeline.P4Program,
		Status:     pipeline.Status,
		CreateTime: timestamppb.New(pipeline.CreateTime),
	})
}

func handleDeletePipeline(w http.ResponseWriter, r *http.Request, svc *PipelineService, pipelineId string) {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	if _, ok := svc.pipelines[pipelineId]; !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": fmt.Sprintf("pipeline %q not found", pipelineId)})
		return
	}

	delete(svc.pipelines, pipelineId)
	log.Printf("[REST] Deleted pipeline: %s", pipelineId)

	w.WriteHeader(http.StatusNoContent)
}

// --- Helpers ---

func writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-OPI-Sandbox", "true")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("[REST] Error encoding response: %v", err)
	}
}
