# Securing the OPI & P4 Developer Sandbox

A security report on how we protect a public sandbox without relying on a third-party security service.

**Version 1.0, November 2026.** Status: design and rollout plan. Layer status below is updated as each one goes live.

---

## Summary

The sandbox lets developers try OPI APIs and P4 programs without hardware. Because anyone can reach it, it will attract scanning, brute-force attempts and floods. We defend it with four self-hosted layers (firewall, reverse proxy, authentication, container limits) plus monitoring, all built from open-source tools that we run ourselves.

## What we defend against

| Threat | Example | Our answer |
|---|---|---|
| Traffic floods | Thousands of HTTP or SYN requests per second | Firewall limits, nginx rate limits, container caps |
| Key guessing | Automated attempts at API keys | Long random keys, lockouts, fail2ban |
| Malicious P4 programs | Oversized or crafted uploads aimed at the compiler | Size limits, input checks, jailed compiles |
| Exposed internals | Switch control ports reachable from the internet | Internal-only network, only ports 80 and 443 public |
| Supply chain | Tampered images or dependencies | Pinned versions, signed images, scanning |

## The defense layers

| Layer | What it does | Status |
|---|---|---|
| 1. Host firewall (nftables) | Blocks everything by default, limits connections per IP, SYN cookies | Planned |
| 2. Reverse proxy (nginx) | TLS, request rate limits, body size caps, timeouts | Planned |
| 3. Authentication and fail2ban | API keys required, per-key quotas, automatic bans | Planned |
| 4. Container sandbox | CPU, memory and time caps, P4 compiles in a jail with no network | Planned |
| Monitoring | Metrics, alerts and uptime checks | Planned |

```
Internet
   │
[ 1. Host firewall: nftables ]     drop junk, SYN cookies, per-IP limits
   │
[ 2. nginx reverse proxy ]         TLS, rate limits, body limits, timeouts
   │
[ 3. Auth + fail2ban ]             API keys, quotas, automatic bans
   │
[ 4. Container sandbox ]           CPU/memory/time caps, compile jail
   │
[ OPI server ] ── [ P4 BMv2 switch ]    internal network only
```

## Rollout plan

1. Close exposure: firewall, SSH hardening, TLS, control ports internal only.
2. Add authentication, quotas and automatic bans.
3. Harden containers and jail P4 compilation.
4. Add monitoring and alerts, then load-test our own limits.
5. Secure the supply chain: signed images, pinned dependencies and repository protections.

## Beyond the sandbox

When SmartNIC hardware joins the platform, the same layered thinking extends into the card: secure boot, signed firmware and pipelines, locked debug ports, inline encryption and tenant isolation. That work is a later phase and nothing in this report claims it exists today.

---

*OpenLabNetworks: Accelerating cloud-native infrastructure*
