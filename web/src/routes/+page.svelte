<script lang="ts">
import { onMount } from 'svelte'
import { login, verifySession } from "$lib/auth";

let username = $state("")
let password = $state("")
let loaded = $state(false)

const handleLogin = async () => await login(username, password)

onMount(async () => {
    try {
        await verifySession()
    } catch {}

    loaded = true
})
</script>

{#if loaded}
    <section>
        <fieldset>
            <legend>Login</legend>

            <form onsubmit={handleLogin}>
                <div>
                    <label for="username">Username</label>
                    <input type="text" id="username" bind:value={username} />
                </div>
                <div>
                    <label for="password">Password</label>
                    <input type="password" id="password" bind:value={password} />
                </div>
                <div>
                    <input type="submit" value="Login" />
                </div>
            </form>
        </fieldset>
    </section>
{/if}
