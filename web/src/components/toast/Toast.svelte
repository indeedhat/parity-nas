<script lang="ts">
// This comomnet is used to display a toast notification
// it makes use of the sveltestrap Alert component to do so

import { fade } from 'svelte/transition'

import { Alert } from 'flowbite-svelte';
import { InfoCircleSolid } from 'flowbite-svelte-icons';

let { type, dismissible, dismiss, children } = $props()
let color: string = $state("primary");

switch (type) {
case 'alert':
    color = "green";
    break;
case 'error':
    color = "red";
    break;
case 'info':
    color = "blue";
    break;
}
const handleDismiss = () => {
    if (dismissible) {
        dismiss()
    }
}
</script>

<article transition:fade onclick={ handleDismiss } class="m-2">
    <Alert border color={ color } class="shadow">
        <InfoCircleSolid slot="icon" class="w-5 h-5" />
        {@render children()}
    </Alert>
</article>

