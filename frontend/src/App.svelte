<script lang="ts">
  import "./style.css"
  import RemotesList from "./components/Remotes/RemotesList.svelte"
  import Prompt from "./components/Prompt.svelte"
  import Toasts from "./components/Toast/Toasts.svelte"
  import RemoteDetails from "./components/Remotes/RemoteDetails.svelte"
  import { selectedRemoteStore } from "./lib/store"
</script>

<Prompt />
<Toasts />

<main class="h-full grid">
  <RemotesList />
  {#if $selectedRemoteStore}
    <RemoteDetails />
  {:else}
    <div class="flex items-center justify-center">
      <h1 class="text-3xl text-gray-500">Select a remote to connect</h1>
    </div>
  {/if}
</main>

<style>
  main {
    grid-template-columns: minmax(0, 1fr) 5fr;
  }
</style>
