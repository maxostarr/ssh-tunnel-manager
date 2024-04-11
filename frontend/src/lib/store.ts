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

  remotesStore.set(remotesData)
}

export const loadTunnels = async (remoteId: string) => {
  const tunnelsData = (await GetTunnels(remoteId)) || []
  tunnelsStore.set(tunnelsData)
}

export const addRemote = async (remote: NewRemote) => {
  await AddRemote(remote.name, remote.host, remote.port, remote.username)
  return loadRemotes()
}

export const addTunnel = async (remote: Remote, tunnel: NewTunnel) => {
  await AddTunnel(
    remote.id,
    tunnel.local_port,
    tunnel.remote_host,
    tunnel.remote_port,
  )
  return loadTunnels(tunnel.remote_id)
}

export const selectRemote = (remote: Remote) => {
  selectedRemoteStore.set(remote)
  loadTunnels(remote.id)
}
