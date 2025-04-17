<script lang="ts">
import { onMount, onDestroy } from "svelte";
import { jwt } from "$lib/stores";
let monitorData = $state("{}")

let stream: EventSource;
onMount(() => {
    stream = new EventSource(`http://localhost:8080/api/system/monitor?bearer=${$jwt}`)

    stream.onerror = e => monitorData = `{"error": "${e}"}`
    stream.addEventListener("message", ({ data }) => {
        console.log(data)
        monitorData = data
    })
})

onDestroy(() => {
    stream?.close()
})
</script>

<pre id="monitor">
    { JSON.stringify(JSON.parse(monitorData), null, 4) }
</pre>
