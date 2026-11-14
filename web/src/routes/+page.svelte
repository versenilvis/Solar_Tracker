<script lang="ts">
    import { ws } from "$lib/ws.svelte";
    import SensorChart from "$lib/components/SensorChart.svelte";
    import AngleDisplay from "$lib/components/AngleDisplay.svelte";
    import LdrBar from "$lib/components/LdrBar.svelte";

    const volt  = $derived(ws.latest?.volt ?? null);
    const curr  = $derived(ws.latest?.curr ?? null);
    const power = $derived(ws.latest?.power ?? null);
    const energy = $derived(ws.latest?.energy ?? null);
    const azi   = $derived(ws.latest?.azi ?? 90);
    const ele   = $derived(ws.latest?.ele ?? 90);
    const ldr   = $derived(ws.latest?.ldr ?? [0, 0, 0, 0] as [number, number, number, number]);
    const mode  = $derived(ws.latest?.mode ?? "---");

    const ts = $derived(
        ws.latest
            ? new Date(ws.latest.timestamp).toLocaleTimeString("vi-VN", { hour12: false })
            : "--:--:--",
    );

    const todayDate = new Date().toLocaleDateString("vi-VN", {
        weekday: "long",
        year: "numeric",
        month: "long",
        day: "numeric",
    });
</script>

<main>
    <header>
        <div class="brand">
            <svg width="22" height="22" viewBox="0 0 24 24" fill="none">
                <circle cx="12" cy="12" r="4" fill="#f59e0b" />
                <g stroke="#f59e0b" stroke-width="2" stroke-linecap="round">
                    <line x1="12" y1="2" x2="12" y2="5" />
                    <line x1="12" y1="19" x2="12" y2="22" />
                    <line x1="2" y1="12" x2="5" y2="12" />
                    <line x1="19" y1="12" x2="22" y2="12" />
                    <line x1="4.22" y1="4.22" x2="6.34" y2="6.34" />
                    <line x1="17.66" y1="17.66" x2="19.78" y2="19.78" />
                    <line x1="19.78" y1="4.22" x2="17.66" y2="6.34" />
                    <line x1="6.34" y1="17.66" x2="4.22" y2="19.78" />
                </g>
            </svg>
            <span>Solar Tracker</span>
        </div>

        <div class="status-group">
            <div class="date-label">{todayDate}</div>
            <div class="mode-badge" class:auto={mode === "auto"}>
                {mode.toUpperCase()}
            </div>
            <div class="status" class:online={ws.connected}>
                <span class="dot"></span>
                {ws.connected ? "Live" : "Connecting..."}
            </div>
        </div>
    </header>

    <!-- stat cards -->
    <section class="stat-row">
        <div class="stat-card power">
            <div class="stat-label">
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none">
                    <path d="M13 2L4.5 13.5H11L10 22l9.5-12H14L13 2z" stroke="#f59e0b" stroke-width="2" stroke-linejoin="round"/>
                </svg>
                Power
            </div>
            <div class="stat-value power-val">
                {power !== null ? power.toFixed(1) : "---"}
                <span class="unit">mW</span>
            </div>
            <div class="stat-sub">Peak: {ws.stats.max_power.toFixed(1)} mW</div>
        </div>

        <div class="stat-card volt">
            <div class="stat-label">
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none">
                    <path d="M12 2v20M2 12h20" stroke="#38bdf8" stroke-width="2" stroke-linecap="round"/>
                </svg>
                Voltage
            </div>
            <div class="stat-value volt-val">
                {volt !== null ? volt.toFixed(2) : "---"}
                <span class="unit">V</span>
            </div>
            <div class="stat-sub">{curr !== null ? curr.toFixed(1) : "---"} mA</div>
        </div>

        <div class="stat-card energy">
            <div class="stat-label">
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none">
                    <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z" stroke="#4ade80" stroke-width="2"/>
                </svg>
                Energy
            </div>
            <div class="stat-value energy-val">
                {energy !== null ? energy.toFixed(2) : "---"}
                <span class="unit">Wh</span>
            </div>
            <div class="stat-sub">tích lũy</div>
        </div>

        <div class="stat-card meta">
            <div class="stat-label">Last updated</div>
            <div class="stat-ts">{ts}</div>
            <div class="stat-label" style="margin-top: 10px">Data points</div>
            <div class="stat-ts">{ws.history.length}<span class="unit"> / 60</span></div>
        </div>
    </section>

    <!-- chart -->
    <section class="chart-section">
        <div class="chart-card">
            <div class="chart-header">
                <span class="chart-title power-text">Power</span>
                <span class="chart-badge power-badge">mW</span>
            </div>
            <div class="chart-wrap">
                <SensorChart
                    history={ws.history}
                    field="power"
                    color="#f59e0b"
                    glow="rgba(245,158,11,0.08)"
                    unit="mW" />
            </div>
        </div>

        <div class="chart-card">
            <div class="chart-header">
                <span class="chart-title volt-text">Voltage / Current</span>
                <span class="chart-badge volt-badge">V</span>
            </div>
            <div class="chart-wrap">
                <SensorChart
                    history={ws.history}
                    field="volt"
                    color="#38bdf8"
                    glow="rgba(56,189,248,0.08)"
                    unit="V" />
            </div>
        </div>
    </section>

    <!-- angle + ldr -->
    <section class="bottom-row">
        <div class="stat-card bottom-card">
            <div class="stat-label" style="margin-bottom: 8px">Servo Angles</div>
            <AngleDisplay {azi} {ele} />
        </div>

        <div class="stat-card bottom-card">
            <div class="stat-label" style="margin-bottom: 8px">LDR Sensors</div>
            <LdrBar {ldr} />
        </div>
    </section>
</main>

<style>
    main {
        max-width: 1200px;
        margin: 0 auto;
        padding: 24px 20px 48px;
        display: flex;
        flex-direction: column;
        gap: 24px;
    }

    /* header */
    header {
        display: flex;
        align-items: center;
        justify-content: space-between;
    }

    .brand {
        display: flex;
        align-items: center;
        gap: 10px;
        font-size: 1.1rem;
        font-weight: 600;
        letter-spacing: -0.02em;
        color: var(--text-primary);
    }

    .status-group {
        display: flex;
        align-items: center;
        gap: 12px;
    }

    .date-label {
        font-size: 0.9rem;
        font-weight: 500;
        color: var(--text-secondary);
    }

    .mode-badge {
        font-size: 0.72rem;
        font-weight: 700;
        letter-spacing: 0.08em;
        font-family: var(--font-mono);
        padding: 4px 10px;
        border-radius: 999px;
        background: rgba(100, 100, 120, 0.15);
        border: 1px solid rgba(100, 100, 120, 0.25);
        color: var(--text-secondary);
    }

    .mode-badge.auto {
        background: rgba(245, 158, 11, 0.12);
        border-color: rgba(245, 158, 11, 0.3);
        color: var(--solar-primary);
    }

    .status {
        display: flex;
        align-items: center;
        gap: 7px;
        font-size: 0.8rem;
        font-weight: 500;
        color: var(--offline);
        background: rgba(239, 68, 68, 0.1);
        border: 1px solid rgba(239, 68, 68, 0.2);
        padding: 5px 12px;
        border-radius: 999px;
        transition: all 0.3s;
    }

    .status.online {
        color: var(--online);
        background: rgba(34, 197, 94, 0.1);
        border-color: rgba(34, 197, 94, 0.2);
    }

    .dot {
        width: 7px;
        height: 7px;
        border-radius: 50%;
        background: currentColor;
    }

    .status.online .dot {
        animation: pulse 2s infinite;
    }

    @keyframes pulse {
        0%, 100% { opacity: 1; transform: scale(1); }
        50%       { opacity: 0.5; transform: scale(0.8); }
    }

    /* stat cards */
    .stat-row {
        display: grid;
        grid-template-columns: 1fr 1fr 1fr 200px;
        gap: 16px;
    }

    .stat-card {
        background: var(--bg-surface);
        border: 1px solid var(--border);
        border-radius: var(--radius-md);
        padding: 20px 22px;
        display: flex;
        flex-direction: column;
        gap: 8px;
        transition: border-color 0.2s, box-shadow 0.2s;
    }

    .stat-card.power:hover { border-color: rgba(245,158,11,0.4); box-shadow: 0 0 20px rgba(245,158,11,0.07); }
    .stat-card.volt:hover  { border-color: rgba(56,189,248,0.4); box-shadow: 0 0 20px rgba(56,189,248,0.07); }
    .stat-card.energy:hover { border-color: rgba(74,222,128,0.4); box-shadow: 0 0 20px rgba(74,222,128,0.07); }

    .stat-label {
        display: flex;
        align-items: center;
        gap: 6px;
        font-size: 0.75rem;
        font-weight: 500;
        color: var(--text-secondary);
        text-transform: uppercase;
        letter-spacing: 0.06em;
    }

    .stat-value {
        font-size: 2.4rem;
        font-weight: 700;
        font-family: var(--font-mono);
        letter-spacing: -0.03em;
        line-height: 1;
    }

    .power-val  { color: var(--solar-primary); }
    .volt-val   { color: var(--blue-primary); }
    .energy-val { color: var(--green-primary); }

    .unit {
        font-size: 0.42em;
        font-weight: 400;
        color: var(--text-secondary);
        vertical-align: super;
        margin-left: 2px;
    }

    .stat-sub {
        font-size: 0.8rem;
        color: var(--text-secondary);
        font-family: var(--font-mono);
    }

    .stat-ts {
        font-family: var(--font-mono);
        font-size: 1.25rem;
        font-weight: 500;
        color: var(--text-primary);
    }

    /* charts */
    .chart-section {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 16px;
    }

    .chart-card {
        background: var(--bg-surface);
        border: 1px solid var(--border);
        border-radius: var(--radius-md);
        padding: 20px 22px;
        display: flex;
        flex-direction: column;
        gap: 16px;
    }

    .chart-header {
        display: flex;
        align-items: center;
        justify-content: space-between;
    }

    .chart-title { font-size: 0.85rem; font-weight: 600; }
    .power-text  { color: var(--solar-primary); }
    .volt-text   { color: var(--blue-primary); }

    .chart-badge {
        font-family: var(--font-mono);
        font-size: 0.72rem;
        padding: 2px 8px;
        border-radius: 4px;
        font-weight: 500;
    }

    .power-badge { background: rgba(245,158,11,0.12); color: var(--solar-primary); }
    .volt-badge  { background: rgba(56,189,248,0.12); color: var(--blue-primary); }

    .chart-wrap {
        height: 220px;
        position: relative;
    }

    /* bottom row */
    .bottom-row {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 16px;
    }

    .bottom-card {
        min-height: 180px;
    }

    /* responsive */
    @media (max-width: 1000px) {
        .stat-row { grid-template-columns: 1fr 1fr; }
        .stat-card.meta { grid-column: 1 / -1; flex-direction: row; gap: 24px; }
        .chart-section { grid-template-columns: 1fr; }
        .bottom-row { grid-template-columns: 1fr; }
    }

    @media (max-width: 540px) {
        .stat-row { grid-template-columns: 1fr; }
        .stat-value { font-size: 2rem; }
        .date-label { display: none; }
    }
</style>
