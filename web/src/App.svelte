<script lang="ts">
import './style/main.css';

import { onMount } from 'svelte';
import { NotAuthorized } from '$lib/request';
import { user } from '$lib/stores';
import { logout, redirectGuests } from '$lib/auth';
import { currentRoute } from './routes';
import ToastRack from '$components/toast/ToastRack.svelte';
import { Navbar, NavBrand, NavLi, NavUl, NavHamburger } from 'flowbite-svelte';

let { children } = $props();

onMount(() => {
    redirectGuests($user)

    window.onunhandledrejection = (e) => {
        e.stopPropagation();

        if (e instanceof NotAuthorized) {
            window.location.href = "/account/login";
        }
    };
});
</script>

<Navbar>
    <NavBrand href="/home">
        <span class="self-center whitespace-nowrap text-xl font-semibold dark:text-white">Parity NAS</span>
    </NavBrand>
    <NavHamburger  />
    <NavUl>
        {#if $user?.name}
            <NavLi href="/home">Home</NavLi>
            <NavLi href="/tty">Terminal</NavLi>
            <NavLi onclick={logout}>Logout</NavLi>
        {:else}
            <NavLi href="/account/login">Login</NavLi>
        {/if}
    </NavUl>
</Navbar>

<ToastRack />
{@render $currentRoute()}
