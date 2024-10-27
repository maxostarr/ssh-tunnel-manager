<script lang="ts">
  import { selectedRemoteStore } from "../../lib/store"
  import { openRemote, closeRemote } from "../../lib/utils"
  import Tunnels from "../Tunnels/Tunnels.svelte"
  import NewTunnel from "../Tunnels/NewTunnel.svelte"

  let showNewTunnel: () => void
</script>

<NewTunnel bind:show={showNewTunnel} />

<div class="flex-1 p-2">
  {#if $selectedRemoteStore.status !== "connected"}
    <button
      class="btn btn-primary"
      on:click={() => openRemote($selectedRemoteStore.id)}>Connect</button
    >
  {/if}
  {#if $selectedRemoteStore.status === "connected"}
    <button
      class="btn btn-primary"
      on:click={() => closeRemote($selectedRemoteStore.id)}>Disconnect</button
    >
  {/if}

  <button class="btn btn-primary" on:click={showNewTunnel}>New Tunnel</button>
  <Tunnels />
</div>
