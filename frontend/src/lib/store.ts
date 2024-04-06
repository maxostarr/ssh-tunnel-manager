import { writable } from "svelte/store"
import { AddRemote, GetRemotes } from "../../wailsjs/go/main/App.js"

export interface Remote {
  Name: string
  ID: string
  Host: string
  Port: number
  Username: string
  Password: string
}

export type NewRemote = Omit<Remote, "ID">

export const remotesStore = writable([] as Remote[])

export const loadRemotes = async () => {
  const remotesData = await GetRemotes()
  remotesStore.set(remotesData)
}

export const addRemote = async (remote: NewRemote) => {
  await AddRemote(
    remote.Name,
    remote.Host,
    remote.Port,
    remote.Username,
    remote.Password,
  )
  return loadRemotes()
}
