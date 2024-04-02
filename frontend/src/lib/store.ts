import { writable } from "svelte/store"
import { GetRemotes } from "../../wailsjs/go/main/App.js"

export const remotesStore = writable([])

export const loadRemotes = async () => {
  const remotesData = await GetRemotes()
  console.log('🚀 ~ file: store.ts:8 ~ loadRemotes ~ remotesData:', remotesData)
  // remotes.set(remotes)
}

export const addRemote = (remote) => {
  remotesStore.update((remotes) => [...remotes, remote])
}
