import { get, writable } from "svelte/store"

export interface Toast {
  id: number
  type: "info" | "success" | "warning" | "error"
  message: string
  dismissible?: boolean
  timeout?: number
}

export const toasts = writable([])

const getId = () => {
  let id = Math.floor(Math.random() * 10000)
  const existingToasts = get(toasts)
  while (existingToasts.find((t) => t.id === id)) {
    id = Math.floor(Math.random() * 10000)
  }
  return id
}

export const addToast = (toast: Omit<Toast, "id">) => {
  // Create a unique ID so we can easily find/remove it
  // if it is dismissible/has a timeout.
  const id = getId()

  // Setup some sensible defaults for a toast.
  const defaults = {
    id,
    type: "info",
    dismissible: true,
    timeout: 3000,
  }

  const merged = { ...defaults, ...toast }

  // Push the toast to the top of the list of toasts
  toasts.update((all) => [merged, ...all])

  // If toast is dismissible, dismiss it after "timeout" amount of time.
  if (merged.timeout) {
    setTimeout(() => dismissToast(id), merged.timeout)
  }
}

export const dismissToast = (id) => {
  toasts.update((all) => all.filter((t) => t.id !== id))
}
