// Package handlers - gRPC type definitions and service descriptors.
// These types mirror the OPI API protobuf messages in a lightweight way
// so the sandbox compiles without requiring the full OPI protobuf toolchain.
//
// For the sandbox, we use a simplified approach: a standard net/http JSON API
// server that emulates OPI-style resource management, combined with gRPC
// server reflection so tools like grpcurl can discover the service.
package handlers

import (
	"google.golang.org/protobuf/types/known/timestamppb"
)

// --- Network Service Types ---

// PortInfo represents a network port resource following OPI naming conventions.
type PortInfo struct {
	Name       string                 `json:"name"`
	Id         string                 `json:"id"`
	MacAddress string                 `json:"mac_address"`
	Mtu        int32                  `json:"mtu"`
	AdminState string                 `json:"admin_state"`
	OperState  string                 `json:"oper_state"`
	Speed      string                 `json:"speed"`
	CreateTime *timestamppb.Timestamp `json:"create_time"`
}

// PortListResponse is the response for ListPorts.
type PortListResponse struct {
	Ports []*PortInfo `json:"ports"`
}

// GetPortRequest requests a specific port by ID.
type GetPortRequest struct {
	PortId string `json:"port_id"`
}

// CreatePortRequest defines fields for creating a new port.
type CreatePortRequest struct {
	MacAddress string `json:"mac_address"`
	Mtu        int32  `json:"mtu"`
}

// --- Pipeline Service Types ---

// PipelineInfo represents a P4 pipeline resource.
type PipelineInfo struct {
	Name       string                 `json:"name"`
	Id         string                 `json:"id"`
	P4Program  string                 `json:"p4_program"`
	Status     string                 `json:"status"`
	CreateTime *timestamppb.Timestamp `json:"create_time"`
}

// PipelineListResponse is the response for ListPipelines.
type PipelineListResponse struct {
	Pipelines []*PipelineInfo `json:"pipelines"`
}

// GetPipelineRequest requests a specific pipeline by ID.
type GetPipelineRequest struct {
	PipelineId string `json:"pipeline_id"`
}

// CreatePipelineRequest defines fields for creating a new pipeline.
type CreatePipelineRequest struct {
	P4Program string `json:"p4_program"`
}

// --- Health Types ---

// HealthResponse represents the server health status.
type HealthResponse struct {
	Status   string            `json:"status"`
	Services map[string]string `json:"services"`
}

// --- API Info ---

// APIInfoResponse provides metadata about the OPI sandbox server.
type APIInfoResponse struct {
	Name     string   `json:"name"`
	Version  string   `json:"version"`
	Services []string `json:"services"`
	DocsURL  string   `json:"docs_url"`
}
