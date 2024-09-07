<script lang="ts">
  import { fade } from "svelte/transition"
  import type { ssh_manager } from "../../../wailsjs/go/models"
  import { deleteRemote, updateRemote } from "../../lib/store"
  import Prompt from "../Prompt.svelte"
  import { prompt } from "../../lib/promptStore"

  let clientX: number = 0
  let clientY: number = 0
  let remote: ssh_manager.SshManagerRemoteData = {} as any

  export const openContextMenu: (
    inpRemote: ssh_manager.SshManagerRemoteData,
    mouseEvent: MouseEvent,
  ) => void = (inpRemote, mouseEvent) => {
    remote = inpRemote
    clientX = mouseEvent.clientX
    clientY = mouseEvent.clientY
  }

  export const closeContextMenu: () => void = () => {
    clientX = 0
    clientY = 0
    remote = null
  }

  const handleRenameRemote: () => Promise<void> = async () => {
    console.log("Rename remote")
    const { newName } = await prompt({
      type: "prompt",
      cancelText: "Cancel",
      confirmText: "Rename",
      label: "Enter new name for remote",
      inputs: [
        {
          type: "text",
          label: "New name",
          key: "newName",
        },
      ],
    })
    updateRemote({
      ...remote,
      name: newName,
    })
  }

  const handleEditHost: () => Promise<void> = async () => {
    console.log("Edit host")

    const { newHost } = await prompt({
      type: "prompt",
      cancelText: "Cancel",
      confirmText: "Edit",
      label: "Enter new host for remote",
      inputs: [
        {
          type: "text",
          label: "New host",
          key: "newHost",
          // placeholder: "New host",
          // required: true,
        },
      ],
    })
    updateRemote({
      ...remote,
      host: newHost,
    })
  }

  const handleEditUser: () => Promise<void> = async () => {
    console.log("Edit user")

    const { newUser } = await prompt({
      type: "prompt",
      cancelText: "Cancel",
      confirmText: "Edit",
      label: "Enter new user for remote",
      inputs: [
        {
          type: "text",
          label: "New user",
          key: "newUser",
          // placeholder: "New user",
          // required: true,
        },
      ],
    })
    updateRemote({
      ...remote,
      username: newUser,
    })
  }

  const handleDeleteRemote: () => void = () => {
    deleteRemote(remote.id)
  }

  // Close context menu on escape key
  window.addEventListener("keydown", (event) => {
    if (event.key === "Escape") {
      closeContextMenu()
    }
  })

  // Close context menu on click outside
  window.addEventListener("click", closeContextMenu)
</script>

<!-- <Prompt bind:prompt /> -->

{#if clientX && clientY}
  <ul
    class="menu bg-base-200 rounded-box z-10 transition-all"
    style="position: absolute; top: {clientY}px; left: {clientX}px;"
    in:fade={{ duration: 150 }}
    out:fade={{ duration: 150 }}
  >
    <li class="menu-title">{remote.name}</li>
    <!-- svelte-ignore a11y-invalid-attribute -->
    <li on:click={handleRenameRemote} on:keydown={handleRenameRemote}>
      <a href="javascript:void(0)">Rename</a>
    </li>
    <!-- svelte-ignore a11y-invalid-attribute -->
    <li on:click={handleEditHost} on:keydown={handleEditHost}>
      <a href="javascript:void(0)">Edit Host</a>
    </li>

    <!-- svelte-ignore a11y-invalid-attribute -->
    <li on:click={handleEditUser} on:keydown={handleEditUser}>
      <a href="javascript:void(0)">Edit User</a>
    </li>

    <!-- svelte-ignore a11y-invalid-attribute -->
    <li
      on:click={handleDeleteRemote}
      on:keydown={handleDeleteRemote}
      class="text-error"
    >
      <a href="javascript:void(0)">Delete</a>
    </li>
  </ul>
{/if}
