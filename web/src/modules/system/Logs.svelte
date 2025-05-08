<script lang="ts">
import { onMount, onDestroy } from "svelte";

import { jwt } from "$lib/stores";
import toast from "$lib/toast";


let logs = $state([])


let stream: EventSource;
onMount(() => {
    stream = new EventSource(`http://localhost:8080/api/system/logs?bearer=${$jwt}`)

    stream.onerror = e => toast.error(JSON.stringify(e))
    stream.addEventListener("message", ({ data }) => {
        console.log(data)
        logs.push(data)
    })
})

onDestroy(() => {
    stream?.close()
})
</script>

<section id="system-logs">
    {#each logs as log}
        <article class="log-entry">{ log }</article>
    {/each}
</section>
