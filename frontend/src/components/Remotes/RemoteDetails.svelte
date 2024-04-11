<script lang="ts">
  import { selectedRemoteStore } from "../../lib/store"
  import { openRemote } from "../../lib/utils"
  import Tunnels from "../Tunnels/Tunnels.svelte"
</script>

<div>
  <button
    class="btn btn-primary"
    on:click={() => openRemote($selectedRemoteStore.id)}>Connect</button
  >
  <Tunnels />
</div>
