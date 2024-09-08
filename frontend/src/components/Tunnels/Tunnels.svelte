<script lang="ts">
  import { selectedRemoteStore } from "../../lib/store"
</script>

{#if !$selectedRemoteStore.tunnels.length}
  <div class="card bg-neutral w-max h-max absolute inset-1/2">
    <div class="card-body">
      <h2 class="text-center card-title">
        No tunnels configured for this host
      </h2>
      <p class="text-center">
        Click the
        <span class="text-primary"> New Tunnel </span>

        button to create one!
      </p>
    </div>
  </div>
{/if}

{#if $selectedRemoteStore.tunnels.length}
  <div class="overflow-x-auto">
    <table class="table">
      <!-- head -->
      <thead>
        <tr>
          <th></th>
          <th>Local Port</th>
          <th>Remote Host</th>
          <th>Remote Port</th>
        </tr>
      </thead>
      <tbody>
        <!-- body -->
        {#each $selectedRemoteStore.tunnels as tunnel, i}
          <tr>
            <td>{i + 1}</td>
            <td>{tunnel.local_port}</td>
            <td>{tunnel.remote_host}</td>
            <td>{tunnel.remote_port}</td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
{/if}
