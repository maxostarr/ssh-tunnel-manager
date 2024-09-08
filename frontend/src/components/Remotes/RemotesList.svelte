<script lang="ts">
  import { PlusIcon } from "lucide-svelte"
  import {
    loadRemotes,
    remotesStore,
    selectRemote,
    selectedRemoteStore,
  } from "../../lib/store"

  import { onMount } from "svelte"
  import NewRemote from "./NewRemote.svelte"
  import { openRemote } from "../../lib/utils"
  import RemoteContext from "./RemoteContext.svelte"
  import type { ssh_manager } from "../../../wailsjs/go/models"

  let debugRemotes: string = ""
  let showNewRemote: () => void
  let openContextMenu: (
    remote: ssh_manager.SshManagerRemoteData,
    event: MouseEvent,
  ) => void

  onMount(async () => {
    await loadRemotes()
  })

  remotesStore.subscribe((value) => {
    debugRemotes = JSON.stringify(value, null, 2)
  })
</script>

<div class="card-bordered w-96 bg-base-100 shadow-xl h-full">
  <NewRemote bind:show={showNewRemote} />
  <RemoteContext bind:openContextMenu />
  <div class="card-body">
    <h2 class="flex justify-between items-center">
      <span>Remotes</span>

      <button
        class="btn btn-primary btn-circle btn-outline"
        on:click={showNewRemote}
      >
        <PlusIcon />
      </button>
    </h2>
    <div class="divider"></div>
    <div class="">
      <table class="table w-full table-fixed">
        <thead>
          <tr>
            <th>Name</th>
            <th>User</th>
            <th>Host</th>
          </tr>
        </thead><tbody>
          {#each $remotesStore as remote}
            <tr
              on:click={() => selectRemote(remote)}
              on:dblclick={() => openRemote(remote.id)}
              on:contextmenu|preventDefault={(event) =>
                openContextMenu(remote, event)}
              tabindex="0"
              role="button"
              class="cursor-pointer"
              class:selected={remote.id === $selectedRemoteStore.id}
            >
              <td>
                <h2>{remote.name}</h2>
              </td>

              <td>
                <p class="text-ellipsis overflow-hidden">{remote.username}</p>
              </td>

              <td class="tooltip w-full text-left" data-tip={remote.host}>
                <p class="text-ellipsis overflow-hidden">{remote.host}</p>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  </div>
</div>

<style>
  .table tr:hover {
    background-color: oklch(var(--p));
    color: oklch(var(--pc));
  }

  .table tr.selected {
    /* background-color: oklch(var(--p));
    color: oklch(var(--pc)); */

    outline: 1px solid oklch(var(--p));
  }
</style>
