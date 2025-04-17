// This file contains a very basic toast notification management framewok

import { writable } from 'svelte/store'
import type { ToastNotification } from './types';

export const toasts = writable<ToastNotification[]>([])

/**
 * Dismiss a currently showing toast notification
 */
export const dismissToast = (id: number | null) => {
    if (!id) {
        return;
    }

    toasts.update(
        all => all.filter(
            (t: ToastNotification) => t.id !== id
        )
    )
}

/**
 * Add a toast notification to the screen
 */
export const addToast = (toast: ToastNotification): void => {
    const id = Math.floor(Math.random() * 10000)

    // Setup some sensible defaults for a toast.
    const defaults = {
        id,
        type: 'info',
        dismissible: true,
        timeout: 12_000,
    }

    // Push the toast to the top of the list of toasts
    const t:ToastNotification = { ...defaults, ...toast }
    toasts.update((all: ToastNotification[]) => {
        return [t, ...all]
    })

    // If toast is dismissible, dismiss it after "timeout" amount of time.
    if (t.timeout) {
        setTimeout(() => dismissToast(id), t.timeout)
    }
}

const error = (message: string) => addToast(<ToastNotification>{type: "error", message});
const alert = (message: string) => addToast(<ToastNotification>{type: "alert", message});
const info = (message: string) => addToast(<ToastNotification>{type: "info", message});

export default {
    error,
    alert,
    info,
    add: addToast,
    dismis: dismissToast
}
