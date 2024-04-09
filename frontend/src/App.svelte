<script lang="ts">
  import "./style.css"
  import Remotes from "./components/Remotes/Remotes.svelte"
  import Prompt from "./components/Prompt.svelte"
  import Toasts from "./components/Toast/Toasts.svelte"
  import { addToast } from "./lib/toastStore"

  const addToastHandler = () => {
    addToast({
      message:
        "This is a toast" +
        Array.from({ length: Math.floor(10 * Math.random()) })
          .map(() => "!")
          .join(""),
      type: (["info", "success", "warning", "error"] as const)[
        Math.floor(4 * Math.random())
      ],
      dismissible: true,
      timeout: null,
    })
  }
</script>

<Prompt />
<Toasts />

<main class="h-full grid">
  <Remotes />
  <button on:click={addToastHandler}>Add Toast</button>
</main>

<style>
  main {
    grid-template-columns: minmax(0, 1fr) 5fr;
  }
</style>
