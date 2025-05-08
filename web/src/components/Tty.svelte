<script lang="ts">
import "@xterm/xterm/css/xterm.css"

import { onDestroy, onMount } from "svelte";
import type { Action } from "svelte/action";
import { Terminal } from '@xterm/xterm'
import { FitAddon } from '@xterm/addon-fit'

import { jwt } from "$lib/stores";
import toast from "$lib/toast";


let initialized = $state(false)
let { open } = $props()


let xterm: Terminal
let fit: FitAddon
let sock: WebSocket
let termElement: HTMLElement


const termElementMounted: Action = (node: HTMLElement) => {
    termElement = node
}

const resizeTerminal = () => {
    console.log("resize:" + JSON.stringify({
        cols: xterm.cols,
        rows: xterm.rows,
    }))
    fit.fit()
    sock.send("resize:" + JSON.stringify({
        cols: xterm.cols,
        rows: xterm.rows,
    }))
}


$effect(() => {
    console.log(open, initialized)
    if (!open || initialized) {
        return
    }

    initialized = true
    xterm = new Terminal()
    fit = new FitAddon()

    xterm.loadAddon(fit)
    xterm.open(termElement)
    fit.fit()

    sock = new WebSocket(`ws://localhost:8080/api/system/tty?bearer=${$jwt}`)
    sock.onmessage = e => {
        const str = new String(e.data)
        const splitI = str.indexOf(":")
        const type = str.substring(0, splitI)
        const msg = str.substring(splitI + 1)

        switch (type) {
        case "io":
            xterm.write(msg.replace('\\r\\n', '\r\n'))
            break
        case "notice":
            toast.info(msg)
            break
        }
    }

    xterm.onData(data => {
        sock.send(`io:${data}`)
    })

    window.addEventListener('resize', resizeTerminal)
    sock.addEventListener('open', resizeTerminal)
})

onDestroy(() => {
    window.removeEventListener('resize', resizeTerminal)
    sock.close()
})
</script>


<svelte:head>
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/xterm/css/xterm.css" />
    <script src="https://cdn.jsdelivr.net/npm/xterm/lib/xterm.js"></script>
    <script src="https://cdn.jsdelivr.net/npm/xterm-addon-fit/lib/xterm-addon-fit.js"></script>
</svelte:head>


<article id="terminal"
    use:termElementMounted
    class={ open && initialized ? "" : "hidden" }
></article>

