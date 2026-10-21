// Package handlers provides HTTP/JSON REST handlers that emulate OPI APIs.
// This file implements the REST API server that runs alongside or instead of
// raw gRPC, making the sandbox accessible via curl without needing grpcurl.
package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

const maxRequestBodyBytes = 1 << 20

// Validation patterns for user-supplied input.
var (
	// macRegex matches standard 6-octet colon-separated MAC addresses.
	macRegex = regexp.MustCompile(`^([0-9a-fA-F]{2}:){5}[0-9a-fA-F]{2}$`)
	// p4ProgramRegex allows alphanumeric, underscores, hyphens, and dots; must end in .p4
	p4ProgramRegex = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.\-]{0,62}\.p4$`)
)

// decodeJSONBody validates and decodes a request body. An empty body is valid
// and leaves dst at its zero value so callers can apply their defaults.
func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) bool {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaType != "application/json" {
		writeJSON(w, http.StatusUnsupportedMediaType, map[string]string{
			"error": "Content-Type must be application/json",
		})
		return false
	}

	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxRequestBodyBytes))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dst); err != nil && !errors.Is(err, io.EOF) {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid JSON body: " + err.Error(),
		})
		return false
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid JSON body: must contain a single JSON object",
		})
		return false
	}
	return true
}

// RegisterRESTRoutes registers all OPI-style REST API routes on the given mux.
func RegisterRESTRoutes(mux *http.ServeMux, netSvc *NetworkService, pipeSvc *PipelineService) {
	// Web Testing Dashboard
	mux.HandleFunc("/ui", ServeUI)
	mux.HandleFunc("/ui/", ServeUI)

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
	// Content negotiation: serve Web UI if requested by browser
	if strings.Contains(r.Header.Get("Accept"), "text/html") {
		ServeUI(w, r)
		return
	}
	writeJSON(w, http.StatusOK, APIInfoResponse{
		Name:    "Intili Sandbox API Server",
		Version: "0.1.0",
		Services: []string{
			"opi.v1.NetworkService",
			"opi.v1.PipelineService",
		},
		DocsURL: "https://github.com/openlabnetworks/opi-p4-sandbox",
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
	if !decodeJSONBody(w, r, &req) {
		return
	}

	// --- Input validation ---
	parsedMAC, err := net.ParseMAC(req.MacAddress)
	if req.MacAddress != "" && (err != nil || len(parsedMAC) != 6 || !macRegex.MatchString(req.MacAddress)) {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid mac_address format; expected XX:XX:XX:XX:XX:XX",
		})
		return
	}
	if req.Mtu != 0 && (req.Mtu < 68 || req.Mtu > 9216) {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "mtu must be between 68 and 9216",
		})
		return
	}

	svc.mu.Lock()
	defer svc.mu.Unlock()

	// Use monotonic counter — never reuses IDs even after deletes
	seq := svc.nextID.Add(1) - 1
	id := fmt.Sprintf("port-%d", seq)
	name := fmt.Sprintf("ports/%s", id)

	// Defensive: reject if ID already exists (should never happen with a counter)
	if _, exists := svc.ports[id]; exists {
		writeJSON(w, http.StatusConflict, map[string]string{
			"error": fmt.Sprintf("port %q already exists", id),
		})
		return
	}

	if req.MacAddress == "" {
		req.MacAddress = fmt.Sprintf("00:1a:2b:%02x:%02x:%02x", (seq>>16)&0xff, (seq>>8)&0xff, seq&0xff)
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
	if !decodeJSONBody(w, r, &req) {
		return
	}

	// --- Input validation ---
	if req.P4Program != "" && !p4ProgramRegex.MatchString(req.P4Program) {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid p4_program name; must be alphanumeric with underscores/hyphens, ending in .p4",
		})
		return
	}

	svc.mu.Lock()
	defer svc.mu.Unlock()

	// Use monotonic counter — never reuses IDs even after deletes
	seq := svc.nextID.Add(1) - 1
	id := fmt.Sprintf("pipeline-%d", seq)
	name := fmt.Sprintf("pipelines/%s", id)

	// Defensive: reject if ID already exists (should never happen with a counter)
	if _, exists := svc.pipelines[id]; exists {
		writeJSON(w, http.StatusConflict, map[string]string{
			"error": fmt.Sprintf("pipeline %q already exists", id),
		})
		return
	}

	if req.P4Program == "" {
		req.P4Program = "basic_forwarding.p4"
	}

	pipeline := &Pipeline{
		Name:       name,
		Id:         id,
		P4Program:  req.P4Program,
		Status:     "REGISTERED",
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
	w.Header().Set("X-Intili-Sandbox", "true")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("[REST] Error encoding response: %v", err)
	}
}
