<script lang="ts">
  import { fade } from "svelte/transition"
  import type { ssh_manager } from "../../../wailsjs/go/models"
  import { deleteRemote, updateRemote } from "../../lib/store"
  import { prompt } from "../../lib/promptStore"
  import { onDestroy } from "svelte"

  let clientX: number = 0
  let clientY: number = 0
  let remote: ssh_manager.SshManagerRemoteData = {} as any

  export const openContextMenu: (
    inpRemote: ssh_manager.SshManagerRemoteData,
    mouseEvent: MouseEvent,
  ) => void = (inpRemote, mouseEvent) => {
    registerCloseListeners()
    remote = inpRemote
    clientX = mouseEvent.clientX
    clientY = mouseEvent.clientY
  }

  export const closeContextMenu: () => void = () => {
    console.log("Close context menu")
    deregCloseListeners()
    clientX = 0
    clientY = 0
    remote = null
  }

  const handleRenameRemote: () => Promise<void> = async () => {
    deregCloseListeners()
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
    closeContextMenu()
  }

  const handleEditHost: () => Promise<void> = async () => {
    deregCloseListeners()
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
    closeContextMenu()
  }

  const handleEditUser: () => Promise<void> = async () => {
    deregCloseListeners()
    console.log("Edit user", remote)

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

    console.log("New user", newUser, remote)

    updateRemote({
      ...remote,
      username: newUser,
    })
    closeContextMenu()
  }

  const handleDeleteRemote: () => void = () => {
    deleteRemote(remote.id)
  }

  const checkEsc = (ev: KeyboardEvent) => (any) => {
    return (event) => {
      if (event.key === "Escape") {
        closeContextMenu()
      }
    }
  }

  const checkOutside = (event: MouseEvent) => {
    if (!document.querySelector(".menu").contains(event.target as Node)) {
      closeContextMenu()
    }
  }

  const registerCloseListeners = () => {
    window.addEventListener("click", checkOutside)
    window.addEventListener("keydown", checkEsc)
  }

  registerCloseListeners()

  const deregCloseListeners = () => {
    window.removeEventListener("click", checkOutside)
    window.removeEventListener("keydown", checkEsc)
  }
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
