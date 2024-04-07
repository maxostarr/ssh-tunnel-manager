<script lang="ts">
  import { EventsEmit, EventsOn } from "../../wailsjs/runtime"
  let resolve
  let reject

  let promptText = ""

  export const prompt = (prompt: string) => {
    promptText = prompt
    const promise = new Promise((res, rej) => {
      resolve = res
      reject = rej
    })
    ;(document.getElementById("newRemote") as HTMLDialogElement).showModal()
    return promise
  }

  const close = () => {
    ;(document.getElementById("newRemote") as HTMLDialogElement).close()
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

  EventsOn("prompt", async (prompt) => {
    console.log("🚀 ~ EventsOn ~ prompt:", prompt)
    const res = await prompt(prompt)
    console.log("🚀 ~ EventsOn ~ res:", res)

    EventsEmit("prompt-response", res)
  })
</script>

<dialog class="modal card" id="newRemote">
  <div class="modal-box card-body">
    <h2 class="card-title">Prompt</h2>
    <div class="divider"></div>
    <form class="form-control flex gap-2" on:submit|preventDefault={submit}>
      <p>{{ prompt }}</p>
      <label class="input input-bordered flex items-center">
        <input type="text" name="response" placeholder="Response" required />
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
