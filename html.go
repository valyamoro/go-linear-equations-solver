package main

const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width,initial-scale=1">
  <title>Linear Equations System</title>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link href="https://fonts.googleapis.com/css2?family=Space+Mono:wght@400;700&family=Oxanium:wght@300;400;600;800&display=swap" rel="stylesheet">
  <script src="https://cdn.plot.ly/plotly-2.32.0.min.js" charset="utf-8"></script>
  <style>
    :root {
      --bg:        #06101e;
      --surface:   #0a1628;
      --surface2:  #0e1f38;
      --border:    rgba(100,160,255,0.1);
      --border2:   rgba(100,160,255,0.2);
      --text:      #c5d8f0;
      --muted:     #4a6a9a;
      --muted2:    #6a8ab0;
      --accent:    #00e6a0;
      --accent2:   #00b87a;
      --blue:      #4d9eff;
      --red:       #ff4d6d;
      --yellow:    #ffd166;
    }

    *, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }

    body {
      background: var(--bg);
      color: var(--text);
      font-family: 'Space Mono', 'Courier New', monospace;
      min-height: 100vh;
      display: flex;
      flex-direction: column;
      overflow: hidden;
    }

    body::before {
      content: '';
      position: fixed; inset: 0;
      background: repeating-linear-gradient(
        0deg,
        transparent,
        transparent 2px,
        rgba(0,0,0,0.04) 2px,
        rgba(0,0,0,0.04) 4px
      );
      pointer-events: none;
      z-index: 1000;
    }

    header {
      background: var(--surface);
      border-bottom: 1px solid var(--border2);
      padding: 0 24px;
      height: 56px;
      display: flex;
      align-items: center;
      gap: 16px;
      flex-shrink: 0;
      position: relative;
      overflow: hidden;
    }
    header::after {
      content: '';
      position: absolute;
      bottom: 0; left: 0; right: 0;
      height: 1px;
      background: linear-gradient(90deg, transparent, var(--accent), transparent);
      opacity: 0.6;
    }

    .logo-mark {
      width: 34px; height: 34px;
      border: 1.5px solid var(--accent);
      border-radius: 6px;
      display: flex; align-items: center; justify-content: center;
      font-family: 'Oxanium', sans-serif;
      font-weight: 800;
      font-size: 16px;
      color: var(--accent);
      flex-shrink: 0;
      box-shadow: 0 0 12px rgba(0,230,160,0.2);
      animation: pulse-border 3s ease infinite;
    }
    @keyframes pulse-border {
      0%, 100% { box-shadow: 0 0 8px rgba(0,230,160,0.2); }
      50%       { box-shadow: 0 0 20px rgba(0,230,160,0.4); }
    }

    .header-title {
      font-family: 'Oxanium', sans-serif;
      font-weight: 600;
      font-size: 15px;
      letter-spacing: .05em;
      color: var(--text);
    }
    .header-title span {
      color: var(--accent);
    }

    .header-tag {
      margin-left: auto;
      font-size: 10px;
      color: var(--muted);
      letter-spacing: .08em;
      border: 1px solid var(--border);
      padding: 3px 8px;
      border-radius: 3px;
      background: rgba(0,0,0,0.2);
    }

    .workspace {
      flex: 1;
      display: flex;
      min-height: 0;
    }

    aside {
      width: 288px;
      flex-shrink: 0;
      background: var(--surface);
      border-right: 1px solid var(--border);
      display: flex;
      flex-direction: column;
      overflow-y: auto;
      overflow-x: hidden;
    }

    .sidebar-section {
      padding: 18px 16px;
      border-bottom: 1px solid var(--border);
    }
    .sidebar-section:last-child {
      border-bottom: none;
    }

    .section-label {
      font-size: 9px;
      font-weight: 700;
      letter-spacing: .15em;
      text-transform: uppercase;
      color: var(--muted);
      margin-bottom: 12px;
      display: flex;
      align-items: center;
      gap: 8px;
    }
    .section-label::after {
      content: '';
      flex: 1;
      height: 1px;
      background: var(--border);
    }

    .eq-list { display: flex; flex-direction: column; gap: 7px; }
    .eq-row {
      display: flex;
      align-items: center;
      gap: 10px;
      background: var(--surface2);
      border: 1px solid var(--border);
      border-radius: 5px;
      padding: 8px 11px;
      font-size: 12px;
      transition: border-color .2s, background .2s;
      animation: slide-in .35s ease backwards;
    }
    .eq-row:hover {
      background: rgba(0,230,160,0.05);
      border-color: rgba(0,230,160,0.25);
    }
    .color-pip {
      width: 8px; height: 8px;
      border-radius: 50%;
      flex-shrink: 0;
      box-shadow: 0 0 6px currentColor;
    }
    .eq-label {
      color: var(--text);
      font-size: 12px;
    }

    .status-box {
      display: flex;
      align-items: center;
      gap: 10px;
      padding: 11px 14px;
      border-radius: 6px;
      font-size: 12px;
      font-family: 'Oxanium', sans-serif;
      font-weight: 600;
      letter-spacing: .03em;
    }
    .status-box.ok {
      background: rgba(0,230,160,0.08);
      border: 1px solid rgba(0,230,160,0.3);
      color: var(--accent);
    }
    .status-box.fail {
      background: rgba(255,77,109,0.08);
      border: 1px solid rgba(255,77,109,0.3);
      color: var(--red);
    }
    .status-icon { font-size: 16px; }

    .sol-list { display: flex; flex-direction: column; gap: 5px; }
    .sol-row {
      background: var(--surface2);
      border: 1px solid var(--border);
      border-radius: 4px;
      padding: 6px 10px;
      font-size: 11px;
      color: var(--accent);
      display: flex;
      align-items: center;
      gap: 8px;
    }
    .sol-row::before {
      content: '◆';
      font-size: 8px;
      color: var(--accent2);
    }

    .range-row {
      display: flex;
      flex-direction: column;
      gap: 5px;
    }
    .range-item {
      background: var(--surface2);
      border: 1px solid var(--border);
      border-radius: 4px;
      padding: 6px 10px;
      font-size: 11px;
      color: var(--muted2);
    }
    .range-item strong { color: var(--blue); }

    .chart-wrap {
      flex: 1;
      position: relative;
      min-width: 0;
      min-height: 0;
    }
    #plot {
      width: 100%;
      height: 100%;
    }

    @keyframes slide-in {
      from { opacity: 0; transform: translateX(-8px); }
      to   { opacity: 1; transform: translateX(0); }
    }
    .eq-row:nth-child(1) { animation-delay: 0.05s; }
    .eq-row:nth-child(2) { animation-delay: 0.10s; }
    .eq-row:nth-child(3) { animation-delay: 0.15s; }
    .eq-row:nth-child(4) { animation-delay: 0.20s; }
    .eq-row:nth-child(5) { animation-delay: 0.25s; }
    .eq-row:nth-child(6) { animation-delay: 0.30s; }
    .eq-row:nth-child(7) { animation-delay: 0.35s; }
    .eq-row:nth-child(8) { animation-delay: 0.40s; }

    ::-webkit-scrollbar { width: 5px; }
    ::-webkit-scrollbar-track { background: transparent; }
    ::-webkit-scrollbar-thumb {
      background: rgba(100,160,255,0.2);
      border-radius: 3px;
    }
    ::-webkit-scrollbar-thumb:hover {
      background: rgba(100,160,255,0.35);
    }
  </style>
</head>
<body>

<!-- HEADER -->
<header>
  <div class="logo-mark">=</div>
  <div class="header-title">
    Linear <span>Equations System</span>
  </div>
  <div class="header-tag">Plotly.js · 2D · Go</div>
</header>

<!-- WORKSPACE -->
<div class="workspace">

  <!-- SIDEBAR -->
  <aside>

    <!-- Equation list -->
    <div class="sidebar-section">
      <div class="section-label">System</div>
      <div class="eq-list">
        {{range .Equations}}
        <div class="eq-row">
          <div class="color-pip" style="background:{{.Color}};color:{{.Color}}"></div>
          <span class="eq-label">{{.Label}}</span>
        </div>
        {{end}}
      </div>
    </div>

    <!-- Status -->
    <div class="sidebar-section">
      <div class="section-label">Result</div>
      {{if .IsEmpty}}
      <div class="status-box fail">
        <span class="status-icon">✗</span>
        System incompatible
      </div>
      {{else}}
      <div class="status-box ok">
        <span class="status-icon">✓</span>
        Solutions found
      </div>
      {{end}}
    </div>

    <!-- Solutions -->
    {{if not .IsEmpty}}
    <div class="sidebar-section">
      <div class="section-label">Solutions ({{len .Solutions}})</div>
      <div class="sol-list">
        {{range .Solutions}}
        <div class="sol-row">{{fmtPt .X .Y}}</div>
        {{end}}
      </div>
    </div>
    {{end}}

    <!-- Viewport -->
    <div class="sidebar-section">
      <div class="section-label">View range</div>
      <div class="range-row">
        <div class="range-item"><strong>x</strong> ∈ [{{fmtF .XMin}}, {{fmtF .XMax}}]</div>
        <div class="range-item"><strong>y</strong> ∈ [{{fmtF .YMin}}, {{fmtF .YMax}}]</div>
      </div>
    </div>

  </aside>

  <!-- CHART -->
  <div class="chart-wrap">
    <div id="plot"></div>
  </div>

</div><!-- .workspace -->

<script>
(function () {
  var traces = {{.TracesJSON}};
  var layout = {{.LayoutJSON}};
  var config = {
    responsive: true,
    displaylogo: false,
    modeBarButtonsToRemove: ['lasso2d', 'select2d', 'autoScale2d'],
    toImageButtonOptions: { format: 'png', filename: 'linear-equations', scale: 2 }
  };
  Plotly.newPlot('plot', traces, layout, config);
})();
</script>

</body>
</html>`
