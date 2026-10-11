package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTestServer() *http.ServeMux {
	mux := http.NewServeMux()
	netSvc := NewNetworkService()
	pipeSvc := NewPipelineService()
	RegisterRESTRoutes(mux, netSvc, pipeSvc)
	return mux
}

func TestAPIInfo(t *testing.T) {
	mux := setupTestServer()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var resp APIInfoResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Name != "OPI P4 Sandbox API Server" {
		t.Errorf("unexpected name: %s", resp.Name)
	}
	if len(resp.Services) != 2 {
		t.Errorf("expected 2 services, got %d", len(resp.Services))
	}
}

func TestNotFound(t *testing.T) {
	mux := setupTestServer()

	req := httptest.NewRequest(http.MethodGet, "/unknown-path", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}
}

func TestHealthCheck(t *testing.T) {
	mux := setupTestServer()

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var resp HealthResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Status != "SERVING" {
		t.Errorf("expected status SERVING, got %s", resp.Status)
	}
}

func TestPortLifecycle(t *testing.T) {
	mux := setupTestServer()

	// 1. List default ports
	req := httptest.NewRequest(http.MethodGet, "/v1/ports", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var listResp PortListResponse
	if err := json.NewDecoder(rec.Body).Decode(&listResp); err != nil {
		t.Fatalf("failed to decode list response: %v", err)
	}
	if len(listResp.Ports) != 1 || listResp.Ports[0].Id != "port-0" {
		t.Errorf("expected default port-0, got %v", listResp.Ports)
	}

	// 2. Get specific port
	req = httptest.NewRequest(http.MethodGet, "/v1/ports/port-0", nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200 for port-0, got %d", rec.Code)
	}

	// 3. Create a new port
	createBody := []byte(`{"mac_address":"00:11:22:33:44:55","mtu":9000}`)
	req = httptest.NewRequest(http.MethodPost, "/v1/ports", bytes.NewBuffer(createBody))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201 on create port, got %d", rec.Code)
	}

	var createdPort PortInfo
	if err := json.NewDecoder(rec.Body).Decode(&createdPort); err != nil {
		t.Fatalf("failed to decode created port: %v", err)
	}
	if createdPort.MacAddress != "00:11:22:33:44:55" || createdPort.Mtu != 9000 {
		t.Errorf("unexpected port fields: %+v", createdPort)
	}

	// 4. Delete port
	req = httptest.NewRequest(http.MethodDelete, "/v1/ports/"+createdPort.Id, nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status 204 on delete port, got %d", rec.Code)
	}

	// 5. Verify deleted port returns 404
	req = httptest.NewRequest(http.MethodGet, "/v1/ports/"+createdPort.Id, nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404 for deleted port, got %d", rec.Code)
	}

	// 6. Delete nonexistent port
	req = httptest.NewRequest(http.MethodDelete, "/v1/ports/nonexistent", nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404 deleting nonexistent port, got %d", rec.Code)
	}
}

func TestPipelineLifecycle(t *testing.T) {
	mux := setupTestServer()

	// 1. List default pipelines
	req := httptest.NewRequest(http.MethodGet, "/v1/pipelines", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var listResp PipelineListResponse
	if err := json.NewDecoder(rec.Body).Decode(&listResp); err != nil {
		t.Fatalf("failed to decode list response: %v", err)
	}
	if len(listResp.Pipelines) != 1 || listResp.Pipelines[0].Id != "pipeline-default" {
		t.Errorf("expected default pipeline-default, got %v", listResp.Pipelines)
	}

	// 2. Get specific pipeline
	req = httptest.NewRequest(http.MethodGet, "/v1/pipelines/pipeline-default", nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	// 3. Create a new pipeline
	createBody := []byte(`{"p4_program":"firewall_drop.p4"}`)
	req = httptest.NewRequest(http.MethodPost, "/v1/pipelines", bytes.NewBuffer(createBody))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201 on create pipeline, got %d", rec.Code)
	}

	var createdPipe PipelineInfo
	if err := json.NewDecoder(rec.Body).Decode(&createdPipe); err != nil {
		t.Fatalf("failed to decode created pipeline: %v", err)
	}
	if createdPipe.P4Program != "firewall_drop.p4" || createdPipe.Status != "LOADED" {
		t.Errorf("unexpected pipeline fields: %+v", createdPipe)
	}

	// 4. Delete pipeline
	req = httptest.NewRequest(http.MethodDelete, "/v1/pipelines/"+createdPipe.Id, nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status 204 on delete pipeline, got %d", rec.Code)
	}

	// 5. Verify deleted pipeline returns 404
	req = httptest.NewRequest(http.MethodGet, "/v1/pipelines/"+createdPipe.Id, nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404 for deleted pipeline, got %d", rec.Code)
	}
}

func TestMethodNotAllowedAndValidation(t *testing.T) {
	mux := setupTestServer()

	// Method not allowed on /v1/ports
	req := httptest.NewRequest(http.MethodPatch, "/v1/ports", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405 on PATCH /v1/ports, got %d", rec.Code)
	}

	// Empty port ID on /v1/ports/
	req = httptest.NewRequest(http.MethodGet, "/v1/ports/", nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 on empty port id, got %d", rec.Code)
	}

	// Method not allowed on /v1/pipelines
	req = httptest.NewRequest(http.MethodPut, "/v1/pipelines", nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405 on PUT /v1/pipelines, got %d", rec.Code)
	}

	// Empty pipeline ID on /v1/pipelines/
	req = httptest.NewRequest(http.MethodGet, "/v1/pipelines/", nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 on empty pipeline id, got %d", rec.Code)
	}
}
