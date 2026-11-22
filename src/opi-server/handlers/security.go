package handlers

import (
	"net/http"
)

// ServeSecurityPage renders the public security overview page.
// Sensitive information (trade-offs, known gaps, accepted risks) is kept in
// docs/security-internal-team.md and is NOT shown on this public page.
func ServeSecurityPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(securityHTML))
}

const securityHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Security | ProtonNIC Sandbox</title>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=DM+Mono:wght@400;500&family=Manrope:wght@500;700;800&display=swap" rel="stylesheet">
  <style>
    :root {
      --ink: #102022;
      --muted: #587074;
      --line: #c9d7d3;
      --paper: #f4f8f5;
      --card: #ffffff;
      --accent: #d96b3b;
      --teal: #0e7771;
      --red: #c0392b;
      --green-bg: #e8f5e9;
      --green-txt: #1b5e20;
      --amber-bg: #fff8e1;
      --amber-txt: #e65100;
      --mono: "DM Mono", monospace;
      --sans: "Manrope", sans-serif;
    }
    * { box-sizing: border-box; }
    body { margin: 0; color: var(--ink); background: var(--paper); font-family: var(--sans); }
    a { color: inherit; }

    /* Nav */
    .sec-nav {
      display: flex; justify-content: space-between; align-items: center;
      padding: 24px 7vw; border-bottom: 1px solid var(--line);
    }
    .logo { text-decoration: none; font-weight: 800; letter-spacing: -.04em; }
    .logo span { color: var(--accent); }
    .nav-links { display: flex; gap: 28px; align-items: center; }
    .nav-links a {
      font: 12px var(--mono); text-transform: uppercase; text-decoration: none;
      color: var(--muted); letter-spacing: .08em; transition: color .15s;
    }
    .nav-links a:hover { color: var(--ink); }
    .nav-links a.active { color: var(--ink); font-weight: 500; }

    /* Hero */
    .sec-hero { max-width: 1180px; margin: 0 auto; padding: 84px 7vw 0; }
    .eyebrow {
      color: var(--teal); font: 12px var(--mono);
      letter-spacing: .12em; text-transform: uppercase;
    }
    .sec-hero h1 {
      max-width: 820px; margin: 18px 0 24px;
      font-size: clamp(40px, 7vw, 88px); line-height: .96; letter-spacing: -.06em;
    }
    .intro {
      max-width: 700px; color: var(--muted); font-size: 19px; line-height: 1.65;
    }
    .version-badge {
      display: inline-block; margin-top: 20px; padding: 6px 14px;
      font: 11px var(--mono); letter-spacing: .06em; text-transform: uppercase;
      border: 1px solid var(--line); color: var(--muted); background: var(--card);
    }

    /* Sections */
    .sec-main { max-width: 1180px; margin: 0 auto; padding: 0 7vw 110px; }
    .section { margin-top: 90px; }
    h2 { margin: 0 0 14px; font-size: 32px; letter-spacing: -.05em; }
    .section-lead { max-width: 650px; color: var(--muted); line-height: 1.7; margin-bottom: 28px; }

    /* Threat table */
    .threat-table { width: 100%; border-collapse: collapse; margin-top: 24px; }
    .threat-table th {
      text-align: left; font: 11px var(--mono); text-transform: uppercase;
      letter-spacing: .08em; color: var(--muted); padding: 12px 16px;
      border-bottom: 2px solid var(--ink); font-weight: 500;
    }
    .threat-table td {
      padding: 16px; border-bottom: 1px solid var(--line);
      font-size: 15px; line-height: 1.5;
    }
    .threat-table tr:hover td { background: rgba(14,119,113,.04); }
    .threat-name { font-weight: 700; }
    .threat-example { color: var(--muted); font-size: 14px; }
    .threat-answer { color: var(--teal); font-weight: 600; }

    /* Defense layer stack */
    .layer-stack {
      border-left: 3px solid var(--accent); margin-top: 32px;
    }
    .layer-item {
      padding: 22px 26px; border-bottom: 1px solid var(--line);
      background: var(--card); display: flex; justify-content: space-between;
      align-items: flex-start; gap: 20px;
      transition: background .15s;
    }
    .layer-item:hover { background: #f0f5f2; }
    .layer-item:last-child { border-bottom: none; }
    .layer-info { flex: 1; }
    .layer-num {
      color: var(--accent); font: 12px var(--mono); font-weight: 500;
      margin-bottom: 6px;
    }
    .layer-title { font-weight: 700; font-size: 17px; margin-bottom: 4px; }
    .layer-desc { color: var(--muted); font-size: 14px; line-height: 1.6; }
    .layer-stops {
      margin-top: 8px; font: 12px var(--mono); color: var(--teal);
    }
    .status-badge {
      flex-shrink: 0; padding: 4px 12px; font: 11px var(--mono);
      text-transform: uppercase; letter-spacing: .06em; border: 1px solid;
    }
    .status-planned {
      color: var(--amber-txt); border-color: #ffcc80; background: var(--amber-bg);
    }
    .status-live {
      color: var(--green-txt); border-color: #a5d6a7; background: var(--green-bg);
    }

    /* Architecture diagram */
    .arch-diagram {
      background: #111215; color: #a1a1aa; font: 13px var(--mono);
      padding: 28px 32px; margin-top: 32px; line-height: 1.8;
      overflow-x: auto; white-space: pre;
    }
    .arch-diagram .highlight { color: #60a5fa; }
    .arch-diagram .label { color: #71717a; }

    /* Rollout steps */
    .rollout-list {
      list-style: none; padding: 0; margin: 28px 0 0; counter-reset: step;
    }
    .rollout-list li {
      counter-increment: step; display: flex; gap: 20px;
      padding: 20px 0; border-bottom: 1px solid var(--line);
      align-items: flex-start;
    }
    .rollout-list li:last-child { border-bottom: none; }
    .step-num {
      flex-shrink: 0; width: 40px; height: 40px;
      display: flex; align-items: center; justify-content: center;
      font: 14px var(--mono); font-weight: 500; color: #fff;
      background: var(--ink); border-radius: 0;
    }
    .step-num::before { content: counter(step); }
    .step-text { font-size: 16px; line-height: 1.6; padding-top: 7px; }

    /* Beyond section */
    .future-card {
      margin-top: 28px; padding: 28px 32px;
      border: 1px solid var(--line); background: var(--card);
    }
    .future-card p { color: var(--muted); line-height: 1.65; margin: 0; }
    .future-card .note {
      display: block; margin-top: 14px; font: 12px var(--mono);
      color: var(--accent); text-transform: uppercase; letter-spacing: .06em;
    }

    /* Simulator */
    .sim-section {
      margin-top: 90px; padding: 48px 0; border-top: 2px solid var(--ink);
    }
    .sim-controls {
      display: grid; grid-template-columns: 1.2fr 0.8fr 0.8fr auto;
      gap: 16px; align-items: flex-end; padding: 24px;
      background: var(--card); border: 1px solid var(--line); margin-top: 28px;
    }
    .sim-field { display: flex; flex-direction: column; gap: 6px; }
    .sim-field label {
      font: 10px var(--mono); text-transform: uppercase;
      letter-spacing: .08em; color: var(--muted); font-weight: 500;
    }
    .sim-field select, .sim-field input {
      font: 13px var(--mono); padding: 10px 12px;
      border: 1px solid var(--line); background: var(--paper);
      color: var(--ink); outline: none;
    }
    .sim-field select:focus, .sim-field input:focus {
      border-color: var(--teal);
    }
    .sim-btn {
      display: inline-flex; align-items: center; gap: 10px;
      background: var(--ink); color: #fff; border: 1px solid var(--ink);
      padding: 10px 22px; font: 12px var(--mono); font-weight: 600;
      text-transform: uppercase; letter-spacing: .06em; cursor: pointer;
      transition: background .15s; height: 42px;
    }
    .sim-btn:hover { background: #1a3035; }
    .sim-stages {
      display: grid; grid-template-columns: repeat(4, 1fr) auto;
      gap: 12px; margin: 32px 0;
    }
    .sim-stage {
      border: 1px solid var(--line); padding: 20px 14px; text-align: center;
      background: var(--card); transition: all .25s ease;
    }
    .sim-stage.lit {
      background: var(--teal); color: #fff; border-color: #0a5f5a;
    }
    .sim-stage.lit .sim-stage-caption { color: #b2dfdb; }
    .sim-stage-idx {
      font: 10px var(--mono); font-weight: 500;
      letter-spacing: .1em; margin-bottom: 6px;
    }
    .sim-stage-name { font-weight: 700; font-size: 15px; margin-bottom: 2px; }
    .sim-stage-caption { font: 10px var(--mono); color: var(--muted); }
    .sim-result {
      display: none; padding: 18px 22px; font: 13px var(--mono);
      line-height: 1.6; margin-top: 8px; border: 1px solid;
    }
    .sim-result.blocked {
      background: #fef2f2; border-color: #ef4444; color: #991b1b;
    }
    .sim-result.contained {
      background: #fff8e1; border-color: #f59e0b; color: #92400e;
    }
    .sim-result.allowed {
      background: #e8f5e9; border-color: #10b981; color: #1b5e20;
    }

    /* Footer */
    .sec-footer {
      padding: 26px 7vw; border-top: 1px solid var(--line);
      color: var(--muted); font: 11px var(--mono);
      text-transform: uppercase; letter-spacing: .08em;
    }

    @media (max-width: 760px) {
      .sec-nav { padding: 20px; }
      .sec-hero { padding: 50px 20px 0; }
      .sec-main { padding: 0 20px 70px; }
      .section { margin-top: 60px; }
      .threat-table { display: block; overflow-x: auto; }
      .layer-item { flex-direction: column; gap: 10px; }
      .nav-links { gap: 16px; }
      .nav-links a { font-size: 10px; }
      .arch-diagram { font-size: 11px; padding: 20px 16px; }
      .sim-controls { grid-template-columns: 1fr; }
      .sim-stages { grid-template-columns: repeat(5, minmax(100px, 1fr)); overflow-x: auto; }
      .sim-section { margin-top: 60px; }
    }
  </style>
</head>
<body>

  <!-- Nav -->
  <header class="sec-nav">
    <a class="logo" href="/ui">ProtonNIC Sandbox</a>
    <div class="nav-links">
      <a href="/ui">&larr; Sandbox</a>
      <a href="#simulation">Simulation</a>
      <a href="/about">About</a>
      <a href="/security" class="active">Security</a>
    </div>
  </header>

  <!-- Hero -->
  <section class="sec-hero">
    <div class="eyebrow">Security / Defense Architecture</div>
    <h1>Securing the OPI &amp; P4 Developer Sandbox</h1>
    <p class="intro">How we protect a public sandbox without relying on a third-party security service &mdash; four self-hosted defense layers built from open-source tools that we run ourselves.</p>
    <div class="version-badge">Version 1.0 &bull; November 2026 &bull; Design &amp; Rollout</div>
  </section>

  <main class="sec-main">

    <!-- What we defend against -->
    <section class="section">
      <h2>What we defend against</h2>
      <p class="section-lead">The sandbox is publicly reachable, which means it attracts scanning, brute-force attempts and floods. Each threat class is addressed by at least one defense layer.</p>

      <table class="threat-table">
        <thead>
          <tr>
            <th>Threat</th>
            <th>Example</th>
            <th>Our answer</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td class="threat-name">Traffic floods</td>
            <td class="threat-example">Thousands of HTTP or SYN requests per second</td>
            <td class="threat-answer">Firewall limits, nginx rate limits, container caps</td>
          </tr>
          <tr>
            <td class="threat-name">Key guessing</td>
            <td class="threat-example">Automated attempts at API keys</td>
            <td class="threat-answer">Long random keys, lockouts, fail2ban</td>
          </tr>
          <tr>
            <td class="threat-name">Malicious P4 programs</td>
            <td class="threat-example">Oversized or crafted uploads aimed at the compiler</td>
            <td class="threat-answer">Size limits, input checks, jailed compiles</td>
          </tr>
          <tr>
            <td class="threat-name">Exposed internals</td>
            <td class="threat-example">Switch control ports reachable from the internet</td>
            <td class="threat-answer">Internal-only network, only ports 80 and 443 public</td>
          </tr>
          <tr>
            <td class="threat-name">Supply chain</td>
            <td class="threat-example">Tampered images or dependencies</td>
            <td class="threat-answer">Pinned versions, signed images, scanning</td>
          </tr>
        </tbody>
      </table>
    </section>

    <!-- The defense layers -->
    <section class="section">
      <h2>The defense layers</h2>
      <p class="section-lead">Four concentric layers filter traffic before it reaches the application. Each one reduces the attack surface for the layer below it.</p>

      <div class="layer-stack">
        <div class="layer-item">
          <div class="layer-info">
            <div class="layer-num">Layer 01</div>
            <div class="layer-title">Host firewall &mdash; nftables</div>
            <div class="layer-desc">Blocks everything by default, limits connections per IP, enables SYN cookies to handle SYN flood attempts.</div>
            <div class="layer-stops">Stops: SYN floods, port scans, banned IPs</div>
          </div>
          <span class="status-badge status-planned">Planned</span>
        </div>
        <div class="layer-item">
          <div class="layer-info">
            <div class="layer-num">Layer 02</div>
            <div class="layer-title">Reverse proxy &mdash; nginx</div>
            <div class="layer-desc">Terminates TLS, enforces request rate limits, caps request body size, and sets aggressive timeouts for slow connections.</div>
            <div class="layer-stops">Stops: HTTP floods, slow connections, oversized uploads</div>
          </div>
          <span class="status-badge status-planned">Planned</span>
        </div>
        <div class="layer-item">
          <div class="layer-info">
            <div class="layer-num">Layer 03</div>
            <div class="layer-title">Authentication &amp; fail2ban</div>
            <div class="layer-desc">API keys required for all endpoints, per-key quotas enforced, and fail2ban automatically bans IPs after repeated failures.</div>
            <div class="layer-stops">Stops: Key guessing, anonymous abuse</div>
          </div>
          <span class="status-badge status-planned">Planned</span>
        </div>
        <div class="layer-item">
          <div class="layer-info">
            <div class="layer-num">Layer 04</div>
            <div class="layer-title">Container sandbox</div>
            <div class="layer-desc">CPU, memory and time caps on every container. P4 compiles run in a jail with no network access to contain compiler abuse.</div>
            <div class="layer-stops">Stops: Resource exhaustion, compiler abuse</div>
          </div>
          <span class="status-badge status-planned">Planned</span>
        </div>
        <div class="layer-item">
          <div class="layer-info">
            <div class="layer-num">Monitor</div>
            <div class="layer-title">Monitoring &amp; alerts</div>
            <div class="layer-desc">Metrics collection, threshold alerts and uptime checks. Does not prevent attacks directly but enables rapid detection and response.</div>
            <div class="layer-stops">Enables: Fast detection &amp; incident response</div>
          </div>
          <span class="status-badge status-planned">Planned</span>
        </div>
      </div>

      <div class="arch-diagram"><span class="label">Internet</span>
   &darr;
<span class="highlight">[ 1. Host firewall: nftables ]</span>     <span class="label">drop junk, SYN cookies, per-IP limits</span>
   &darr;
<span class="highlight">[ 2. nginx reverse proxy ]</span>         <span class="label">TLS, rate limits, body limits, timeouts</span>
   &darr;
<span class="highlight">[ 3. Auth + fail2ban ]</span>             <span class="label">API keys, quotas, automatic bans</span>
   &darr;
<span class="highlight">[ 4. Container sandbox ]</span>           <span class="label">CPU/memory/time caps, compile jail</span>
   &darr;
<span class="highlight">[ OPI server ]</span> &mdash;&mdash; <span class="highlight">[ P4 BMv2 switch ]</span>    <span class="label">internal network only</span></div>
    </section>

    <!-- Rollout plan -->
    <section class="section">
      <h2>Rollout plan</h2>
      <p class="section-lead">The layers are deployed in order, each one closing more exposure before the next one begins.</p>

      <ol class="rollout-list">
        <li>
          <div class="step-num"></div>
          <div class="step-text">Close exposure: firewall, SSH hardening, TLS, control ports internal only.</div>
        </li>
        <li>
          <div class="step-num"></div>
          <div class="step-text">Add authentication, quotas and automatic bans.</div>
        </li>
        <li>
          <div class="step-num"></div>
          <div class="step-text">Harden containers and jail P4 compilation.</div>
        </li>
        <li>
          <div class="step-num"></div>
          <div class="step-text">Add monitoring and alerts, then load-test our own limits.</div>
        </li>
        <li>
          <div class="step-num"></div>
          <div class="step-text">Secure the supply chain: signed images, pinned dependencies and repository protections.</div>
        </li>
      </ol>
    </section>

    <!-- Beyond the sandbox -->
    <section class="section">
      <h2>Beyond the sandbox</h2>
      <div class="future-card">
        <p>When SmartNIC hardware joins the platform, the same layered thinking extends into the card: secure boot, signed firmware and pipelines, locked debug ports, inline encryption and tenant isolation. That work is a later phase and nothing in this report claims it exists today.</p>
        <span class="note">Future phase &bull; Hardware security hardening</span>
      </div>
    </section>

    <!-- Attack Simulation Section -->
    <section class="section sim-section" id="simulation">
      <h2>Interactive Defense Simulation</h2>
      <p class="section-lead">Test how each layer responds to different real-world attack vectors in real time. Choose a scenario below to observe how the defense pipeline halts the threat before it reaches internal services.</p>

      <div class="sim-controls">
        <div class="sim-field">
          <label>Attack Vector</label>
          <select id="page-sec-attack">
            <option value="syn_flood">SYN Flood (10,000 req/sec)</option>
            <option value="http_flood">HTTP Slowloris / Flood</option>
            <option value="brute_force">API Key Guessing / Brute Force</option>
            <option value="malicious_p4">Malicious P4 Program Upload</option>
            <option value="port_scan">Reconnaissance Port Scan</option>
            <option value="normal">Normal Authenticated API Call</option>
          </select>
        </div>
        <div class="sim-field">
          <label>Target Service</label>
          <select id="page-sec-target">
            <option value="api">REST API (:443 Public)</option>
            <option value="bmv2">P4 BMv2 Switch (:9559 Internal)</option>
            <option value="ssh">Host SSH (:22 Internal)</option>
          </select>
        </div>
        <div class="sim-field">
          <label>Source IP</label>
          <input type="text" id="page-sec-source" value="198.51.100.77" />
        </div>
        <div>
          <button class="sim-btn" onclick="runPageSecuritySim()">
            <span>Simulate</span>
            <span>&rarr;</span>
          </button>
        </div>
      </div>

      <div class="sim-stages">
        <div class="sim-stage" id="psim-stage-1">
          <div class="sim-stage-idx">LAYER 01</div>
          <div class="sim-stage-name">nftables</div>
          <div class="sim-stage-caption">Host Firewall</div>
        </div>
        <div class="sim-stage" id="psim-stage-2">
          <div class="sim-stage-idx">LAYER 02</div>
          <div class="sim-stage-name">nginx</div>
          <div class="sim-stage-caption">Reverse Proxy</div>
        </div>
        <div class="sim-stage" id="psim-stage-3">
          <div class="sim-stage-idx">LAYER 03</div>
          <div class="sim-stage-name">Auth / fail2ban</div>
          <div class="sim-stage-caption">Access Control</div>
        </div>
        <div class="sim-stage" id="psim-stage-4">
          <div class="sim-stage-idx">LAYER 04</div>
          <div class="sim-stage-name">Sandbox Jail</div>
          <div class="sim-stage-caption">Container Sandbox</div>
        </div>
        <div class="sim-stage" id="psim-stage-5">
          <div class="sim-stage-idx">TARGET</div>
          <div class="sim-stage-name">OPI Server</div>
          <div class="sim-stage-caption">Dataplane Backend</div>
        </div>
      </div>

      <div class="sim-result" id="page-sim-result"></div>
    </section>

  </main>

  <footer class="sec-footer">OpenLabNetworks &bull; Accelerating cloud-native infrastructure</footer>

  <script>
    var pageScenarios = {
      syn_flood: {
        block: 1,
        cssClass: 'blocked',
        verdict: 'BLOCKED AT LAYER 1 (HOST FIREWALL)',
        details: 'nftables SYN cookies engaged and connection rate limit per IP triggered. Attack traffic dropped before kernel socket allocation.'
      },
      http_flood: {
        block: 2,
        cssClass: 'blocked',
        verdict: 'BLOCKED AT LAYER 2 (NGINX REVERSE PROXY)',
        details: 'nginx rate limiter (zone per-IP) exceeded limit_req threshold; client connection timed out and 429 Too Many Requests returned.'
      },
      brute_force: {
        block: 3,
        cssClass: 'blocked',
        verdict: 'BLOCKED AT LAYER 3 (AUTHENTICATION &amp; FAIL2BAN)',
        details: '5 consecutive invalid API key attempts detected. fail2ban dynamically injected temporary drop rule for origin IP.'
      },
      malicious_p4: {
        block: 4,
        cssClass: 'contained',
        verdict: 'CONTAINED AT LAYER 4 (CONTAINER SANDBOX)',
        details: 'Malformed/oversized P4 compile rejected: strict payload limit (1MB max), CPU/memory cgroups enforced, compile jail executed with no network access.'
      },
      port_scan: {
        block: 1,
        cssClass: 'blocked',
        verdict: 'BLOCKED AT LAYER 1 (HOST FIREWALL)',
        details: 'nftables default-drop rule dropped probe packet silently. Internal switch control ports (:9559, :9090) are completely non-routable.'
      },
      normal: {
        block: 0,
        cssClass: 'allowed',
        verdict: 'ALLOWED &bull; ALL LAYERS PASSED',
        details: 'Valid API token verified, rate limits respected. Request safely forwarded to internal OPI server daemon.'
      }
    };

    async function runPageSecuritySim() {
      var attack = document.getElementById('page-sec-attack').value;
      var target = document.getElementById('page-sec-target').value;
      var banner = document.getElementById('page-sim-result');
      banner.style.display = 'none';

      var scenario = pageScenarios[attack];
      if (target === 'bmv2' || target === 'ssh') {
        scenario = {
          block: 1,
          cssClass: 'blocked',
          verdict: 'BLOCKED AT LAYER 1 (HOST FIREWALL)',
          details: 'Default-deny policy enforced: Target port ' + (target === 'bmv2' ? ':9559' : ':22') + ' is not exposed to the public internet.'
        };
      }

      var total = 5;
      for (var i = 1; i <= total; i++) {
        var el = document.getElementById('psim-stage-' + i);
        el.classList.remove('lit');
        el.style.background = '';
        el.style.borderColor = '';
        el.style.color = '';
      }

      for (var s = 1; s <= total; s++) {
        var stage = document.getElementById('psim-stage-' + s);
        stage.classList.add('lit');
        await new Promise(function(resolve) { setTimeout(resolve, 240); });

        if (scenario.block > 0 && s === scenario.block) {
          stage.style.background = '#c0392b';
          stage.style.borderColor = '#962d22';
          stage.style.color = '#ffffff';
          break;
        }

        if (s < total || scenario.block === 0) {
          stage.classList.remove('lit');
        }
      }

      if (scenario.block === 0) {
        var targetStage = document.getElementById('psim-stage-5');
        targetStage.style.background = '#10b981';
        targetStage.style.borderColor = '#059669';
        targetStage.style.color = '#ffffff';
        targetStage.classList.add('lit');
      }

      banner.className = 'sim-result ' + scenario.cssClass;
      banner.innerHTML = '<strong>[' + scenario.verdict + ']</strong><br>' + scenario.details;
      banner.style.display = 'block';

      setTimeout(function() {
        for (var k = 1; k <= total; k++) {
          var stg = document.getElementById('psim-stage-' + k);
          stg.classList.remove('lit');
          stg.style.background = '';
          stg.style.borderColor = '';
          stg.style.color = '';
        }
      }, 4000);
    }
  </script>

</body>
</html>
`
