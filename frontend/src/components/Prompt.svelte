<script lang="ts">
  import { onMount } from "svelte"
  import { EventsEmit, EventsOn, EventsOff } from "../../wailsjs/runtime"
  let resolve
  let reject

  let promptText = ""

  export const prompt = (promptString: string) => {
    promptText = promptString
    const promise = new Promise((res, rej) => {
      resolve = res
      reject = rej
    })
    ;(document.getElementById("prompt") as HTMLDialogElement).showModal()
    return promise
  }

  const close = () => {
    ;(document.getElementById("prompt") as HTMLDialogElement).close()
  }

  const submit = async (event: Event) => {
    event.preventDefault()
    const form = event.target as HTMLFormElement
    const formData = new FormData(form)
    const data = Object.fromEntries(formData.entries())

    resolve(data.response)
    close()
  }

  const cancel = () => {
    reject("cancelled")
    close()
  }

  onMount(() => {
    EventsOn("prompt", async (promptString) => {
      const res = await prompt(promptString)

      EventsEmit("prompt-response", res)
    })

    return () => EventsOff("prompt")
  })
</script>

<dialog class="modal card" id="prompt">
  <div class="modal-box card-body">
    <h2 class="card-title">Prompt</h2>
    <div class="divider"></div>
    <form class="form-control flex gap-2" on:submit|preventDefault={submit}>
      <p>{promptText}</p>
      <label class="input input-bordered flex items-center grow">
        <input
          type="password"
          name="response"
          placeholder="Response"
          required
        />
      </label>
      <div class="join">
        <button class="btn join-item flex-1 btn-primary" type="submit"
          >Submit</button
        >
        <button
          class="btn join-item flex-1 btn-primary btn-outline"
          type="button"
          on:click={() => cancel()}>Cancel</button
        >
      </div>
    </form>
  </div>
</dialog>
