<script lang="ts">
    // This comomnet is used to display a toast notification
    // it makes use of the sveltestrap Alert component to do so

    import { fade } from 'svelte/transition'

    let { type, dismissible, dismiss, children } = $props()
    let color: string = $state("primary");

    switch (type) {
        case 'alert':
            color = "success";
            break;
        case 'error':
            color = "danger";
            break;
    }
    const handleDismiss = () => {
        if (dismissible) {
            dismiss()
        }
    }
</script>

<article transition:fade onclick={ handleDismiss }>
    <div class={ `toast ${color}` }>
        {@render children()}
    </div>
</article>

<style>
.toast {
    padding: 4px 8px;
    margin: 4px;
    border-radius: 4px;
    box-shadow: 2px 7px 8px -1px rgba(102,102,102,1);
}
.alert {
    color: blue;
}
.success {
    background: green;
    color: black;
}
.danger {
    color: red;
}
</style>
