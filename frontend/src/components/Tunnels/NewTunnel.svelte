<script lang="ts">
  import {
    selectedRemoteStore,
    type NewTunnel,
    addTunnel,
  } from "../../lib/store"

  const saveNewTunnel = async (event: Event) => {
    event.preventDefault()
    const form = event.target as HTMLFormElement
    const formData = new FormData(form)
    const data = Object.fromEntries(formData.entries())

    const newTunnel: NewTunnel = {
      local_port: parseInt(data.local_port.toString()),
      remote_host: data.remote_host.toString(),
      remote_port: parseInt(data.remote_port.toString()),
      remote_id: $selectedRemoteStore.id,
    }

    await addTunnel(newTunnel)
    close()
  }

  export const show = () => {
    ;(document.getElementById("newTunnel") as HTMLDialogElement).showModal()
  }

  export const close = () => {
    ;(document.getElementById("newTunnel") as HTMLDialogElement).close()
  }

  // onMount(() => {
  //   show()
  // })
</script>

<dialog class="modal card" id="newTunnel">
  <div class="modal-box card-body">
    <h2 class="card-title">New Remote</h2>
    <div class="divider"></div>
    <form
      action="saveNewRemote"
      class="form-control flex gap-2"
      on:submit|preventDefault={saveNewTunnel}
    >
      <label class="input input-bordered flex items-center">
        <input
          class="grow"
          type="number"
          name="local_port"
          placeholder="Local Port"
          required
        />
      </label>
      <label class="input input-bordered flex items-center">
        <input
          class="grow"
          type="text"
          name="remote_host"
          placeholder="Remote Host"
          required
        />
      </label>
      <label class="input input-bordered flex items-center">
        <input
          class="grow"
          type="number"
          name="remote_port"
          placeholder="Remote Port"
          required
        />
      </label>

      <div class="join">
        <button class="btn join-item flex-1 btn-primary" type="submit"
          >Save</button
        >
        <button
          class="btn join-item flex-1 btn-primary btn-outline"
          type="button"
          on:click={() => close()}>Cancel</button
        >
      </div>
    </form>
  </div>
</dialog>
