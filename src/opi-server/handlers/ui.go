package handlers

import (
	"net/http"
)

// ServeUI renders the minimalist, editorial light-themed testing dashboard
// for the OPENLABNETWORKS Open Sandbox.
func ServeUI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(dashboardHTML))
}

const dashboardHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>OPENLABNETWORKS | Open Sandbox</title>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;700&family=Space+Grotesk:wght@500;700;900&family=Inter:wght@400;500;600&display=swap" rel="stylesheet">
  <style>
    /* -------------------------------------------------------------
       MINIMAL EDITORIAL LIGHT THEME
       Inspired by minimalist Swiss / Brutalist layout & clean typography
    ------------------------------------------------------------- */
    :root {
      --bg: #f4f5f7;
      --surface: #ffffff;
      --ink: #111111;
      --ink-muted: #666666;
      --ink-subtle: #999999;
      --line: #e2e4e8;
      --line-dark: #111111;
      --cobalt: #2563eb;
      --cobalt-dark: #1d4ed8;
      --green: #10b981;
      --red: #ef4444;
      --font-display: 'Space Grotesk', -apple-system, BlinkMacSystemFont, sans-serif;
      --font-body: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
      --font-mono: 'JetBrains Mono', monospace;
    }

    * {
      box-sizing: border-box;
      margin: 0;
      padding: 0;
    }

    body {
      background-color: var(--bg);
      background-image: 
        linear-gradient(to right, rgba(0, 0, 0, 0.035) 1px, transparent 1px),
        linear-gradient(to bottom, rgba(0, 0, 0, 0.035) 1px, transparent 1px);
      background-size: 44px 44px;
      color: var(--ink);
      font-family: var(--font-body);
      line-height: 1.5;
      -webkit-font-smoothing: antialiased;
      overflow-x: hidden;
    }

    /* Minimal Top Navigation */
    header.site-header {
      background: var(--bg);
      border-bottom: 1px solid var(--line-dark);
      padding: 20px 48px;
      display: flex;
      justify-content: space-between;
      align-items: center;
    }

    .brand {
      font-family: var(--font-display);
      font-weight: 900;
      font-size: 16px;
      letter-spacing: -0.03em;
      text-transform: uppercase;
      color: var(--ink);
      text-decoration: none;
      display: flex;
      align-items: center;
      gap: 0;
      white-space: nowrap;
    }
    .brand-dot {
      color: var(--cobalt);
    }

    nav.site-nav {
      display: flex;
      gap: 36px;
      font-family: var(--font-mono);
      font-size: 11px;
      text-transform: uppercase;
      letter-spacing: 0.1em;
    }
    nav.site-nav a {
      color: var(--ink-muted);
      text-decoration: none;
      transition: color 0.15s ease;
    }
    nav.site-nav a:hover {
      color: var(--ink);
    }

    .btn {
      display: inline-flex;
      align-items: center;
      gap: 10px;
      background: var(--ink);
      color: #ffffff;
      border: 1px solid var(--ink);
      padding: 10px 22px;
      font-family: var(--font-mono);
      font-size: 11px;
      font-weight: 600;
      text-transform: uppercase;
      letter-spacing: 0.08em;
      cursor: pointer;
      text-decoration: none;
      transition: all 0.15s ease;
    }
    .btn:hover {
      background: #262626;
      border-color: #262626;
    }
    .btn-ghost {
      background: transparent;
      color: var(--ink);
      border: 1px solid var(--line-dark);
    }
    .btn-ghost:hover {
      background: rgba(0,0,0,0.04);
    }
    .btn-sm {
      padding: 6px 14px;
      font-size: 10px;
    }
    .btn-circle-arrow {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      width: 18px;
      height: 18px;
      border-radius: 50%;
      border: 1px solid rgba(255,255,255,0.4);
      font-size: 10px;
    }

    /* Sub-header rule line */
    .meta-rule-bar {
      border-bottom: 1px solid var(--line-dark);
      padding: 10px 48px;
      display: flex;
      justify-content: space-between;
      font-family: var(--font-mono);
      font-size: 10px;
      text-transform: uppercase;
      letter-spacing: 0.12em;
      color: var(--ink-muted);
    }

    /* HERO SECTION */
    .hero-wrap {
      display: grid;
      grid-template-columns: 1.3fr 0.9fr;
      border-bottom: 1px solid var(--line-dark);
      min-height: 520px;
      position: relative;
    }

    .hero-content {
      padding: 64px 48px;
      display: flex;
      flex-direction: column;
      justify-content: space-between;
      position: relative;
    }

    .hero-title {
      font-family: var(--font-display);
      font-weight: 900;
      text-transform: uppercase;
      line-height: 0.88;
      letter-spacing: -0.05em;
      margin-bottom: 28px;
    }
    .hero-title .title-black {
      font-size: clamp(72px, 9.5vw, 128px);
      color: var(--ink);
      display: block;
    }
    .hero-title .title-stroke {
      font-size: clamp(72px, 9.5vw, 128px);
      color: #ffffff;
      -webkit-text-stroke: 3px var(--ink);
      text-shadow: 4px 4px 0 var(--ink);
      display: block;
      margin-top: 2px;
    }

    .hero-lead {
      display: flex;
      align-items: center;
      gap: 10px;
      font-family: var(--font-mono);
      font-size: 12px;
      color: var(--ink);
      margin-bottom: 32px;
    }
    .dot-online {
      width: 8px;
      height: 8px;
      border-radius: 50%;
      background: var(--green);
      display: inline-block;
      flex-shrink: 0;
    }

    .hero-cta {
      display: flex;
      align-items: center;
      gap: 16px;
      margin-bottom: 40px;
    }

    .hero-footer-meta {
      display: flex;
      justify-content: space-between;
      align-items: center;
      font-family: var(--font-mono);
      font-size: 11px;
      color: var(--ink-muted);
      border-top: 1px solid var(--line);
      padding-top: 18px;
    }

    /* Right Split: Cobalt Diagonal + Floating Card */
    .hero-visual {
      background: var(--cobalt);
      clip-path: polygon(16% 0%, 100% 0%, 100% 100%, 0% 100%);
      display: flex;
      align-items: center;
      justify-content: center;
      padding: 48px;
      position: relative;
    }

    .hardware-card {
      background: #111215;
      color: #ffffff;
      border: 1px solid #22242a;
      width: 100%;
      max-width: 340px;
      padding: 36px 30px;
      position: relative;
      font-family: var(--font-mono);
      box-shadow: 0 20px 40px rgba(0, 0, 0, 0.25);
    }
    .hardware-card .punch-l {
      position: absolute;
      left: -9px;
      top: 50%;
      transform: translateY(-50%);
      width: 18px;
      height: 18px;
      background: var(--cobalt);
      border-radius: 50%;
    }
    .hardware-card .punch-r {
      position: absolute;
      right: -9px;
      top: 50%;
      transform: translateY(-50%);
      width: 18px;
      height: 18px;
      background: var(--cobalt);
      border-radius: 50%;
    }
    .card-top-tag {
      font-size: 9px;
      letter-spacing: 0.16em;
      text-transform: uppercase;
      color: #71717a;
      border-bottom: 1px solid #27272a;
      padding-bottom: 10px;
      margin-bottom: 18px;
      display: flex;
      justify-content: space-between;
    }
    .card-giant-text {
      font-family: var(--font-display);
      font-size: 48px;
      font-weight: 900;
      line-height: 0.9;
      letter-spacing: -0.04em;
      color: #ffffff;
      margin-bottom: 4px;
    }
    .card-accent-text {
      font-family: var(--font-display);
      font-size: 40px;
      font-weight: 900;
      color: #60a5fa;
      letter-spacing: -0.04em;
      margin-bottom: 24px;
    }
    .card-spec-row {
      display: flex;
      justify-content: space-between;
      border-top: 1px solid #27272a;
      padding-top: 14px;
      margin-top: 14px;
      font-size: 10px;
      color: #a1a1aa;
    }
    .card-spec-row b {
      color: #ffffff;
      font-weight: 500;
    }

    /* MINIMAL TABBED STUDIO WORKBENCH */
    .studio-wrap {
      max-width: 1280px;
      margin: 0 auto;
      padding: 64px 48px 96px 48px;
    }

    .studio-header {
      display: flex;
      justify-content: space-between;
      align-items: flex-end;
      border-bottom: 1px solid var(--line-dark);
      padding-bottom: 16px;
      margin-bottom: 32px;
    }
    .studio-title {
      font-family: var(--font-display);
      font-size: 24px;
      font-weight: 900;
      text-transform: uppercase;
      letter-spacing: -0.02em;
    }

    /* Clean Minimalist Tab Bar */
    .tab-bar {
      display: flex;
      gap: 8px;
      font-family: var(--font-mono);
      font-size: 11px;
      text-transform: uppercase;
      letter-spacing: 0.08em;
    }
    .tab-btn {
      background: transparent;
      border: 1px solid transparent;
      color: var(--ink-muted);
      padding: 8px 16px;
      cursor: pointer;
      font-family: inherit;
      font-size: inherit;
      transition: all 0.15s ease;
    }
    .tab-btn:hover {
      color: var(--ink);
    }
    .tab-btn.active {
      background: var(--ink);
      color: #ffffff;
      border-color: var(--ink);
      font-weight: 700;
    }

    /* Tab Panes */
    .tab-pane {
      display: none;
      background: var(--surface);
      border: 1px solid var(--line-dark);
      padding: 36px;
      animation: fadeIn 0.2s ease;
    }
    .tab-pane.active {
      display: block;
    }
    @keyframes fadeIn {
      from { opacity: 0; transform: translateY(4px); }
      to { opacity: 1; transform: translateY(0); }
    }

    .pane-toolbar {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 24px;
      padding-bottom: 16px;
      border-bottom: 1px solid var(--line);
    }
    .pane-desc {
      font-family: var(--font-mono);
      font-size: 11px;
      color: var(--ink-muted);
      text-transform: uppercase;
      letter-spacing: 0.06em;
    }

    /* Minimal Tables */
    table.min-table {
      width: 100%;
      border-collapse: collapse;
      font-family: var(--font-mono);
      font-size: 12px;
      text-align: left;
    }
    table.min-table th {
      border-bottom: 1px solid var(--line-dark);
      padding: 12px 14px;
      font-size: 10px;
      text-transform: uppercase;
      letter-spacing: 0.08em;
      color: var(--ink-muted);
      font-weight: 700;
    }
    table.min-table td {
      border-bottom: 1px solid var(--line);
      padding: 14px;
      color: var(--ink);
    }
    table.min-table tr:hover td {
      background: #fafafa;
    }

    .tag-state {
      display: inline-block;
      padding: 2px 6px;
      font-size: 10px;
      font-weight: 700;
      letter-spacing: 0.05em;
      border: 1px solid var(--line-dark);
    }
    .tag-state.up {
      background: #e6f4ea;
      color: #137333;
      border-color: #34a853;
    }

    /* Minimal P4 Simulator */
    .sim-strip {
      display: grid;
      grid-template-columns: repeat(4, 1fr) auto;
      gap: 16px;
      align-items: flex-end;
      background: #f9fafb;
      border: 1px solid var(--line);
      padding: 20px;
      margin-bottom: 32px;
    }
    .field-box {
      display: flex;
      flex-direction: column;
      gap: 6px;
    }
    .field-box label {
      font-family: var(--font-mono);
      font-size: 10px;
      text-transform: uppercase;
      letter-spacing: 0.08em;
      color: var(--ink-muted);
      font-weight: 700;
    }
    .field-box select, .field-box input {
      font-family: var(--font-mono);
      font-size: 12px;
      padding: 8px 10px;
      border: 1px solid var(--line-dark);
      background: #fff;
      color: var(--ink);
      outline: none;
    }

    .sim-stages-row {
      display: grid;
      grid-template-columns: repeat(5, 1fr);
      gap: 12px;
      margin: 28px 0;
    }
    .stage-item {
      border: 1px solid var(--line-dark);
      padding: 16px 12px;
      text-align: center;
      background: #ffffff;
      transition: all 0.2s ease;
    }
    .stage-item.lit {
      background: var(--cobalt);
      color: #ffffff;
      border-color: var(--cobalt-dark);
    }
    .stage-item.lit .stage-caption {
      color: #dbeafe;
    }
    .stage-idx {
      font-family: var(--font-mono);
      font-size: 9px;
      font-weight: 700;
      letter-spacing: 0.1em;
      margin-bottom: 6px;
    }
    .stage-label {
      font-family: var(--font-display);
      font-size: 13px;
      font-weight: 700;
      text-transform: uppercase;
      margin-bottom: 2px;
    }
    .stage-caption {
      font-family: var(--font-mono);
      font-size: 9px;
      color: var(--ink-subtle);
    }

    .sim-banner {
      display: none;
      padding: 16px 20px;
      border: 1px solid var(--line-dark);
      font-family: var(--font-mono);
      font-size: 12px;
      margin-top: 24px;
    }
    .sim-banner.pass {
      background: #f0fdf4;
      border-color: #16a34a;
      color: #166534;
    }
    .sim-banner.drop {
      background: #fef2f2;
      border-color: #dc2626;
      color: #991b1b;
    }

    /* Minimal API Console */
    .console-split {
      display: grid;
      grid-template-columns: 280px 1fr;
      gap: 20px;
    }
    .api-list {
      display: flex;
      flex-direction: column;
      gap: 8px;
    }
    .api-item {
      display: flex;
      align-items: center;
      justify-content: space-between;
      border: 1px solid var(--line-dark);
      padding: 10px 14px;
      background: #ffffff;
      font-family: var(--font-mono);
      font-size: 11px;
      cursor: pointer;
      text-align: left;
      transition: all 0.15s ease;
    }
    .api-item:hover {
      background: var(--ink);
      color: #ffffff;
    }
    .term-box {
      background: #111215;
      color: #38bdf8;
      border: 1px solid #22242a;
      font-family: var(--font-mono);
      font-size: 11px;
      padding: 16px;
      height: 320px;
      overflow-y: auto;
      overflow-x: auto;
      white-space: pre-wrap;
      overflow-wrap: anywhere;
      line-height: 1.6;
    }

    /* Modal */
    .modal-backdrop {
      display: none;
      position: fixed;
      top: 0;
      left: 0;
      width: 100vw;
      height: 100vh;
      background: rgba(0, 0, 0, 0.4);
      backdrop-filter: blur(2px);
      z-index: 1000;
      align-items: center;
      justify-content: center;
    }
    .modal-box {
      background: #ffffff;
      border: 1px solid var(--line-dark);
      width: 90%;
      max-width: 440px;
      padding: 28px;
      box-shadow: 0 16px 32px rgba(0,0,0,0.15);
    }
    .modal-head {
      font-family: var(--font-display);
      font-size: 16px;
      font-weight: 900;
      text-transform: uppercase;
      margin-bottom: 20px;
      border-bottom: 1px solid var(--line-dark);
      padding-bottom: 8px;
    }

    /* Footer */
    footer.site-footer {
      border-top: 1px solid var(--line-dark);
      padding: 32px 48px;
      display: flex;
      justify-content: space-between;
      gap: 16px;
      font-family: var(--font-mono);
      font-size: 10px;
      color: var(--ink-muted);
      text-transform: uppercase;
      letter-spacing: 0.1em;
    }

    @media (max-width: 920px) {
      header.site-header {
        flex-wrap: wrap;
        gap: 16px 24px;
      }
      nav.site-nav {
        order: 3;
        width: 100%;
        justify-content: flex-start;
        overflow-x: auto;
        padding-bottom: 4px;
      }
      .hero-wrap {
        grid-template-columns: 1fr;
      }
      .hero-visual {
        clip-path: none;
        min-height: 360px;
      }
      .sim-strip {
        grid-template-columns: 1fr;
      }
      .sim-stages-row {
        grid-template-columns: repeat(5, minmax(120px, 1fr));
        overflow-x: auto;
      }
      .console-split {
        grid-template-columns: 1fr;
      }
    }

    @media (max-width: 640px) {
      body {
        background-size: 28px 28px;
      }
      header.site-header,
      .meta-rule-bar {
        padding-left: 20px;
        padding-right: 20px;
      }
      header.site-header {
        align-items: flex-start;
      }
      .brand {
        font-size: 17px;
      }
      header.site-header > div {
        margin-left: auto;
      }
      nav.site-nav {
        gap: 20px;
        font-size: 10px;
        scrollbar-width: none;
      }
      nav.site-nav::-webkit-scrollbar {
        display: none;
      }
      .meta-rule-bar {
        align-items: flex-start;
        flex-direction: column;
        gap: 5px;
        font-size: 9px;
        line-height: 1.4;
      }
      .hero-content {
        padding: 40px 20px 28px;
      }
      .hero-title .title-black,
      .hero-title .title-stroke {
        font-size: clamp(54px, 18vw, 92px);
      }
      .hero-title .title-stroke {
        -webkit-text-stroke-width: 2px;
        text-shadow: 3px 3px 0 var(--ink);
      }
      .hero-lead {
        align-items: flex-start;
        font-size: 11px;
        margin-bottom: 24px;
      }
      .hero-cta {
        align-items: stretch;
        flex-direction: column;
        gap: 10px;
        margin-bottom: 28px;
      }
      .hero-cta .btn {
        justify-content: space-between;
        width: 100%;
      }
      .hero-footer-meta {
        align-items: flex-start;
        flex-direction: column;
        gap: 8px;
        font-size: 9px;
      }
      .hero-visual {
        min-height: 330px;
        padding: 28px 20px;
      }
      .hardware-card {
        padding: 28px 22px;
      }
      .card-giant-text {
        font-size: 40px;
      }
      .card-accent-text {
        font-size: 34px;
        margin-bottom: 20px;
      }
      .studio-wrap {
        padding: 40px 20px 64px;
      }
      .studio-header {
        align-items: flex-start;
        flex-direction: column;
        gap: 18px;
        margin-bottom: 20px;
      }
      .studio-title {
        font-size: 20px;
      }
      .tab-bar {
        width: calc(100vw - 40px);
        gap: 4px;
        overflow-x: auto;
        padding-bottom: 5px;
        scrollbar-width: none;
      }
      .tab-bar::-webkit-scrollbar {
        display: none;
      }
      .tab-btn {
        flex: 0 0 auto;
        padding: 8px 12px;
        font-size: 10px;
      }
      .tab-pane {
        padding: 20px 16px;
      }
      .pane-toolbar {
        align-items: flex-start;
        flex-direction: column;
        gap: 14px;
      }
      .pane-toolbar > div:last-child {
        width: 100%;
        flex-wrap: wrap;
      }
      .pane-toolbar > div:last-child .btn {
        flex: 1 1 auto;
      }
      .pane-desc {
        font-size: 10px;
      }
      table.min-table {
        display: block;
        overflow-x: auto;
        white-space: nowrap;
        -webkit-overflow-scrolling: touch;
      }
      table.min-table th,
      table.min-table td {
        padding: 12px 10px;
      }
      .sim-strip {
        gap: 12px;
        padding: 16px;
      }
      .sim-strip > div:last-child {
        width: 100%;
      }
      .sim-strip > div:last-child .btn {
        width: 100%;
        justify-content: space-between;
      }
      .sim-stages-row {
        gap: 8px;
        margin: 20px 0;
      }
      .stage-item {
        padding: 14px 10px;
      }
      .console-split {
        gap: 14px;
      }
      .term-box {
        height: 260px;
        padding: 12px;
        font-size: 10px;
      }
      .modal-box {
        width: calc(100% - 32px);
        max-height: calc(100vh - 32px);
        overflow-y: auto;
        padding: 22px 18px;
      }
      footer.site-footer {
        align-items: flex-start;
        flex-direction: column;
        padding: 24px 20px;
        font-size: 9px;
      }
    }

    @media (max-width: 380px) {
      header.site-header > div {
        width: 100%;
        margin-left: 0;
      }
      header.site-header > div .btn {
        width: 100%;
        justify-content: space-between;
      }
      .hero-title .title-black,
      .hero-title .title-stroke {
        font-size: 52px;
      }
      .hardware-card {
        padding: 24px 18px;
      }
      .card-spec-row {
        font-size: 9px;
      }
    }
  </style>
</head>
<body>

  <!-- Minimal Nav -->
  <header class="site-header">
    <a href="https://www.openlabnetworks.com/" class="brand" target="_blank" rel="noopener noreferrer">
      <span>OPENLABNETWORKS</span>
    </a>
    <nav class="site-nav">
      <a href="javascript:void(0)" onclick="selectTab('ports')">Ports</a>
      <a href="javascript:void(0)" onclick="selectTab('pipelines')">Pipelines</a>
      <a href="javascript:void(0)" onclick="selectTab('simulator')">Offload Lab</a>
      <a href="javascript:void(0)" onclick="selectTab('console')">APIs</a>
      <a href="javascript:void(0)" onclick="selectTab('ami')">AMI</a>
      <a href="/about">About</a>
    </nav>
    <div>
      <button class="btn btn-sm" onclick="runSmokeTest()">
        <span>Smoke Test</span>
        <span class="btn-circle-arrow">&rarr;</span>
      </button>
    </div>
  </header>

  <!-- Editorial Meta Bar -->
  <div class="meta-rule-bar">
    <span>OPENLABNETWORKS &mdash; HARDWARE OFFLOAD</span>
    <span>DPU / SMARTNIC READY</span>
    <span>STATUS: ONLINE &bull; 0ms LATENCY</span>
  </div>

  <!-- HERO SECTION -->
  <section class="hero-wrap">
    <!-- Left Hero -->
    <div class="hero-content">
      <div style="font-family: var(--font-mono); font-size: 10px; letter-spacing: 0.15em; text-transform: uppercase; color: var(--ink-muted); margin-bottom: 16px;">
        PROGRAMMABLE INFRASTRUCTURE // LINUX // 2026
      </div>

      <div>
        <div class="hero-title">
          <span class="title-black">OPEN</span>
          <span class="title-stroke">SANDBOX</span>
        </div>

        <div class="hero-lead">
          <span class="dot-online"></span>
          <span>Zero hardware dependency. Instant DPU dataplane simulation.</span>
        </div>

        <div class="hero-cta">
          <button class="btn" onclick="scrollToStudio()">
            <span>Open Workbench</span>
            <span class="btn-circle-arrow">&rarr;</span>
          </button>
          <button class="btn btn-ghost" onclick="selectTab('simulator'); scrollToStudio();">
            <span>Explore Offload</span>
            <span style="font-size: 10px;">&darr;</span>
          </button>
        </div>
      </div>

      <div class="hero-footer-meta">
        <div>01 &mdash; 05 // SYSTEM TEST</div>
        <div>{ REST METADATA // BMV2 CONTROL SEPARATE }</div>
      </div>
    </div>

    <!-- Right Hero: Cobalt Diagonal + Floating Card -->
    <div class="hero-visual">
      <div class="hardware-card">
        <div class="punch-l"></div>
        <div class="punch-r"></div>

        <div class="card-top-tag">
          <span>DPU HARDWARE PASS</span>
          <span>DEV #01</span>
        </div>

        <div class="card-giant-text">100G</div>
        <div class="card-accent-text">BMV2</div>

        <div class="card-spec-row">
          <span>TARGET</span>
          <b>simple_switch_grpc</b>
        </div>
        <div class="card-spec-row">
          <span>DEVICE ID</span>
          <b>0</b>
        </div>
        <div class="card-spec-row">
          <span>P4RUNTIME</span>
          <b>:9559 gRPC</b>
        </div>
        <div class="card-spec-row">
          <span>REST API</span>
          <b>:8080 Active</b>
        </div>
      </div>
    </div>
  </section>

  <!-- INTERACTIVE WORKBENCH -->
  <section class="studio-wrap" id="workbench">
    <div class="studio-header">
      <div class="studio-title">Testing Workbench</div>
      <div class="tab-bar">
        <button class="tab-btn active" id="tab-btn-ports" onclick="selectTab('ports')">01 Ports</button>
        <button class="tab-btn" id="tab-btn-pipelines" onclick="selectTab('pipelines')">02 Pipelines</button>
        <button class="tab-btn" id="tab-btn-simulator" onclick="selectTab('simulator')">03 Offload Lab</button>
        <button class="tab-btn" id="tab-btn-console" onclick="selectTab('console')">04 REST Console</button>
        <button class="tab-btn" id="tab-btn-ami" onclick="selectTab('ami')">05 AWS AMI</button>
      </div>
    </div>

    <!-- TAB 1: PORTS -->
    <div class="tab-pane active" id="pane-ports">
      <div class="pane-toolbar">
        <div class="pane-desc">Network Ports &bull; GET /v1/ports</div>
        <div style="display:flex; gap:8px;">
          <button class="btn btn-ghost btn-sm" onclick="fetchPorts()">Refresh</button>
          <button class="btn btn-sm" onclick="openModal('modal-port')">+ Create Port</button>
        </div>
      </div>

      <table class="min-table">
        <thead>
          <tr>
            <th>Port ID</th>
            <th>Name</th>
            <th>MAC Address</th>
            <th>MTU</th>
            <th>Speed</th>
            <th>State</th>
            <th>Action</th>
          </tr>
        </thead>
        <tbody id="ports-tbody">
          <tr><td colspan="7" style="text-align:center; padding:24px; color:var(--ink-muted);">Loading ports...</td></tr>
        </tbody>
      </table>
    </div>

    <!-- TAB 2: PIPELINES -->
    <div class="tab-pane" id="pane-pipelines">
      <div class="pane-toolbar">
        <div class="pane-desc">Dataplane Workflows &bull; GET /v1/pipelines</div>
        <div style="display:flex; gap:8px;">
          <button class="btn btn-ghost btn-sm" onclick="fetchPipelines()">Refresh</button>
          <button class="btn btn-sm" onclick="openModal('modal-pipeline')">+ Register Pipeline</button>
        </div>
      </div>

      <table class="min-table">
        <thead>
          <tr>
            <th>Pipeline ID</th>
            <th>Workflow</th>
            <th>Status</th>
            <th>Action</th>
          </tr>
        </thead>
        <tbody id="pipelines-tbody">
          <tr><td colspan="4" style="text-align:center; padding:24px; color:var(--ink-muted);">Loading pipelines...</td></tr>
        </tbody>
      </table>
    </div>

    <!-- TAB 3: HARDWARE OFFLOAD LAB -->
    <div class="tab-pane" id="pane-simulator">
      <div class="pane-toolbar">
        <div class="pane-desc">V1Model 5-Stage Dataplane Traversal</div>
      </div>

      <div class="sim-strip">
        <div class="field-box">
          <label>Dataplane workflow</label>
          <select id="sim-program">
            <option value="basic_forwarding.p4">basic_forwarding.p4 (L2/L3 Routing)</option>
            <option value="firewall_drop.p4">firewall_drop.p4 (IPv4 Filter)</option>
          </select>
        </div>
        <div class="field-box">
          <label>Ingress Port</label>
          <select id="sim-ingress">
            <option value="port-0">port-0 (100G)</option>
            <option value="port-1">port-1 (100G)</option>
          </select>
        </div>
        <div class="field-box">
          <label>EtherType</label>
          <select id="sim-ethertype">
            <option value="0x0800">0x0800 (IPv4)</option>
            <option value="0x0806">0x0806 (ARP / Non-IPv4)</option>
          </select>
        </div>
        <div class="field-box">
          <label>Destination</label>
          <input type="text" id="sim-dst" value="10.0.1.10:80">
        </div>
        <div>
          <button class="btn btn-sm" style="height:36px;" onclick="runSimulation()">
            <span>Simulate</span>
            <span class="btn-circle-arrow">&rarr;</span>
          </button>
        </div>
      </div>

      <div class="sim-stages-row">
        <div class="stage-item" id="stage-1">
          <div class="stage-idx">STAGE 01</div>
          <div class="stage-label">Parser</div>
          <div class="stage-caption">Ethernet &amp; IPv4</div>
        </div>
        <div class="stage-item" id="stage-2">
          <div class="stage-idx">STAGE 02</div>
          <div class="stage-label">Checksum</div>
          <div class="stage-caption">RFC 1071</div>
        </div>
        <div class="stage-item" id="stage-3">
          <div class="stage-idx">STAGE 03</div>
          <div class="stage-label">Ingress</div>
          <div class="stage-caption">Match &amp; Action</div>
        </div>
        <div class="stage-item" id="stage-4">
          <div class="stage-idx">STAGE 04</div>
          <div class="stage-label">Egress</div>
          <div class="stage-caption">Queuing</div>
        </div>
        <div class="stage-item" id="stage-5">
          <div class="stage-idx">STAGE 05</div>
          <div class="stage-label">Deparser</div>
          <div class="stage-caption">Emit Packet</div>
        </div>
      </div>

      <div class="sim-banner" id="sim-banner"></div>
    </div>

    <!-- TAB 4: REST CONSOLE -->
    <div class="tab-pane" id="pane-console">
      <div class="pane-toolbar">
        <div class="pane-desc">Interactive REST Calls to :8080</div>
        <div id="curl-hint" style="font-family:var(--font-mono); font-size:10px; color:var(--cobalt);">curl http://localhost:8080/healthz</div>
      </div>

      <div class="console-split">
        <div class="api-list">
          <button class="api-item" onclick="callApi('GET', '/healthz')">
            <span>GET /healthz</span>
            <span style="color:#10b981; font-weight:700;">&bull;</span>
          </button>
          <button class="api-item" onclick="callApi('GET', '/v1/ports')">
            <span>GET /v1/ports</span>
            <span>&rarr;</span>
          </button>
          <button class="api-item" onclick="callApi('POST', '/v1/ports', {mac_address: randomMac(), mtu: 1500})">
            <span>POST /v1/ports</span>
            <span>+</span>
          </button>
          <button class="api-item" onclick="callApi('GET', '/v1/pipelines')">
            <span>GET /v1/pipelines</span>
            <span>&rarr;</span>
          </button>
          <button class="api-item" onclick="callApi('POST', '/v1/pipelines', {p4_program: 'firewall_drop.p4'})">
            <span>POST /v1/pipelines</span>
            <span>+</span>
          </button>
        </div>

        <div class="term-box" id="term-box">// Output Console Ready. Click any endpoint on the left...</div>
      </div>
    </div>

    <!-- TAB 5: AWS AMI -->
    <div class="tab-pane" id="pane-ami">
      <div class="pane-toolbar">
        <div class="pane-desc">Option 1: Pre-Baked Custom AMI (Packer)</div>
      </div>

      <p style="font-family:var(--font-mono); font-size:12px; line-height:1.6; margin-bottom:16px;">
        Zero download time at boot. Pre-bakes Docker, BMv2, and OPI images directly into the AMI filesystem.
      </p>

      <div style="background:#111215; color:#e4e4e7; padding:16px; font-family:var(--font-mono); font-size:11px; margin-bottom:16px; line-height:1.5;">
# 1. Initialize Packer plugin<br>
packer init packer/opi-sandbox.pkr.hcl<br><br>
# 2. Build the AMI (AWS credentials required)<br>
packer build -var="aws_region=us-east-1" packer/opi-sandbox.pkr.hcl<br><br>
# 3. Launch from AWS Console -&gt; "My AMIs" -&gt; opi-p4-sandbox-v1.0.0<br>
# Service auto-starts on boot via /etc/systemd/system/opi-sandbox.service
      </div>

      <div style="font-family:var(--font-mono); font-size:11px; color:var(--ink-muted);">
        Template file: <code>packer/opi-sandbox.pkr.hcl</code> &bull; Systemd: <code>/etc/systemd/system/opi-sandbox.service</code>
      </div>
    </div>
  </section>

  <!-- MODAL: ADD PORT -->
  <div class="modal-backdrop" id="modal-port">
    <div class="modal-box">
      <div class="modal-head">Create Virtual Port</div>
      <form onsubmit="submitPort(event)">
        <div class="field-box" style="margin-bottom:14px;">
          <label>MAC Address</label>
          <input type="text" id="m-mac" required>
          <small style="font-family:var(--font-mono); font-size:10px; margin-top:4px;">
            <a href="javascript:void(0)" onclick="document.getElementById('m-mac').value = randomMac()" style="color:var(--cobalt)">Generate MAC</a>
          </small>
        </div>
        <div class="field-box" style="margin-bottom:20px;">
          <label>MTU</label>
          <input type="number" id="m-mtu" value="1500" required>
        </div>
        <div style="display:flex; justify-content:flex-end; gap:8px;">
          <button type="button" class="btn btn-ghost btn-sm" onclick="closeModal('modal-port')">Cancel</button>
          <button type="submit" class="btn btn-sm">Create</button>
        </div>
      </form>
    </div>
  </div>

  <!-- MODAL: REGISTER PIPELINE METADATA -->
  <div class="modal-backdrop" id="modal-pipeline">
    <div class="modal-box">
      <div class="modal-head">Register Pipeline Metadata</div>
      <form onsubmit="submitPipeline(event)">
        <div class="field-box" style="margin-bottom:20px;">
          <label>Program</label>
          <select id="m-pipe-prog">
            <option value="basic_forwarding.p4">basic_forwarding.p4</option>
            <option value="firewall_drop.p4">firewall_drop.p4</option>
          </select>
        </div>
        <div style="display:flex; justify-content:flex-end; gap:8px;">
          <button type="button" class="btn btn-ghost btn-sm" onclick="closeModal('modal-pipeline')">Cancel</button>
          <button type="submit" class="btn btn-sm">Register</button>
        </div>
      </form>
    </div>
  </div>

  <!-- Minimal Footer -->
  <footer class="site-footer">
    <a href="https://www.openlabnetworks.com/" target="_blank" rel="noopener noreferrer">OPENLABNETWORKS</a>
    <div>DPU &bull; SMARTNIC &bull; HARDWARE OFFLOAD</div>
  </footer>

  <script>
    const API = window.location.origin;

    function scrollToStudio() {
      document.getElementById('workbench').scrollIntoView({ behavior: 'smooth' });
    }

    function selectTab(id) {
      document.querySelectorAll('.tab-btn').forEach(function(b) { b.classList.remove('active'); });
      document.querySelectorAll('.tab-pane').forEach(function(p) { p.classList.remove('active'); });

      const btn = document.getElementById('tab-btn-' + id);
      const pane = document.getElementById('pane-' + id);
      if (btn) btn.classList.add('active');
      if (pane) pane.classList.add('active');
    }

    function openModal(id) {
      if (id === 'modal-port') document.getElementById('m-mac').value = randomMac();
      document.getElementById(id).style.display = 'flex';
    }
    function closeModal(id) {
      document.getElementById(id).style.display = 'none';
    }

    function randomMac() {
      const h = '0123456789abcdef';
      let m = '00:1a:2b';
      for (let i = 0; i < 3; i++) {
        m += ':' + h[Math.floor(Math.random()*16)] + h[Math.floor(Math.random()*16)];
      }
      return m;
    }

    async function fetchPorts() {
      try {
        const res = await fetch(API + '/v1/ports');
        const data = await res.json();
        const tbody = document.getElementById('ports-tbody');
        tbody.innerHTML = '';
        const list = data.ports || [];
        if (list.length === 0) {
          const tr = document.createElement('tr');
          const td = document.createElement('td');
          td.colSpan = 7;
          td.style.cssText = 'text-align:center; padding:18px;';
          td.textContent = 'No ports configured.';
          tr.appendChild(td);
          tbody.appendChild(tr);
          return;
        }
        list.forEach(function(p) {
          const tr = document.createElement('tr');

          const tdId = document.createElement('td');
          const strong = document.createElement('strong');
          strong.textContent = p.id;
          tdId.appendChild(strong);
          tr.appendChild(tdId);

          const tdName = document.createElement('td');
          tdName.textContent = p.name || '';
          tr.appendChild(tdName);

          const tdMac = document.createElement('td');
          const code = document.createElement('code');
          code.textContent = p.mac_address;
          tdMac.appendChild(code);
          tr.appendChild(tdMac);

          const tdMtu = document.createElement('td');
          tdMtu.textContent = p.mtu;
          tr.appendChild(tdMtu);

          const tdSpeed = document.createElement('td');
          tdSpeed.textContent = p.speed || '100G';
          tr.appendChild(tdSpeed);

          const tdState = document.createElement('td');
          const stateSpan = document.createElement('span');
          stateSpan.className = 'tag-state up';
          stateSpan.textContent = p.oper_state;
          tdState.appendChild(stateSpan);
          tr.appendChild(tdState);

          const tdAction = document.createElement('td');
          const delBtn = document.createElement('button');
          delBtn.className = 'btn btn-ghost btn-sm';
          delBtn.textContent = 'Delete';
          delBtn.addEventListener('click', function() { delPort(p.id); });
          tdAction.appendChild(delBtn);
          tr.appendChild(tdAction);

          tbody.appendChild(tr);
        });
      } catch (e) {
        console.error(e);
      }
    }

    async function submitPort(e) {
      e.preventDefault();
      const mac = document.getElementById('m-mac').value.trim();
      const mtu = parseInt(document.getElementById('m-mtu').value, 10);
      try {
        const res = await fetch(API + '/v1/ports', {
          method: 'POST',
          headers: {'Content-Type': 'application/json'},
          body: JSON.stringify({ mac_address: mac, mtu: mtu })
        });
        if (res.ok) {
          closeModal('modal-port');
          fetchPorts();
        }
      } catch (err) {
        alert(err.message);
      }
    }

    async function delPort(id) {
      if (!confirm('Delete port ' + id + '?')) return;
      await fetch(API + '/v1/ports/' + id, { method: 'DELETE' });
      fetchPorts();
    }

    async function fetchPipelines() {
      try {
        const res = await fetch(API + '/v1/pipelines');
        const data = await res.json();
        const tbody = document.getElementById('pipelines-tbody');
        tbody.innerHTML = '';
        const list = data.pipelines || [];
        if (list.length === 0) {
          const tr = document.createElement('tr');
          const td = document.createElement('td');
          td.colSpan = 4;
          td.style.cssText = 'text-align:center; padding:18px;';
          td.textContent = 'No pipelines deployed.';
          tr.appendChild(td);
          tbody.appendChild(tr);
          return;
        }
        list.forEach(function(pipe) {
          const tr = document.createElement('tr');

          const tdId = document.createElement('td');
          const strong = document.createElement('strong');
          strong.textContent = pipe.id;
          tdId.appendChild(strong);
          tr.appendChild(tdId);

          const tdProg = document.createElement('td');
          const code = document.createElement('code');
          code.textContent = pipe.p4_program;
          tdProg.appendChild(code);
          tr.appendChild(tdProg);

          const tdStatus = document.createElement('td');
          const statusSpan = document.createElement('span');
          statusSpan.className = 'tag-state up';
          statusSpan.textContent = pipe.status || 'ACTIVE';
          tdStatus.appendChild(statusSpan);
          tr.appendChild(tdStatus);

          const tdAction = document.createElement('td');
          const delBtn = document.createElement('button');
          delBtn.className = 'btn btn-ghost btn-sm';
          delBtn.textContent = 'Delete';
          delBtn.addEventListener('click', function() { delPipe(pipe.id); });
          tdAction.appendChild(delBtn);
          tr.appendChild(tdAction);

          tbody.appendChild(tr);
        });
      } catch (e) {
        console.error(e);
      }
    }

    async function submitPipeline(e) {
      e.preventDefault();
      const prog = document.getElementById('m-pipe-prog').value;
      try {
        const res = await fetch(API + '/v1/pipelines', {
          method: 'POST',
          headers: {'Content-Type': 'application/json'},
          body: JSON.stringify({ p4_program: prog })
        });
        if (res.ok) {
          closeModal('modal-pipeline');
          fetchPipelines();
        }
      } catch (err) {
        alert(err.message);
      }
    }

    async function delPipe(id) {
      if (!confirm('Delete pipeline ' + id + '?')) return;
      await fetch(API + '/v1/pipelines/' + id, { method: 'DELETE' });
      fetchPipelines();
    }

    async function runSimulation() {
      const prog = document.getElementById('sim-program').value;
      const ethertype = document.getElementById('sim-ethertype').value;
      const banner = document.getElementById('sim-banner');
      banner.style.display = 'none';

      for (let i = 1; i <= 5; i++) {
        const el = document.getElementById('stage-' + i);
        el.classList.add('lit');
        await new Promise(function(r) { setTimeout(r, 160); });
        el.classList.remove('lit');
      }

      banner.style.display = 'block';
      if (prog === 'firewall_drop.p4') {
        if (ethertype !== '0x0800') {
          banner.className = 'sim-banner drop';
          banner.innerHTML = '<strong>[DROPPED]</strong> EtherType ' + ethertype + ' != 0x0800. P4 marked packet to drop in Ingress stage.';
        } else {
          banner.className = 'sim-banner pass';
          banner.innerHTML = '<strong>[FORWARDED]</strong> Valid IPv4 packet. Ingress forwarded to Port 1.';
        }
      } else {
        banner.className = 'sim-banner pass';
        banner.innerHTML = '<strong>[FORWARDED]</strong> IPv4 LPM table match. TTL decremented and emitted by Deparser.';
      }
    }

    async function callApi(method, path, body) {
      document.getElementById('curl-hint').innerText = 'curl -X ' + method + ' ' + API + path;
      const opts = { method: method };
      if (body) {
        opts.headers = { 'Content-Type': 'application/json' };
        opts.body = JSON.stringify(body);
      }
      try {
        const res = await fetch(API + path, opts);
        const json = await res.json();
        const box = document.getElementById('term-box');
        box.innerText = '// ' + method + ' ' + path + ' &bull; Status ' + res.status + '\n' + JSON.stringify(json, null, 2) + '\n\n' + box.innerText;
      } catch (err) {
        document.getElementById('term-box').innerText = err.message;
      }
    }

    async function runSmokeTest() {
      selectTab('console');
      scrollToStudio();
      await callApi('GET', '/healthz');
      await callApi('GET', '/v1/ports');
      fetchPorts();
    }

    window.addEventListener('DOMContentLoaded', function() {
      fetchPorts();
      fetchPipelines();
    });
  </script>
</body>
</html>
`

// ServeAbout renders the project overview page separately from the interactive
// workbench so visitors can understand the platform before using it.
func ServeAbout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(aboutHTML))
}

const aboutHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>About | OPENLABNETWORKS</title>
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
      --mono: "DM Mono", monospace;
      --sans: "Manrope", sans-serif;
    }
    * { box-sizing: border-box; }
    body { margin: 0; color: var(--ink); background: var(--paper); font-family: var(--sans); }
    a { color: inherit; }
    .about-nav { display: flex; justify-content: space-between; align-items: center; padding: 24px 7vw; border-bottom: 1px solid var(--line); }
    .logo { text-decoration: none; font-weight: 800; letter-spacing: -.04em; }
    .logo span { color: var(--accent); }
    .back { font: 12px var(--mono); text-transform: uppercase; text-decoration: none; color: var(--muted); }
    .about-main { max-width: 1180px; margin: 0 auto; padding: 84px 7vw 110px; }
    .eyebrow { color: var(--teal); font: 12px var(--mono); letter-spacing: .12em; text-transform: uppercase; }
    h1 { max-width: 800px; margin: 18px 0 24px; font-size: clamp(48px, 8vw, 104px); line-height: .96; letter-spacing: -.07em; }
    .intro { max-width: 680px; color: var(--muted); font-size: 20px; line-height: 1.65; }
    .section { margin-top: 100px; }
    h2 { margin: 0 0 14px; font-size: 34px; letter-spacing: -.05em; }
    .section-lead { max-width: 650px; color: var(--muted); line-height: 1.7; }
    .grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 16px; margin-top: 28px; }
    .card { padding: 28px; min-height: 220px; background: var(--card); border: 1px solid var(--line); }
    .number { color: var(--accent); font: 13px var(--mono); }
    h3 { margin: 38px 0 10px; font-size: 21px; letter-spacing: -.03em; }
    .card p, .architecture p { color: var(--muted); line-height: 1.65; }
    .architecture { display: grid; grid-template-columns: .85fr 1.15fr; gap: 48px; align-items: start; }
    .stack { border-left: 2px solid var(--accent); }
    .layer { padding: 19px 22px; border-bottom: 1px solid var(--line); background: var(--card); }
    .layer strong { display: block; margin-bottom: 5px; }
    .layer span { color: var(--muted); font: 12px var(--mono); }
    .about-footer { padding: 26px 7vw; border-top: 1px solid var(--line); color: var(--muted); font: 11px var(--mono); text-transform: uppercase; letter-spacing: .08em; }
    @media (max-width: 760px) {
      .about-nav { padding: 20px; }
      .about-main { padding: 60px 20px 80px; }
      .grid, .architecture { grid-template-columns: 1fr; }
      .section { margin-top: 70px; }
    }
  </style>
</head>
<body>
  <header class="about-nav">
    <a class="logo" href="/ui">OPENLABNETWORKS</a>
    <a class="back" href="/ui">&larr; Open sandbox</a>
  </header>
  <main class="about-main">
    <div class="eyebrow">OPENLABNETWORKS / Core use cases</div>
    <h1>Move infrastructure work closer to the wire.</h1>
    <p class="intro">OPENLABNETWORKS Open Sandbox is a practical environment for designing, validating, and demonstrating hardware-assisted networking workflows before dedicated DPU or SmartNIC hardware is available.</p>

    <section class="section">
      <h2>Built for the work between idea and hardware.</h2>
      <p class="section-lead">The sandbox gives platform teams a repeatable place to exercise control APIs, packet-processing behavior, and deployment workflows. It shortens feedback loops without pretending that software emulation replaces production silicon.</p>
      <div class="grid">
        <article class="card"><div class="number">01 / DEVELOP</div><h3>Prototype offload services</h3><p>Model virtual ports, service boundaries, and packet-handling decisions while the hardware design is still evolving.</p></article>
        <article class="card"><div class="number">02 / VALIDATE</div><h3>Test control-plane behavior</h3><p>Use a clean, resettable environment to validate client integrations, lifecycle operations, and health checks.</p></article>
        <article class="card"><div class="number">03 / DEMONSTRATE</div><h3>Share a working path</h3><p>Give application and infrastructure teams a common, runnable reference for hardware-assisted networking concepts.</p></article>
      </div>
    </section>

    <section class="section architecture">
      <div>
        <div class="eyebrow">System architecture</div>
        <h2>Small surface. Clear boundaries.</h2>
        <p>The environment separates the API control plane from the packet-processing target. That makes experiments observable and keeps the transition to a physical accelerator straightforward.</p>
      </div>
      <div class="stack">
        <div class="layer"><strong>Developer experience</strong><span>Browser workbench · curl · smoke tests</span></div>
        <div class="layer"><strong>Control plane</strong><span>Go service · REST resources · in-memory sandbox state</span></div>
        <div class="layer"><strong>Acceleration target</strong><span>Software dataplane · virtual interfaces · programmable processing</span></div>
        <div class="layer"><strong>Deployment layer</strong><span>Docker Compose · Linux · optional cloud image</span></div>
      </div>
    </section>
  </main>
  <footer class="about-footer">OPENLABNETWORKS &bull; DPU and SmartNIC enablement</footer>
</body>
</html>
`
