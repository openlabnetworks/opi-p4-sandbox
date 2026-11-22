# Security Analysis — Internal Team Reference

> **CONFIDENTIAL — Do not publish or share outside the team.**
> This document accompanies the public security page and contains the honest trade-off analysis, known limitations, and accepted risks that we deliberately keep off the public-facing site.

**Version 1.0, November 2026**

---

## Pros and cons of our self-hosted approach

We chose to run our own protection rather than put the sandbox behind a third-party security service. Here is the honest trade.

### Strengths

- **No outside dependency.** No company sees our traffic, and no outage or policy change elsewhere can switch our protection off.
- **Full control and transparency.** Every rule lives in our repository, so it can be reviewed, tested and explained to partners.
- **Predictable cost.** Open-source tools with no per-request pricing.
- **Consistent with what we build.** We work on open, programmable infrastructure, and we secure it with open components.
- **Strong against common attacks.** Application-level floods, brute force, scanning and resource exhaustion are all handled.

### Trade-offs

- **Large floods can still win.** If traffic exceeds our network link, it saturates before our firewall sees it. We reduce this with a hosting provider that filters at the network level and with fast rebuilds.
- **More work for us.** Patching, tuning, monitoring and on-call are our responsibility.
- **Real users can be blocked.** Many people share one public IP (colleges, offices, mobile carriers), so per-IP limits need careful tuning.
- **Authentication adds friction.** Visitors need a key before using the API, which slows the quick-try experience.
- **No global edge network.** Visitors far from our server see higher latency and there is no edge caching.
- **Mistakes are ours alone.** A wrong firewall rule can lock us out or leave a gap, so changes need testing and review.

## Risks we accept

- A very large distributed flood could take the sandbox offline for a while. It is a free demo, so we plan to recover quickly rather than promise it never goes down.
- A newly discovered flaw in nginx, the P4 compiler, the BMv2 switch or the operating system could bypass a layer before a patch exists.
- A stolen API key works until we notice and revoke it, so quotas and monitoring matter.

## Defense layer gap analysis

| Layer | Does not stop |
|---|---|
| 1. Host firewall (nftables) | Floods larger than our network link |
| 2. Reverse proxy (nginx) | Low-rate attacks that look like normal use |
| 3. Authentication and fail2ban | A leaked valid key, until it is revoked |
| 4. Container sandbox | Unknown bugs in the toolchain itself |
| Monitoring | Prevention (detection only) |

---

*OpenLabNetworks — Internal use only*
