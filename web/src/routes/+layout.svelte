<script lang="ts">
import '$style/main.css'

import { onMount } from 'svelte';

import { NotAuthorized } from '$lib/request';
import { user } from '$lib/stores'
import { logout } from '$lib/auth'

import ToastRack from '$components/toast/ToastRack.svelte';


let { children } = $props();

onMount(() => {
    window.onunhandledrejection = e => {
        e.stopPropagation();

        if (e instanceof NotAuthorized) {
            window.location.href = "/account/login";
        }
    };
})
</script>

<section id="body">
    <header>
        {#if $user?.name}
            <a href="/tty">Terminal</a>
            <a onclick={logout} href="#">Logout</a>
        {:else}
            <a href="/account/login">Login</a>
        {/if}
    </header>
    <ToastRack />
    {@render children()}
</section>
