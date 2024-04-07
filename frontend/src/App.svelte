<script lang="ts">
  import "./style.css"
  import Remotes from "./components/Remotes/Remotes.svelte"
  import Prompt from "./components/Prompt.svelte"
</script>

<Prompt />

<main class="h-full grid">
  <Remotes />
</main>

<style>
  main {
    grid-template-columns: minmax(0, 1fr) 5fr;
  }
</style>
