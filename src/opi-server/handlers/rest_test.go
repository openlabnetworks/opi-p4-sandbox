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
	if createdPipe.P4Program != "firewall_drop.p4" || createdPipe.Status != "REGISTERED" {
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

func TestServeUI(t *testing.T) {
	mux := setupTestServer()

	for _, path := range []string{"/ui", "/ui/"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200 on %s, got %d", path, rec.Code)
		}
		ct := rec.Header().Get("Content-Type")
		if ct != "text/html; charset=utf-8" {
			t.Errorf("expected text/html content-type on %s, got %s", path, ct)
		}
		body := rec.Body.String()
		if !bytes.Contains([]byte(body), []byte("OPI & P4 Developer Sandbox")) {
			t.Errorf("expected body to contain sandbox title on %s", path)
		}
	}
}

func TestContentNegotiation(t *testing.T) {
	mux := setupTestServer()

	// Browser request with text/html in Accept header
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 on browser GET /, got %d", rec.Code)
	}
	if rec.Header().Get("Content-Type") != "text/html; charset=utf-8" {
		t.Errorf("expected HTML for browser request, got %s", rec.Header().Get("Content-Type"))
	}

	// API/curl request with application/json or */*
	reqAPI := httptest.NewRequest(http.MethodGet, "/", nil)
	recAPI := httptest.NewRecorder()
	mux.ServeHTTP(recAPI, reqAPI)

	if recAPI.Code != http.StatusOK {
		t.Fatalf("expected 200 on API GET /, got %d", recAPI.Code)
	}
	if recAPI.Header().Get("Content-Type") != "application/json" {
		t.Errorf("expected application/json for API request, got %s", recAPI.Header().Get("Content-Type"))
	}
}

func TestCreateValidation(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		contentType string
		body        string
		wantStatus  int
	}{
		{"empty port body uses defaults", "/v1/ports", "application/json", "", http.StatusCreated},
		{"empty pipeline body uses defaults", "/v1/pipelines", "application/json", "", http.StatusCreated},
		{"port requires JSON media type", "/v1/ports", "text/plain", `{"mac_address":"aa:bb:cc:dd:ee:ff"}`, http.StatusUnsupportedMediaType},
		{"pipeline requires JSON media type", "/v1/pipelines", "text/plain", `{"p4_program":"safe.p4"}`, http.StatusUnsupportedMediaType},
		{"malformed port JSON", "/v1/ports", "application/json", "{garbage", http.StatusBadRequest},
		{"malformed pipeline JSON", "/v1/pipelines", "application/json", "{garbage", http.StatusBadRequest},
		{"invalid MAC", "/v1/ports", "application/json", `{"mac_address":"not-a-mac"}`, http.StatusBadRequest},
		{"invalid MTU", "/v1/ports", "application/json", `{"mtu":-5}`, http.StatusBadRequest},
		{"oversized MTU", "/v1/ports", "application/json", `{"mtu":9217}`, http.StatusBadRequest},
		{"unsafe P4 program", "/v1/pipelines", "application/json", `{"p4_program":"<img src=x onerror=alert(1)>"}`, http.StatusBadRequest},
		{"unknown port field", "/v1/ports", "application/json", `{"unexpected":true}`, http.StatusBadRequest},
		{"multiple JSON values", "/v1/pipelines", "application/json", `{} {}`, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := setupTestServer()
			req := httptest.NewRequest(http.MethodPost, tt.path, bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", tt.contentType)
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)
			if rec.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d: %s", tt.wantStatus, rec.Code, rec.Body.String())
			}
		})
	}
}

// TestCreateDeleteCreatePort is a regression test for GitHub issue #1:
// creating a port after a delete must never reuse an ID or silently overwrite
// an existing port. The monotonic counter must always advance.
func TestCreateDeleteCreatePort(t *testing.T) {
	mux := setupTestServer()

	// Step 1: Create port (should be port-1 because port-0 is pre-populated)
	req := httptest.NewRequest(http.MethodPost, "/v1/ports", bytes.NewBufferString(`{"mac_address":"aa:bb:cc:dd:ee:01","mtu":1500}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("step 1: expected 201, got %d", rec.Code)
	}
	var port1 PortInfo
	json.NewDecoder(rec.Body).Decode(&port1)
	if port1.Id != "port-1" {
		t.Fatalf("step 1: expected id port-1, got %s", port1.Id)
	}

	// Step 2: Delete port-1
	req = httptest.NewRequest(http.MethodDelete, "/v1/ports/port-1", nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("step 2: expected 204, got %d", rec.Code)
	}

	// Step 3: Create another port — MUST get port-2, NOT port-1 again
	req = httptest.NewRequest(http.MethodPost, "/v1/ports", bytes.NewBufferString(`{"mac_address":"aa:bb:cc:dd:ee:02","mtu":9000}`))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("step 3: expected 201, got %d", rec.Code)
	}
	var port2 PortInfo
	json.NewDecoder(rec.Body).Decode(&port2)
	if port2.Id != "port-2" {
		t.Fatalf("step 3: expected id port-2, got %s (ID was reused!)", port2.Id)
	}

	// Step 4: Verify port-0 (pre-populated) was NOT overwritten
	req = httptest.NewRequest(http.MethodGet, "/v1/ports/port-0", nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("step 4: port-0 should still exist, got %d", rec.Code)
	}
	var port0 PortInfo
	json.NewDecoder(rec.Body).Decode(&port0)
	if port0.MacAddress != "00:1a:2b:3c:4d:5e" {
		t.Fatalf("step 4: port-0 was silently overwritten! MAC = %s", port0.MacAddress)
	}

	// Step 5: Verify all ports have unique IDs (list them)
	req = httptest.NewRequest(http.MethodGet, "/v1/ports", nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	var listResp PortListResponse
	json.NewDecoder(rec.Body).Decode(&listResp)
	if len(listResp.Ports) != 2 {
		t.Fatalf("step 5: expected 2 ports (port-0 + port-2), got %d", len(listResp.Ports))
	}
	ids := map[string]bool{}
	for _, p := range listResp.Ports {
		if ids[p.Id] {
			t.Fatalf("step 5: duplicate port ID %s detected!", p.Id)
		}
		ids[p.Id] = true
	}
}

// TestCreateDeleteCreatePipeline is a regression test for GitHub issue #1:
// creating a pipeline after a delete must never reuse an ID or silently
// overwrite an existing pipeline.
func TestCreateDeleteCreatePipeline(t *testing.T) {
	mux := setupTestServer()

	// Step 1: Create pipeline (should be pipeline-0, default is "pipeline-default")
	req := httptest.NewRequest(http.MethodPost, "/v1/pipelines", bytes.NewBufferString(`{"p4_program":"firewall.p4"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("step 1: expected 201, got %d", rec.Code)
	}
	var pipe1 PipelineInfo
	json.NewDecoder(rec.Body).Decode(&pipe1)
	if pipe1.Id != "pipeline-0" {
		t.Fatalf("step 1: expected id pipeline-0, got %s", pipe1.Id)
	}

	// Step 2: Delete pipeline-0
	req = httptest.NewRequest(http.MethodDelete, "/v1/pipelines/pipeline-0", nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("step 2: expected 204, got %d", rec.Code)
	}

	// Step 3: Create another pipeline — MUST get pipeline-1, NOT pipeline-0 again
	req = httptest.NewRequest(http.MethodPost, "/v1/pipelines", bytes.NewBufferString(`{"p4_program":"acl.p4"}`))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("step 3: expected 201, got %d", rec.Code)
	}
	var pipe2 PipelineInfo
	json.NewDecoder(rec.Body).Decode(&pipe2)
	if pipe2.Id != "pipeline-1" {
		t.Fatalf("step 3: expected id pipeline-1, got %s (ID was reused!)", pipe2.Id)
	}

	// Step 4: Verify pipeline-default was NOT overwritten
	req = httptest.NewRequest(http.MethodGet, "/v1/pipelines/pipeline-default", nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("step 4: pipeline-default should still exist, got %d", rec.Code)
	}
	var pipeDefault PipelineInfo
	json.NewDecoder(rec.Body).Decode(&pipeDefault)
	if pipeDefault.P4Program != "basic_forwarding.p4" {
		t.Fatalf("step 4: pipeline-default was silently overwritten! P4 = %s", pipeDefault.P4Program)
	}

	// Step 5: Verify all pipelines have unique IDs
	req = httptest.NewRequest(http.MethodGet, "/v1/pipelines", nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	var listResp PipelineListResponse
	json.NewDecoder(rec.Body).Decode(&listResp)
	if len(listResp.Pipelines) != 2 {
		t.Fatalf("step 5: expected 2 pipelines (pipeline-default + pipeline-1), got %d", len(listResp.Pipelines))
	}
	ids := map[string]bool{}
	for _, p := range listResp.Pipelines {
		if ids[p.Id] {
			t.Fatalf("step 5: duplicate pipeline ID %s detected!", p.Id)
		}
		ids[p.Id] = true
	}
}
