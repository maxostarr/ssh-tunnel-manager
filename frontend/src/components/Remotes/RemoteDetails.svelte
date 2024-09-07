<script lang="ts">
  import { selectedRemoteStore } from "../../lib/store"
  import { openRemote } from "../../lib/utils"
  import Tunnels from "../Tunnels/Tunnels.svelte"
  import NewTunnel from "../Tunnels/NewTunnel.svelte"

  let showNewTunnel: () => void
</script>

<NewTunnel bind:show={showNewTunnel} />

<div class="flex-1">
  <button
    class="btn btn-primary"
    on:click={() => openRemote($selectedRemoteStore.id)}>Connect</button
  >

  <button class="btn btn-primary" on:click={showNewTunnel}>New Tunnel</button>
  <Tunnels />
</div>
