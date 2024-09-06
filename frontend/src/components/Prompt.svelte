<script lang="ts">
  import { onMount } from "svelte"
  import { EventsEmit, EventsOn, EventsOff } from "../../wailsjs/runtime"
  import {
    DEFAULT_PROMPT_OPTIONS,
    PromptOptions,
    prompts,
  } from "../lib/promptStore"
  let resolve
  let reject

  let promptText = ""
  let config = {}

  export const prompt = (inpConfig: PromptOptions = DEFAULT_PROMPT_OPTIONS) => {
    promptText = inpConfig.label
    config = inpConfig

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

  prompts.subscribe(async (value) => {
    if (value === null || value.length === 0) {
      return
    }

    const inpConfig = value[0]
    await prompt(inpConfig)
  })

  // onMount(() => {
  //   EventsOn("prompt", async (promptString) => {
  //     const res = await prompt(promptString, {
  //       type: "password",
  //     }).catch((err) => null)

  //     if (res === null) {
  //       EventsEmit("prompt-response", "cancelled", "")
  //       return
  //     }

  //     EventsEmit("prompt-response", "submitted", res)
  //   })

  //   return () => EventsOff("prompt")
  // })
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
