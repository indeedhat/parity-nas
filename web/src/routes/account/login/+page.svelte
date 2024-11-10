<script lang="ts">
import { login } from "$lib/auth";
    import { hydrate } from "svelte";

let username = $state("")
let password = $state("")

const handleLogin = async () => await login(username, password)
</script>

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
