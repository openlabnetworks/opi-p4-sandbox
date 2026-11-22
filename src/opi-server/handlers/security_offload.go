package handlers

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

const securityOffloadSimulationDisclaimer = "This is a control-plane and datapath model only. The sandbox OPI server does not currently program BMv2 through P4Runtime, and no physical SmartNIC, DPU, cryptographic engine, or Kubernetes cluster is configured by this request."

var namespaceNameRegex = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]{0,61}[a-z0-9])?$`)

// handleSecurityOffloadSimulation models compilation of a Kubernetes policy
// into a P4 ACL entry. Keeping the model behind a REST endpoint lets external
// clients use the same deterministic simulation as the security page.
func handleSecurityOffloadSimulation(w http.ResponseWriter, r *http.Request) {
	var req SecurityOffloadSimulationRequest
	if !decodeJSONBody(w, r, &req) {
		return
	}

	if !namespaceNameRegex.MatchString(req.SourceNamespace) || !namespaceNameRegex.MatchString(req.DestinationNamespace) {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "source_namespace and destination_namespace must be valid Kubernetes namespace names"})
		return
	}
	if req.Protocol != "TCP" && req.Protocol != "UDP" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "protocol must be TCP or UDP"})
		return
	}
	if req.DestinationPort < 1 || req.DestinationPort > 65535 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "destination_port must be between 1 and 65535"})
		return
	}
	if req.PolicyAction != "allow" && req.PolicyAction != "deny" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "policy_action must be allow or deny"})
		return
	}

	if req.AttackType == "" {
		req.AttackType = "normal"
	}
	validAttackTypes := map[string]string{
		"normal":            "Authorized workload traffic",
		"syn_flood":         "Volumetric SYN flood",
		"port_scan":         "Reconnaissance / port scan",
		"dns_amplification": "DDoS reflection / amplification",
		"lateral_movement":  "Unauthorized east-west movement",
		"malformed_packet":  "Malformed / protocol violation",
		"ddos_http":         "Application-layer DDoS",
	}
	threatClass, validAttack := validAttackTypes[req.AttackType]
	if !validAttack {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "attack_type is not supported"})
		return
	}

	decision := strings.ToUpper(req.PolicyAction)
	action := "permit"
	details := "The modeled SmartNIC P4 ACL forwards the matching packet to the host workload."
	if req.PolicyAction == "deny" || req.AttackType != "normal" {
		action = "drop"
		details = fmt.Sprintf("The modeled SmartNIC P4 ACL drops %s before it reaches the host network stack.", threatClass)
		decision = "DENY"
	}

	policyMatch := fmt.Sprintf("src_ns=%s,dst_ns=%s,proto=%s,dst_port=%d,threat=%s", req.SourceNamespace, req.DestinationNamespace, req.Protocol, req.DestinationPort, req.AttackType)
	writeJSON(w, http.StatusOK, SecurityOffloadSimulationResponse{
		Mode:             "SIMULATION",
		Decision:         decision,
		EnforcementPoint: "SmartNIC/DPU P4 ACL (modeled)",
		P4TableEntry:     fmt.Sprintf("k8s_policy_acl[%s] -> %s", policyMatch, action),
		ControlPlane: []string{
			"Kubernetes NetworkPolicy selects the workload traffic.",
			"A custom security operator resolves workload identity and compiles an ACL intent.",
			"An OPI/gRPC control-plane service would write the matching P4 table entry to the DPU.",
		},
		Datapath: []string{
			"The NIC classifies the packet against the modeled P4 ACL.",
			fmt.Sprintf("Matching action: %s.", action),
		},
		HostCPUImpact: "Modeled as zero host CPU cycles for the packet-enforcement decision.",
		ThreatClass:   threatClass,
		Details:       details,
		Disclaimer:    securityOffloadSimulationDisclaimer,
	})
}
