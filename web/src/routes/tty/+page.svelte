<script lang="ts">
import "@xterm/xterm/css/xterm.css"

import { onDestroy, onMount } from "svelte";
import type { Action } from "svelte/action";
import { Terminal } from '@xterm/xterm'
import { FitAddon } from '@xterm/addon-fit'

import { jwt } from "$lib/stores";


let xterm: Terminal
let fit: FitAddon
let sock: WebSocket
let termElement: HTMLElement


const termElementMounted: Action = (node: HTMLElement) => {
    termElement = node
}

const resizeTerminal = () => {
    fit.fit()
    sock.send("resize:" + JSON.stringify({
        cols: xterm.cols,
        rows: xterm.rows,
    }))
}


onMount(() => {
    xterm = new Terminal()
    fit = new FitAddon()

    xterm.loadAddon(fit)
    xterm.open(termElement)
    fit.fit()

    sock = new WebSocket(`ws://localhost:8080/api/system/tty?bearer=${jwt}`)
    sock.onmessage = e => {
        console.log({ onmessage: e.data })

        let [ type, msg ] = new String(e.data).split(":", 2)
        switch (type) {
        case "io":
            xterm.write(msg)
            break
        case "notice":
            console.log({ msg })
            break
        }
    }

    xterm.onData(data => {
        console.log({ onData: data })
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


<article id="terminal" use:termElementMounted></article>
