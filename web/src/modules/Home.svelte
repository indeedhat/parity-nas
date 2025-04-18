<script lang="ts">
import { onMount, onDestroy } from "svelte";

import { jwt } from "$lib/stores";

import { Alert, Card, Progressbar } from 'flowbite-svelte'
import { InfoCircleSolid, LabelOutline } from 'flowbite-svelte-icons';

const percent = (total: number, free: number): number => {
    return (100 - (100 / total * free)).toFixed(1)
}

let monitorData = $state({})
let cpuTotal = $derived.by(() => {
    if (!monitorData.cpu) {
        return 0
    }

    let total = 0;
    let idle = 0;

    for (let i in monitorData.cpu) {
        total += monitorData.cpu[i].total
        idle += monitorData.cpu[i].idle
    }

    return percent(total, idle)
})

let memoryTotal = $derived.by(() => {
    if (!monitorData.memory) {
        return 0
    }

    return percent(monitorData.memory.total, monitorData.memory.total - monitorData.memory.used)
})

const formatBytes = (bytes: number, decimals: number = 2): string => {
    if (!+bytes) return '0 Bytes'

    const k = 1024
    const dm = decimals < 0 ? 0 : decimals
    const sizes = ['Bytes', 'KiB', 'MiB', 'GiB', 'TiB', 'PiB']

    const i = Math.floor(Math.log(bytes) / Math.log(k))

    return `${parseFloat((bytes / Math.pow(k, i)).toFixed(dm))} ${sizes[i]}`
}

const formatBits = (bits: number, decimals: number = 2): string => {
    if (!+bits) return '0'

    const k = 1024
    const dm = decimals < 0 ? 0 : decimals
    const sizes = ['b/s', 'kb/s', 'mb/s', 'gb/s', 'tb/s', 'pb/s']

    const i = Math.floor(Math.log(bits) / Math.log(k))

    return `${parseFloat((bits / Math.pow(k, i)).toFixed(dm))} ${sizes[i]}`
}

let stream: EventSource;
onMount(() => {
    stream = new EventSource(`http://localhost:8080/api/system/monitor?bearer=${$jwt}`)

    stream.onerror = e => monitorData = { error: e }
    stream.addEventListener("message", ({ data }) => {
        monitorData = JSON.parse(data)
    })
})

onDestroy(() => {
    stream?.close()
})
</script>

<section id="sysmon">
    {#if monitorData.error}
        <Alert border color="red" class="shadow">
            <InfoCircleSolid slot="icon" class="w-5 h-5" />
            {@render children()}
        </Alert>
    {:else if Object.keys(monitorData).length}
        <Card>
            <h3 class="text-xl font-medium text-gray-900 dark:text-white">Uptime</h3>
            {monitorData.uptime}
        </Card>
        <Card>
            <h3 class="text-xl font-medium text-gray-900 dark:text-white">Memory {memoryTotal}%</h3>
            <div class="flex justify-end mb-1">
                <span class="text-sm font-medium text-blue-700 dark:text-white">
                    {formatBytes(monitorData.memory.used)} /
                    {formatBytes(monitorData.memory.total)}
                </span>
            </div>
            <Progressbar progress={memoryTotal} labelOutside={``} />
        </Card>
        <Card>
            <h3 class="text-xl font-medium text-gray-900 dark:text-white">CPU {cpuTotal}%</h3>
            {#each Object.entries(monitorData.cpu) as [key, core] (key)}
                <Progressbar progress={percent(core.total, core.idle)} labelOutside={key} />
            {/each}
        </Card>
        {#if monitorData.network}
            <Card>
                <h3 class="text-xl font-medium text-gray-900 dark:text-white">Network</h3>
                <!-- TODO: fix values -->
                {#each Object.entries(monitorData.network) as [key, iface] (key)}
                    <div class="text-sm font-medium text-blue-700 dark:text-white">{key}</div>
                    <div class="text-sm font-medium text-blue-700 dark:text-white">TX: {formatBits(iface.tx)}</div>
                    <div class="text-sm font-medium text-blue-700 dark:text-white">RX: {formatBits(iface.rx)}</div>
                {/each}
            </Card>
        {/if}
        {#if monitorData.network}
            <Card>
                <h3 class="text-xl font-medium text-gray-900 dark:text-white">Mounts</h3>
                {#each Object.entries(monitorData.mounts) as [key, disk] (key)}
                    <Progressbar progress={percent(disk.total, disk.total - disk.used)} labelOutside={key} />
                    <div class="flex justify-end mb-1">
                        <span class="text-sm font-medium text-blue-700 dark:text-white">
                            {formatBytes(disk.used)} /
                            {formatBytes(disk.total)}
                        </span>
                    </div>
                {/each}
            </Card>
        {/if}
    {/if}
</section>
