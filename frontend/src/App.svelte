<script lang="ts">
  import "./style.css"
  import { Connect } from "../wailsjs/go/main/App.js"
  import Remotes from "./components/Remotes.svelte"

  let localPort: number = 8181
  let remotePort: number = 9443
  let remoteHost: string = "localhost"
  let resultText: string = ""

  function listen(): void {
    console.log(localPort, remotePort, remoteHost)
    Connect(localPort, remoteHost, remotePort).catch((err) => {
      resultText = err
    })
  }
</script>

<main>
  <Remotes />
  <div>
    <input
      type="number"
      class="input"
      bind:value={localPort}
      placeholder="Local Port"
    />
    <input
      type="text"
      class="input"
      bind:value={remoteHost}
      placeholder="Remote Host"
    />
    <input
      type="number"
      class="input"
      bind:value={remotePort}
      placeholder="Remote Port"
    />
  </div>

  <button on:click={listen}>Listen</button>

  <pre>
    {resultText}
  </pre>
</main>

<style>
</style>
