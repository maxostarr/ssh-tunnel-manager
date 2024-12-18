import { get, writable } from "svelte/store";
import {
  AddRemote,
  AddTunnel,
  DeleteRemote,
  GetRemote,
  GetRemotes,
  GetTunnels,
  UpdateRemote,
} from "../../wailsjs/go/app/App.js";
import { EventsOn } from "../../wailsjs/runtime/runtime.js";

export type RemoteData = Awaited<ReturnType<typeof GetRemotes>>[0];
export type RemoteFull = Awaited<ReturnType<typeof GetRemote>>;
export type Tunnel = Awaited<ReturnType<typeof GetTunnels>>[0];

export type NewRemote = Omit<RemoteData, "id">;
export type NewTunnel = Omit<Tunnel, "id">;

export const remotesStore = writable([] as RemoteData[]);
export const selectedRemoteStore = writable({ tunnels: [] } as RemoteFull);

export const loadRemotes = async () => {
  const remotesData = await GetRemotes();
  remotesStore.set(remotesData);
};

export const loadRemoteDetails = async (remoteId: string) => {
  const remoteDetails = await GetRemote(remoteId);
  remoteDetails.tunnels = remoteDetails.tunnels ?? [];
  selectedRemoteStore.set(remoteDetails);
};

export const refreshSelectedRemote = async () => {
  const selectedRemote = get(selectedRemoteStore);
  if (!selectedRemote.id) return;
  await loadRemoteDetails(selectedRemote.id);
};

export const addRemote = async (remote: NewRemote) => {
  await AddRemote(remote.name, remote.host, remote.port, remote.username);
  return loadRemotes();
};

export const addTunnel = async (tunnel: NewTunnel) => {
  await AddTunnel(
    tunnel.remote_id,
    tunnel.local_port,
    tunnel.remote_host,
    tunnel.remote_port,
  );
  return loadRemoteDetails(tunnel.remote_id);
};

export const selectRemote = (remote: RemoteData) => {
  loadRemoteDetails(remote.id);
};

export const deleteRemote = async (remoteId: string) => {
  await DeleteRemote(remoteId);
  return loadRemotes();
};

export const updateRemote = async (remote: RemoteData) => {
  console.log("🚀 ~ updateRemote ~ remote:", remote);
  await UpdateRemote(
    remote.id,
    remote.name,
    remote.host,
    remote.port,
    remote.username,
  );
  return loadRemotes();
};

EventsOn("remotes-updated", async () => {
  await loadRemotes();
  await refreshSelectedRemote();
});
