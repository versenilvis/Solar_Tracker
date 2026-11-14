<script lang="ts">
    type Props = {
        azi: number;
        ele: number;
    };

    let { azi, ele }: Props = $props();
</script>

<div class="angle-wrap">
    <div class="axis">
        <span class="axis-label">Azimuth</span>
        <div class="angle-value">
            {azi}<span class="unit">°</span>
        </div>
        <div class="arc-bar">
            <div class="arc-fill" style="width: {(azi / 180) * 100}%"></div>
        </div>
        <div class="axis-range">
            <span>0°</span><span>180°</span>
        </div>
    </div>

    <div class="axis">
        <span class="axis-label">Elevation</span>
        <div class="angle-value">
            {ele}<span class="unit">°</span>
        </div>
        <div class="arc-bar">
            <div class="arc-fill ele" style="width: {(ele / 180) * 100}%"></div>
        </div>
        <div class="axis-range">
            <span>0°</span><span>180°</span>
        </div>
    </div>
</div>

<style>
    .angle-wrap {
        display: flex;
        flex-direction: column;
        gap: 20px;
        height: 100%;
        justify-content: center;
    }

    .axis {
        display: flex;
        flex-direction: column;
        gap: 6px;
    }

    .axis-label {
        font-size: 0.72rem;
        font-weight: 600;
        letter-spacing: 0.06em;
        text-transform: uppercase;
        color: var(--text-secondary);
    }

    .angle-value {
        font-family: var(--font-mono);
        font-size: 2rem;
        font-weight: 700;
        letter-spacing: -0.03em;
        color: var(--solar-primary);
        line-height: 1;
    }

    .unit {
        font-size: 0.5em;
        font-weight: 400;
        color: var(--text-secondary);
        vertical-align: super;
        margin-left: 2px;
    }

    .arc-bar {
        height: 4px;
        border-radius: 2px;
        background: #252a3a;
        overflow: hidden;
    }

    .arc-fill {
        height: 100%;
        border-radius: 2px;
        background: var(--solar-primary);
        transition: width 0.3s ease;
    }

    .arc-fill.ele {
        background: var(--blue-primary);
    }

    .axis-range {
        display: flex;
        justify-content: space-between;
        font-size: 0.65rem;
        color: var(--text-dim);
        font-family: var(--font-mono);
    }
</style>
