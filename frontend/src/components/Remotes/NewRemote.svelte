<script lang="ts">
  import { addRemote, type NewRemote } from "../../lib/store"

  import { onMount } from "svelte"

  const saveNewRemote = async (event: Event) => {
    event.preventDefault()
    const form = event.target as HTMLFormElement
    const formData = new FormData(form)
    const data = Object.fromEntries(formData.entries())

    const newRemote: NewRemote = {
      name: data.name.toString(),
      host: data.host.toString(),
      port: parseInt(data.port.toString()),
      username: data.username.toString(),
    }

    await addRemote(newRemote)
    close()
  }

  export const show = () => {
    ;(document.getElementById("newRemote") as HTMLDialogElement).showModal()
  }

  export const close = () => {
    ;(document.getElementById("newRemote") as HTMLDialogElement).close()
  }

  // onMount(() => {
  //   show()
  // })
</script>

<dialog class="modal card" id="newRemote">
  <div class="modal-box card-body">
    <h2 class="card-title">New Remote</h2>
    <div class="divider"></div>
    <form
      action="saveNewRemote"
      class="form-control flex gap-2"
      on:submit|preventDefault={saveNewRemote}
    >
      <label class="input input-bordered flex items-center">
        <input
          class="grow"
          type="text"
          name="name"
          placeholder="Name"
          required
        />
      </label>
      <label class="input input-bordered flex items-center">
        <input
          class="grow"
          type="text"
          name="host"
          placeholder="Host"
          required
        />
      </label>
      <label class="input input-bordered flex items-center">
        <input
          class="grow"
          type="number"
          name="port"
          placeholder="Port"
          required
        />
      </label>
      <label class="input input-bordered flex items-center">
        <input
          class="grow"
          type="text"
          name="username"
          placeholder="Username"
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
