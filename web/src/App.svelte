<script lang="ts">
import './style/main.css';

import { onMount } from 'svelte';
import { Navbar, NavBrand, NavLi, NavUl, NavHamburger } from 'flowbite-svelte';

import { NotAuthorized } from '$lib/request';
import { user } from '$lib/stores';
import { logout } from '$lib/auth';
import { currentRoute } from './routes';
import { hasAccess, Permission } from './lib/auth';

import Tty from '$components/Tty.svelte'
import ToastRack from '$components/toast/ToastRack.svelte';

let ttyOpen = $state(false)

const toggleTty = () => {
    ttyOpen = !ttyOpen
}

onMount(() => {
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
            {#if hasAccess($user, Permission.Admin)}
                <NavLi href="/system/logs">Logs</NavLi>
                <NavLi onclick={ toggleTty }>Terminal</NavLi>
            {/if}
            <NavLi onclick={logout}>Logout</NavLi>
        {:else}
            <NavLi href="/account/login">Login</NavLi>
        {/if}
    </NavUl>
</Navbar>

<ToastRack />
<Tty open={ ttyOpen }/>
{@render $currentRoute()}
