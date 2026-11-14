<script lang="ts">
    import "../app.css";

    let { children } = $props();
</script>

<svelte:head>
    <title>Solar Tracker — CE224</title>
    <meta name="description" content="Hệ thống bám sát năng lượng mặt trời 2 trục" />
</svelte:head>

{@render children()}
