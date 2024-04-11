import { writable } from "svelte/store"
import {
  AddRemote,
  AddTunnel,
  GetRemotes,
  GetTunnels,
} from "../../wailsjs/go/main/App.js"

export type Remote = Awaited<ReturnType<typeof GetRemotes>>[0]
export type Tunnel = Awaited<ReturnType<typeof GetTunnels>>[0]

export type NewRemote = Omit<Remote, "ID">
export type NewTunnel = Omit<Tunnel, "ID">

export const remotesStore = writable([] as Remote[])
export const selectedRemoteStore = writable(null as Remote | null)
export const tunnelsStore = writable([] as Tunnel[])

export const loadRemotes = async () => {
  const remotesData = await GetRemotes()
  const singleRemote = remotesData[0]
  const test = singleRemote.test

  remotesStore.set(remotesData)
}

export const loadTunnels = async (remoteId: string) => {
  const tunnelsData = (await GetTunnels(remoteId)) || []
  tunnelsStore.set(tunnelsData)
}

export const addRemote = async (remote: NewRemote) => {
  await AddRemote(remote.Name, remote.Host, remote.Port, remote.Username)
  return loadRemotes()
}

export const addTunnel = async (remote: Remote, tunnel: NewTunnel) => {
  await AddTunnel(
    remote.ID,
    tunnel.LocalPort,
    tunnel.RemoteHost,
    tunnel.RemotePort,
  )
  return loadTunnels(tunnel.RemoteID)
}

export const selectRemote = (remote: Remote) => {
  selectedRemoteStore.set(remote)
  loadTunnels(remote.ID)
}
