#!/bin/bash
# ============================================================
# Virtual Ethernet Pair Setup Script
# ============================================================
# Creates virtual ethernet pairs for P4 BMv2 packet I/O.
# Run with: sudo ./setup-veth.sh
# ============================================================
set -e

NUM_PAIRS="${1:-4}"

echo "Creating ${NUM_PAIRS} virtual ethernet pairs..."

for i in $(seq 0 $((NUM_PAIRS - 1))); do
    VETH="veth${i}"
    PEER="veth${i}_peer"

    if ip link show "$VETH" &>/dev/null; then
        echo "  ${VETH} already exists, skipping"
        continue
    fi

    echo "  Creating: ${VETH} <-> ${PEER}"
    ip link add "$VETH" type veth peer name "$PEER"
    ip link set "$VETH" up
    ip link set "$PEER" up

    # Disable IPv6 on veth interfaces (avoids unnecessary traffic)
    sysctl -q net.ipv6.conf.${VETH}.disable_ipv6=1 2>/dev/null || true
    sysctl -q net.ipv6.conf.${PEER}.disable_ipv6=1 2>/dev/null || true

    # Set MTU to jumbo frames
    ip link set "$VETH" mtu 9500
    ip link set "$PEER" mtu 9500
done

echo "Done! Virtual ethernet pairs are ready."
echo ""
echo "Verify with: ip link show type veth"
