<script lang="ts">
  import { AddRemote } from "../../../wailsjs/go/main/App"
  import { loadRemotes, remotesStore } from "../../lib/store"
  import type { Remote } from "../../lib/store"

  import { onMount } from "svelte"
  import NewRemote from "./NewRemote.svelte"

  let remotes: Remote[] = []
  let debugRemotes: string = ""
  let showNewRemote: () => void

  onMount(async () => {
    await loadRemotes()
  })

  remotesStore.subscribe((value) => {
    remotes = value
    debugRemotes = JSON.stringify(value, null, 2)
  })
</script>

<div class="card-bordered w-96 bg-base-100 shadow-xl">
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
    <ol>
      {#each remotes as remote}
        <li>
          <h2>{remote.Name}</h2>
          <p>{remote.Host}</p>
          <!-- <p>{remote.Port}</p>
          <p>{remote.Username}</p>
          <p>{remote.Password}</p> -->
        </li>
      {/each}
    </ol>
  </div>
</div>

<style>
  li,
  .card-title {
    display: flex;
    justify-content: space-between;
  }

  li p {
    flex: 0;
  }

  svg {
    fill: currentColor;
  }
</style>
