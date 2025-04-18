<script lang="ts">
import { onMount } from 'svelte'
import { login, verifySession } from "$lib/auth";

import { Card, Button, Label, Input } from 'flowbite-svelte';

let username = $state("")
let password = $state("")
let loaded = $state(false)

const handleLogin = async (e: Event) => {
    e.preventDefault()

    await login(username, password)
}

onMount(async () => {
    try {
        await verifySession()
    } catch {}

    loaded = true
})
</script>

{#if loaded}
    <main class="flex justify-center">
        <div class="content-center">
        <Card>
            <form class="flex flex-col space-y-6" onsubmit={handleLogin}>
                <h3 class="text-xl font-medium text-gray-900 dark:text-white">Login</h3>
                <Label class="space-y-2">
                    <span>Username</span>
                    <Input type="text" placeholder="Username" bind:value={username} required />
                </Label>
                <Label class="space-y-2">
                    <span>Password</span>
                    <Input type="password" name="password" placeholder="•••••" bind:value={password} required />
                </Label>
                <Button type="submit" class="w-full">Login</Button>
            </form>
        </Card>
        </div>
    </main>
{/if}

