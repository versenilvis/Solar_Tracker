<script lang="ts">
    type Props = {
        ldr: [number, number, number, number];
    };

    let { ldr }: Props = $props();

    const MAX_ADC = 4095;
    const labels = ["TL", "TR", "BL", "BR"];
    const colors = ["#f59e0b", "#38bdf8", "#4ade80", "#a78bfa"];
</script>

<div class="ldr-wrap">
    {#each ldr as val, i}
        <div class="ldr-row">
            <span class="ldr-label" style="color: {colors[i]}">{labels[i]}</span>
            <div class="bar-track">
                <div
                    class="bar-fill"
                    style="width: {(val / MAX_ADC) * 100}%; background: {colors[i]}"
                ></div>
            </div>
            <span class="ldr-val">{val}</span>
        </div>
    {/each}
</div>

<style>
    .ldr-wrap {
        display: flex;
        flex-direction: column;
        gap: 10px;
        justify-content: center;
        height: 100%;
    }

    .ldr-row {
        display: grid;
        grid-template-columns: 24px 1fr 48px;
        align-items: center;
        gap: 10px;
    }

    .ldr-label {
        font-size: 0.7rem;
        font-weight: 700;
        letter-spacing: 0.05em;
        font-family: var(--font-mono);
    }

    .bar-track {
        height: 6px;
        border-radius: 3px;
        background: #252a3a;
        overflow: hidden;
    }

    .bar-fill {
        height: 100%;
        border-radius: 3px;
        transition: width 0.3s ease;
        opacity: 0.85;
    }

    .ldr-val {
        font-family: var(--font-mono);
        font-size: 0.72rem;
        color: var(--text-secondary);
        text-align: right;
    }
</style>
