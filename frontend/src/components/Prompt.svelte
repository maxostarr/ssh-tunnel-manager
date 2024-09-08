<script lang="ts">
  import type { ConfirmData, PromptData } from "../lib/promptStore"
  import { prompts } from "../lib/promptStore"
  let resolve
  let reject

  let promptText = ""
  let config: PromptData<any> | ConfirmData = {} as any

  export const prompt = (inpConfig: PromptData<any> | ConfirmData) => {
    promptText = inpConfig.label
    config = inpConfig

    resolve = inpConfig.resolve
    reject = inpConfig.reject
    ;(document.getElementById("prompt") as HTMLDialogElement).showModal()
  }

  const close = () => {
    ;(document.getElementById("prompt") as HTMLDialogElement).close()
    reset()
  }

  const submit = async (event: Event) => {
    event.preventDefault()
    const form = event.target as HTMLFormElement
    const formData = new FormData(form)
    const data = Object.fromEntries(formData.entries())

    resolve(data)
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
    prompt(inpConfig)
  })

  const reset = () => {
    promptText = ""
    config = {} as any
  }

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
      {#if config.type === "prompt"}
        {#each config.inputs as input, i}
          <label class="input input-bordered flex items-center grow gap-2">
            {#if input.label}
              {input.label}
            {/if}
            <input
              type={input.type}
              class="grow"
              name={input.key}
              required
              autofocus={i === 0}
            />
          </label>
        {/each}
      {/if}
      <div class="join">
        <button class="btn join-item flex-1 btn-primary" type="submit"
          >{config.confirmText}</button
        >
        <button
          class="btn join-item flex-1 btn-primary btn-outline"
          type="button"
          on:click={() => cancel()}>{config.cancelText}</button
        >
      </div>
    </form>
  </div>
</dialog>
