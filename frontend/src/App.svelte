<script lang="ts">
  import { Connect } from "../wailsjs/go/main/App.js"

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
  #logo {
    display: block;
    width: 50%;
    height: 50%;
    margin: auto;
    padding: 10% 0 0;
    background-position: center;
    background-repeat: no-repeat;
    background-size: 100% 100%;
    background-origin: content-box;
  }

  .result {
    height: 20px;
    line-height: 20px;
    margin: 1.5rem auto;
  }

  .input-box .btn {
    width: 60px;
    height: 30px;
    line-height: 30px;
    border-radius: 3px;
    border: none;
    margin: 0 0 0 20px;
    padding: 0 8px;
    cursor: pointer;
  }

  .input-box .btn:hover {
    color: #333333;
  }

  .input-box .input {
    border: none;
    border-radius: 3px;
    outline: none;
    height: 30px;
    line-height: 30px;
    padding: 0 10px;
    background-color: rgba(240, 240, 240, 1);
    -webkit-font-smoothing: antialiased;
  }

  .input-box .input:hover {
    border: none;
    background-color: rgba(255, 255, 255, 1);
  }

  .input-box .input:focus {
    border: none;
    background-color: rgba(255, 255, 255, 1);
  }
</style>
