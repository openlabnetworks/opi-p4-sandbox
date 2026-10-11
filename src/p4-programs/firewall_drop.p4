/* -*- P4_16 -*- */
/* 
 * Firewall Drop — Custom P4 Program
 * Drops all non-IPv4 traffic.
 */

#include <core.p4>
#include <v1model.p4>

/* --- Headers --- */
header ethernet_t {
    bit<48> dstAddr;
    bit<48> srcAddr;
    bit<16> etherType;
}

struct headers {
    ethernet_t ethernet;
}

struct metadata {
    // empty
}

/* --- Parser --- */
parser CustomParser(
    packet_in packet,
    out headers hdr,
    inout metadata meta,
    inout standard_metadata_t standard_metadata
) {
    state start {
        // Always read the first 14 bytes as an Ethernet header
        packet.extract(hdr.ethernet);
        transition accept;
    }
}

/* --- Ingress (The Brain) --- */
control CustomIngress(
    inout headers hdr,
    inout metadata meta,
    inout standard_metadata_t standard_metadata
) {
    apply {
        // If the Ethernet Type is NOT 0x0800 (IPv4)
        if (hdr.ethernet.etherType != 0x0800) {
            // Instantly drop the packet!
            mark_to_drop(standard_metadata);
        } else {
            // Otherwise, let it pass to Port 1
            standard_metadata.egress_spec = 1;
        }
    }
}

/* --- Boilerplate blocks --- */
control CustomVerifyChecksum(inout headers hdr, inout metadata meta) { apply { } }
control CustomEgress(inout headers hdr, inout metadata meta, inout standard_metadata_t smeta) { apply { } }
control CustomComputeChecksum(inout headers hdr, inout metadata meta) { apply { } }

/* --- Deparser --- */
control CustomDeparser(packet_out packet, in headers hdr) {
    apply {
        packet.emit(hdr.ethernet);
    }
}

/* --- Instantiate the Switch --- */
V1Switch(
    CustomParser(),
    CustomVerifyChecksum(),
    CustomIngress(),
    CustomEgress(),
    CustomComputeChecksum(),
    CustomDeparser()
) main;
