package handlers

import "net/http"

// ServeSecurityPage renders the hardware-security offload demonstration.
func ServeSecurityPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(securityHTML))
}

const securityHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Hardware Security | ProtonNIC Sandbox</title>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;700&family=Space+Grotesk:wght@500;700;900&family=Inter:wght@400;500;600&display=swap" rel="stylesheet">
  <style>
    :root{--bg:#f3f4f6;--surface:#fff;--ink:#111216;--muted:#69707d;--line:#dfe2e7;--line-dark:#111216;--cobalt:#2563eb;--green:#10b981;--red:#ef4444;--amber:#f59e0b;--display:"Space Grotesk",sans-serif;--body:"Inter",sans-serif;--mono:"JetBrains Mono",monospace}
    *{box-sizing:border-box;margin:0;padding:0}
    body{min-height:100vh;color:var(--ink);background:var(--bg);font-family:var(--body);background-image:linear-gradient(rgba(17,18,22,.035) 1px,transparent 1px),linear-gradient(90deg,rgba(17,18,22,.035) 1px,transparent 1px);background-size:44px 44px;-webkit-font-smoothing:antialiased}
    a{color:inherit}
    .site-header{height:76px;padding:0 48px;display:grid;grid-template-columns:minmax(220px,1fr) auto minmax(220px,1fr);align-items:center;border-bottom:1px solid var(--line-dark);background:rgba(243,244,246,.96)}
    .brand{font:900 16px var(--display);letter-spacing:-.04em;text-decoration:none;text-transform:uppercase}
    .site-nav{display:flex;justify-content:center;gap:34px;font:11px var(--mono);letter-spacing:.09em;text-transform:uppercase}
    .site-nav a{color:var(--muted);text-decoration:none}
    .site-nav a.active,.site-nav a:hover{color:var(--ink)}
    .header-action{justify-self:end}
    .btn{min-height:42px;padding:10px 15px;border:1px solid var(--line-dark);background:var(--ink);color:#fff;display:inline-flex;align-items:center;justify-content:center;gap:12px;font:700 10px var(--mono);letter-spacing:.08em;text-transform:uppercase;text-decoration:none;cursor:pointer}
    .btn:hover{background:#292b31}.btn:disabled{opacity:.55;cursor:wait}.btn-ghost{background:var(--surface);color:var(--ink)}.btn-ghost:hover{background:#fff}
    .circle{width:20px;height:20px;border:1px solid rgba(255,255,255,.5);border-radius:50%;display:inline-grid;place-items:center}
    .meta{height:40px;padding:0 48px;border-bottom:1px solid var(--line);display:grid;grid-template-columns:1fr auto 1fr;align-items:center;color:var(--muted);font:9px var(--mono);letter-spacing:.1em;text-transform:uppercase}
    .meta span:last-child{text-align:right}.online{display:inline-flex;align-items:center;gap:8px}.online i{width:7px;height:7px;border-radius:50%;background:var(--green);box-shadow:0 0 0 3px rgba(16,185,129,.12)}
    main{overflow:hidden}.lab-shell{max-width:1440px;margin:auto;padding:46px 48px 76px}
    .hero-grid{display:grid;grid-template-columns:minmax(0,1.15fr) minmax(360px,.85fr);border:1px solid var(--line-dark);background:var(--surface)}
    .hero-copy{min-height:430px;padding:44px 46px 34px;display:flex;flex-direction:column;justify-content:space-between}
    .eyebrow,.section-label{color:var(--muted);font:9px var(--mono);letter-spacing:.13em;text-transform:uppercase}.eyebrow{display:flex;justify-content:space-between;margin-bottom:26px}.eyebrow b{color:var(--cobalt)}
    h1{max-width:760px;font:900 clamp(58px,7vw,98px)/.86 var(--display);letter-spacing:-.075em;text-transform:uppercase}
    h1 span{display:block;color:var(--cobalt)}.hero-lead{max-width:650px;margin:24px 0 30px;display:flex;gap:12px;color:var(--muted);font:11px/1.75 var(--mono)}.hero-lead i{width:8px;height:8px;margin-top:6px;flex:none;border-radius:50%;background:var(--green)}
    .hero-actions{display:flex;gap:10px}.hero-status{margin-top:34px;padding-top:17px;border-top:1px solid var(--line);display:flex;justify-content:space-between;color:var(--muted);font:9px var(--mono);letter-spacing:.08em;text-transform:uppercase}
    .topology{position:relative;min-height:430px;padding:28px;background:#14161b;color:#fff;overflow:hidden}
    .topology:before{content:"";position:absolute;inset:0;background-image:linear-gradient(rgba(255,255,255,.055) 1px,transparent 1px),linear-gradient(90deg,rgba(255,255,255,.055) 1px,transparent 1px);background-size:28px 28px}
    .topology-head{position:relative;display:flex;justify-content:space-between;color:#8e95a3;font:9px var(--mono);letter-spacing:.1em;text-transform:uppercase}
    .topology-stage{position:absolute;z-index:1;left:28px;right:28px;top:86px;display:grid;grid-template-columns:1fr 30px 1.15fr 30px 1fr;align-items:center}
    .node{min-height:112px;border:1px solid #3a3d46;background:#1a1d23;padding:16px}.node.primary{border-color:#5185ef;background:rgba(37,99,235,.16)}.node small{display:block;color:#9299a7;font:8px var(--mono);letter-spacing:.1em;text-transform:uppercase}.node strong{display:block;margin:13px 0 7px;font:700 16px var(--display);text-transform:uppercase}.node span{color:#aeb4bf;font:9px var(--mono)}.wire{color:#60a5fa;text-align:center;font:18px var(--mono)}
    .telemetry{position:absolute;z-index:2;left:28px;right:28px;bottom:28px;display:grid;grid-template-columns:repeat(3,1fr);border:1px solid #343740}
    .metric{padding:15px 16px;border-right:1px solid #343740}.metric:last-child{border:0}.metric span{display:block;color:#858c99;font:8px var(--mono);letter-spacing:.1em;text-transform:uppercase}.metric b{display:block;margin-top:7px;color:#fff;font:700 20px var(--mono)}.metric b.blue{color:#60a5fa}
    .section{margin-top:48px}.section-head{display:flex;align-items:end;justify-content:space-between;gap:24px;margin-bottom:18px}.section-head h2{margin-top:8px;font:900 36px/.95 var(--display);letter-spacing:-.055em;text-transform:uppercase}.section-head p{max-width:430px;color:var(--muted);font-size:13px;line-height:1.65}
    .lab-grid{display:grid;grid-template-columns:repeat(4,1fr);border-top:1px solid var(--line-dark);border-left:1px solid var(--line-dark)}
    .lab-card{min-height:210px;padding:22px;background:rgba(255,255,255,.72);border-right:1px solid var(--line-dark);border-bottom:1px solid var(--line-dark)}.lab-card:hover{background:#fff}.lab-card .num{color:var(--cobalt);font:9px var(--mono);letter-spacing:.1em}.lab-card h3{margin:30px 0 10px;font:700 20px var(--display);text-transform:uppercase}.lab-card p{color:var(--muted);font-size:12px;line-height:1.65}
    .simulator{border:1px solid var(--line-dark);background:var(--surface)}.sim-head{min-height:74px;padding:18px 20px;border-bottom:1px solid var(--line-dark);display:flex;align-items:center;justify-content:space-between}.sim-head h2{font:900 22px var(--display);letter-spacing:-.04em;text-transform:uppercase}.status-pill{padding:7px 10px;border:1px solid var(--line);color:var(--muted);font:8px var(--mono);letter-spacing:.08em;text-transform:uppercase}.status-pill b{color:var(--green)}
    .controls{display:grid;grid-template-columns:1.1fr 1.1fr .8fr .8fr 1.1fr .9fr 140px;border-bottom:1px solid var(--line-dark)}.field{min-width:0;padding:14px;border-right:1px solid var(--line)}.field:last-of-type{border:0}.field label{display:block;margin-bottom:7px;color:var(--muted);font:8px var(--mono);letter-spacing:.1em;text-transform:uppercase}.field input,.field select{width:100%;height:38px;padding:0 9px;border:1px solid var(--line);border-radius:0;background:#fff;color:var(--ink);font:10px var(--mono);outline:none}.field input:focus,.field select:focus{border-color:var(--cobalt)}.run-field{display:flex;align-items:end;padding:14px}.run-field .btn{width:100%;height:38px}
    .packet-path{display:grid;grid-template-columns:repeat(6,1fr);gap:8px;padding:18px}.packet-step{min-height:105px;padding:13px 12px;border:1px solid var(--line);background:#fafafa;transition:.24s ease}.packet-step.active{background:var(--cobalt);border-color:var(--cobalt);color:#fff;transform:translateY(-3px)}.packet-step.drop{background:var(--red);border-color:var(--red);color:#fff}.packet-step.pass{background:var(--green);border-color:var(--green);color:#fff}.packet-step small{display:block;color:var(--cobalt);font:8px var(--mono);letter-spacing:.08em}.packet-step.active small,.packet-step.drop small,.packet-step.pass small{color:#fff}.packet-step b{display:block;margin-top:14px;font:700 13px var(--display);text-transform:uppercase}.packet-step span{display:block;margin-top:7px;color:var(--muted);font:8px/1.35 var(--mono)}.packet-step.active span,.packet-step.drop span,.packet-step.pass span{color:#fff}
    .result{display:none;padding:18px 20px;background:#14161b;color:#fff;border-top:1px solid var(--line-dark)}.result.show{display:grid;grid-template-columns:220px 1fr;gap:20px}.result-title{font:900 19px var(--display);text-transform:uppercase}.result-title.allow{color:#6ee7b7}.result-title.deny{color:#fca5a5}.result-copy{color:#cbd0d8;font:9px/1.65 var(--mono)}.result-copy code{display:block;margin-top:8px;padding:9px;background:#20232a;color:#93c5fd;overflow:auto}
    .rest-layout{display:grid;grid-template-columns:.72fr 1.28fr;border:1px solid var(--line-dark);background:var(--surface)}.rest-info{padding:24px;border-right:1px solid var(--line-dark)}.rest-info h3{margin:8px 0 12px;font:700 23px var(--display);letter-spacing:-.04em;text-transform:uppercase}.rest-info p{color:var(--muted);font-size:12px;line-height:1.65}.endpoint{margin-top:18px;padding:12px;border:1px solid var(--line);background:#f8f8f9;font:9px var(--mono);overflow-wrap:anywhere}.terminal{min-width:0;background:#101216;color:#dfe3e8}.terminal-bar{height:40px;padding:0 14px;border-bottom:1px solid #30343b;display:flex;align-items:center;justify-content:space-between;color:#8d95a3;font:8px var(--mono);letter-spacing:.08em;text-transform:uppercase}.copy-btn{border:0;background:transparent;color:#93c5fd;font:8px var(--mono);cursor:pointer;text-transform:uppercase}.terminal pre{padding:18px 20px 20px;overflow:auto;font:10px/1.75 var(--mono);white-space:pre-wrap}.terminal .prompt{color:#60a5fa}.terminal .json{color:#a7f3d0}
    .notice{margin-top:16px;padding:12px 14px;border-left:3px solid var(--cobalt);color:var(--muted);background:rgba(37,99,235,.04);font:9px/1.65 var(--mono)}
    .site-footer{min-height:72px;padding:0 48px;border-top:1px solid var(--line-dark);display:flex;align-items:center;justify-content:space-between;color:var(--muted);font:9px var(--mono);letter-spacing:.09em;text-transform:uppercase}
    @media(max-width:1100px){.hero-grid{grid-template-columns:1fr}.topology{min-height:390px}.lab-grid{grid-template-columns:repeat(2,1fr)}.controls{grid-template-columns:repeat(3,1fr)}.field{border-bottom:1px solid var(--line)}.run-field{grid-column:span 1}.packet-path{grid-template-columns:repeat(3,1fr)}.rest-layout{grid-template-columns:1fr}.rest-info{border-right:0;border-bottom:1px solid var(--line-dark)}}
    @media(max-width:760px){.site-header{height:auto;min-height:68px;padding:14px 20px;grid-template-columns:minmax(0,1fr) auto}.brand{font-size:14px}.header-action{display:none}.site-nav{grid-column:1/-1;grid-row:2;width:100%;margin-top:12px;padding-top:12px;border-top:1px solid var(--line);justify-content:flex-start;gap:18px;overflow-x:auto}.meta{height:auto;min-height:38px;padding:9px 20px;grid-template-columns:1fr auto}.meta span:nth-child(2){display:none}.lab-shell{padding:28px 20px 48px}.hero-copy{min-height:390px;padding:34px 20px 24px}.hero-grid h1{font-size:58px}.hero-status{flex-direction:column;gap:7px}.topology{min-height:360px;padding:20px}.topology-stage{left:20px;right:20px;grid-template-columns:1fr 18px 1.15fr 18px 1fr}.node{min-height:100px;padding:11px}.node strong{font-size:13px}.telemetry{left:20px;right:20px;bottom:20px}.metric{padding:12px 8px}.metric b{font-size:14px}.section{margin-top:38px}.section-head{align-items:flex-start;flex-direction:column}.section-head h2{font-size:31px}.lab-grid{grid-template-columns:1fr}.lab-card{min-height:180px}.controls{grid-template-columns:1fr}.field{border-right:0}.run-field{grid-column:auto}.packet-path{grid-template-columns:1fr 1fr}.result.show{grid-template-columns:1fr}.site-footer{min-height:88px;padding:18px 20px;align-items:flex-start;flex-direction:column;justify-content:center;gap:8px}}
  </style>
</head>
<body>
  <header class="site-header">
    <a class="brand" href="/ui">PROTONNIC SANDBOX</a>
    <nav class="site-nav"><a href="/ui">Workbench</a><a href="/about">About</a><a class="active" href="/security">Hardware Security</a></nav>
    <a class="btn header-action" href="#simulator">Run Simulation <span class="circle">&rarr;</span></a>
  </header>
  <div class="meta"><span>Hardware Security / SmartNIC Lab</span><span class="online"><i></i> API Online</span><span>OPI + gRPC + P4</span></div>
  <main>
    <div class="lab-shell">
      <section class="hero-grid">
        <div class="hero-copy">
          <div>
            <div class="eyebrow"><span>Cloud-native security offload</span><b>Hardware boundary 01</b></div>
            <h1>SECURE<span>THE WIRE.</span></h1>
            <div class="hero-lead"><i></i><span>Move Kubernetes policy enforcement from the host CPU to a programmable SmartNIC/DPU datapath. Test the control plane and packet decision here without real accelerator hardware.</span></div>
            <div class="hero-actions"><a class="btn" href="#simulator">Open Simulator <span class="circle">&rarr;</span></a><a class="btn btn-ghost" href="#architecture">View Architecture &darr;</a></div>
          </div>
          <div class="hero-status"><span>Mode: deterministic REST model</span><span>Physical DPU: not required</span></div>
        </div>
        <div class="topology">
          <div class="topology-head"><span>Data plane topology</span><span>Live model</span></div>
          <div class="topology-stage">
            <div class="node"><small>Policy source</small><strong>Kubernetes</strong><span>NetworkPolicy</span></div>
            <div class="wire">&rarr;</div>
            <div class="node primary"><small>Control plane</small><strong>OPI / gRPC</strong><span>Operator compiler</span></div>
            <div class="wire">&rarr;</div>
            <div class="node"><small>Datapath</small><strong>P4 ACL</strong><span>Permit / drop</span></div>
          </div>
          <div class="telemetry"><div class="metric"><span>Enforcement</span><b class="blue">Line rate</b></div><div class="metric"><span>Host CPU</span><b>0*</b></div><div class="metric"><span>Boundary</span><b>NIC</b></div></div>
        </div>
      </section>

      <section class="section" id="architecture">
        <div class="section-head"><div><div class="section-label">01 / Offload surfaces</div><h2>Security outside the host path.</h2></div><p>Select hardware-compatible policy and filtering functions, while keeping the actual cost and capability of cryptography and L7 services target-specific.</p></div>
        <div class="lab-grid">
          <article class="lab-card"><div class="num">01 / L3-L4</div><h3>Micro-segmentation</h3><p>Compile Kubernetes workload identity and NetworkPolicy rules into P4 match-action ACL entries.</p></article>
          <article class="lab-card"><div class="num">02 / CRYPTO</div><h3>Encryption boundary</h3><p>Model IPsec or transport crypto when production hardware provides isolated key storage and crypto engines.</p></article>
          <article class="lab-card"><div class="num">03 / L7</div><h3>Service policy</h3><p>Represent JWT, rate, and API controls in a DPU ARM control-plane workflow; behavior remains target-dependent.</p></article>
          <article class="lab-card"><div class="num">04 / SIGNAL</div><h3>Telemetry and DDoS</h3><p>Filter floods at ingress and stream flow context for monitoring instead of relying on host packet collection.</p></article>
        </div>
      </section>

      <section class="section" id="simulator">
        <div class="section-head"><div><div class="section-label">02 / Interactive policy lab</div><h2>NetworkPolicy attack simulator.</h2></div><p>Choose workload namespaces, protocol, port, policy action, and a threat scenario. The REST endpoint returns the modeled P4 result.</p></div>
        <div class="simulator">
          <div class="sim-head"><h2>Policy-to-datapath console</h2><div class="status-pill">Endpoint <b>simulation</b></div></div>
          <div class="controls">
            <div class="field"><label for="source">Source namespace</label><input id="source" value="payments" maxlength="63"></div>
            <div class="field"><label for="destination">Destination namespace</label><input id="destination" value="inventory" maxlength="63"></div>
            <div class="field"><label for="protocol">Protocol</label><select id="protocol"><option>TCP</option><option>UDP</option></select></div>
            <div class="field"><label for="port">Port</label><input id="port" type="number" min="1" max="65535" value="443"></div>
            <div class="field"><label for="action">Policy action</label><select id="action"><option value="allow">Allow</option><option value="deny">Deny</option></select></div>
            <div class="field"><label for="attack">Threat scenario</label><select id="attack"><option value="normal">Normal traffic</option><option value="syn_flood">SYN flood</option><option value="port_scan">Port scan</option><option value="dns_amplification">DNS amplification</option><option value="lateral_movement">Lateral movement</option><option value="malformed_packet">Malformed packet</option><option value="ddos_http">HTTP DDoS</option></select></div>
            <div class="run-field"><button class="btn" id="run-button" onclick="runOffloadSimulation()">Run <span class="circle">&rarr;</span></button></div>
          </div>
          <div class="packet-path">
            <div class="packet-step" id="flow-1"><small>01 / SOURCE</small><b>NetworkPolicy</b><span>Workload identity</span></div>
            <div class="packet-step" id="flow-2"><small>02 / COMPILE</small><b>Operator</b><span>Resolve policy intent</span></div>
            <div class="packet-step" id="flow-3"><small>03 / CONTROL</small><b>OPI / gRPC</b><span>Write ACL intent</span></div>
            <div class="packet-step" id="flow-4"><small>04 / DATAPATH</small><b>P4 pipeline</b><span>Line-rate match</span></div>
            <div class="packet-step" id="flow-5"><small>05 / NIC</small><b>Packet verdict</b><span>Permit or drop</span></div>
            <div class="packet-step" id="flow-6"><small>06 / HOST</small><b>Workload</b><span>Allowed traffic only</span></div>
          </div>
          <div class="result" id="result"><div class="result-title" id="status"></div><div class="result-copy"><span id="details"></span><code id="entry"></code></div></div>
        </div>
      </section>

      <section class="section">
        <div class="section-head"><div><div class="section-label">03 / Terminal REST interface</div><h2>Run it through the API.</h2></div><p>Start the sandbox, then call the same deterministic endpoint used by this page. Use <code>jq</code> to format the JSON response.</p></div>
        <div class="rest-layout">
          <div class="rest-info"><div class="section-label">POST / REST / JSON</div><h3>One policy request.</h3><p>The request is validated by the Go API and converted into a modeled P4 ACL table entry. It does not change BMv2 or configure host networking.</p><div class="endpoint">POST /api/v1/security/offload/simulations</div><div class="notice">This is a safe software model. A physical DPU, P4Runtime target, Kubernetes cluster, and crypto engine are not configured.</div></div>
          <div class="terminal"><div class="terminal-bar"><span>Terminal / bash</span><button class="copy-btn" onclick="copyCurl()">Copy curl</button></div><pre id="curl-example"><span class="prompt">$</span> curl -sS -X POST http://localhost:8080/api/v1/security/offload/simulations \
  -H 'Content-Type: application/json' \
  -d '{
    "source_namespace": "internet",
    "destination_namespace": "payments",
    "protocol": "TCP",
    "destination_port": 443,
    "policy_action": "allow",
    "attack_type": "syn_flood"
  }' | jq

<span class="json">{ "mode": "SIMULATION", "decision": "DENY", "threat_class": "Volumetric SYN flood" }</span></pre></div>
        </div>
      </section>
    </div>
  </main>
  <footer class="site-footer"><span>PROTONNIC SANDBOX &bull; DPU &amp; SMARTNIC ENABLEMENT</span><span>SOFTWARE MODEL &bull; HARDWARE READY</span></footer>
  <script>
    function pause(ms){return new Promise(function(resolve){setTimeout(resolve,ms)})}
    function resetFlow(){for(var i=1;i<=6;i++){document.getElementById('flow-'+i).className='packet-step'}}
    async function runOffloadSimulation(){
      var button=document.getElementById('run-button'),result=document.getElementById('result');
      resetFlow();result.className='result';button.disabled=true;
      var payload={source_namespace:document.getElementById('source').value.trim(),destination_namespace:document.getElementById('destination').value.trim(),protocol:document.getElementById('protocol').value,destination_port:Number(document.getElementById('port').value),policy_action:document.getElementById('action').value,attack_type:document.getElementById('attack').value};
      try{
        var response=await fetch('/api/v1/security/offload/simulations',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(payload)}),data=await response.json();
        if(!response.ok)throw new Error(data.error||'Simulation request failed');
        for(var i=1;i<=5;i++){var step=document.getElementById('flow-'+i);step.classList.add('active');await pause(260);step.classList.remove('active')}
        if(data.decision==='DENY')document.getElementById('flow-5').classList.add('drop');
        else{document.getElementById('flow-5').classList.add('pass');document.getElementById('flow-6').classList.add('pass')}
        var status=document.getElementById('status');status.textContent=data.decision==='DENY'?'DROP AT SMARTNIC':'ALLOW TO HOST';status.className='result-title '+(data.decision==='DENY'?'deny':'allow');
        document.getElementById('details').textContent=(data.threat_class?data.threat_class+'. ':'')+data.details+' '+data.host_cpu_impact+' '+data.disclaimer;
        document.getElementById('entry').textContent=data.p4_table_entry;
      }catch(err){var status=document.getElementById('status');status.textContent='INVALID REQUEST';status.className='result-title deny';document.getElementById('details').textContent=err.message;document.getElementById('entry').textContent=''}
      finally{result.className='result show';button.disabled=false}
    }
    function copyCurl(){
      var text="curl -sS -X POST http://localhost:8080/api/v1/security/offload/simulations \\\n  -H 'Content-Type: application/json' \\\n  -d '{\"source_namespace\":\"internet\",\"destination_namespace\":\"payments\",\"protocol\":\"TCP\",\"destination_port\":443,\"policy_action\":\"allow\",\"attack_type\":\"syn_flood\"}' | jq";
      if(navigator.clipboard){navigator.clipboard.writeText(text)}
      var button=document.querySelector('.copy-btn');button.textContent='Copied';
      setTimeout(function(){button.textContent='Copy curl'},1500);
    }
  </script>
</body>
</html>
`
