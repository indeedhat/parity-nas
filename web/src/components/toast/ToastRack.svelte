<script lang="ts">
    // This component provides a container elemest for toast nootifications to be displayed in
    import Toast from '$components/toast/Toast.svelte';

    import { dismissToast, toasts } from '$lib/toast';
</script>

{#if $toasts}
    <section class="toastc">
        {#each $toasts as toast (toast.id)}
            <Toast
                type={ toast.type }
                dismissible={ toast.dismissible }
                dismiss={ () => dismissToast(toast.id) }
            >
                { toast.message }
            </Toast>
        {/each}
    </section>
{/if}

<style lang="postcss">
    section {
        position: fixed;
        top:0;
        right: 0;
        width: auto;
        max-width: 90%;
        z-index: 1000;
    }
</style>
