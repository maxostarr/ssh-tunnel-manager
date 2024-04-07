<script lang="ts">
  import "./style.css"
  import Remotes from "./components/Remotes/Remotes.svelte"
</script>

<main class="h-full grid">
  <Remotes />
</main>

<style>
  main {
    grid-template-columns: minmax(0, 1fr) 5fr;
  }
</style>
