<script lang="ts">
  import { tunnelsStore } from "../../lib/store"
</script>

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
      {#each $tunnelsStore as tunnel, i}
        <tr>
          <td>{i + 1}</td>
          <td>{tunnel.LocalPort}</td>
          <td>{tunnel.RemoteHost}</td>
          <td>{tunnel.RemotePort}</td>
        </tr>
      {/each}
    </tbody>
  </table>
</div>
