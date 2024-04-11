<script lang="ts">
  import { loadRemotes, remotesStore, selectRemote } from "../../lib/store"

  import { onMount } from "svelte"
  import NewRemote from "./NewRemote.svelte"
  import { openRemote } from "../../lib/utils"

  let debugRemotes: string = ""
  let showNewRemote: () => void

  onMount(async () => {
    await loadRemotes()
  })

  remotesStore.subscribe((value) => {
    debugRemotes = JSON.stringify(value, null, 2)
  })
</script>

<div class="card-bordered w-96 bg-base-100 shadow-xl h-full">
  <NewRemote bind:show={showNewRemote} />
  <div class="card-body">
    <h2 class="card-title">
      Remotes

      <button
        class="btn btn-primary btn-circle btn-outline"
        on:click={showNewRemote}
      >
        <svg
          height="24px"
          width="24px"
          version="1.1"
          id="Layer_1"
          xmlns="http://www.w3.org/2000/svg"
          xmlns:xlink="http://www.w3.org/1999/xlink"
          viewBox="0 0 455 455"
          xml:space="preserve"
        >
          <polygon
            points="455,212.5 242.5,212.5 242.5,0 212.5,0 212.5,212.5 0,212.5 0,242.5 212.5,242.5 212.5,455 242.5,455 242.5,242.5 
	455,242.5 "
          />
        </svg>
      </button>
    </h2>
    <div class="divider"></div>
    <div class="overflow-x-auto">
      <table class="table">
        <tbody>
          {#each $remotesStore as remote}
            <tr
              on:click={() => selectRemote(remote)}
              on:dblclick={() => openRemote(remote.id)}
              tabindex="0"
              role="button"
            >
              <td>
                <h2>{remote.name}</h2>
              </td>
              <td>
                <p>{remote.host}</p>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  </div>
</div>

<style>
  .card-title {
    display: flex;
    justify-content: space-between;
  }
  svg {
    fill: currentColor;
  }

  .table tr:hover {
    background-color: oklch(var(--p));
    color: oklch(var(--pc));
  }
</style>
