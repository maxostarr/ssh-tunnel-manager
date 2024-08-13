import { writable } from "svelte/store"
import { EventsEmit, EventsOn, EventsOff } from "../../wailsjs/runtime"

export const prompts = writable([])

export interface Prompt {
  id: number
  message: string
  options: PromptOptions
}

export type PromptOptions =
  | {
      type: "confirm"
      confirmText: string
      cancelText: string
    }
  | {
      type: "prompt"
      confirmText: string
      cancelText: string
      inputs: Array<{
        label: string
        type: "text" | "password"
      }>
    }

const DEFAULT_PROMPT_OPTIONS = {
  type: "confirm",
  confirmText: "OK",
  cancelText: "Cancel",
} as const

// TODO: Move ID management to go side to support replying to events

export const prompt = (
  message: string,
  options: PromptOptions = DEFAULT_PROMPT_OPTIONS,
) => {
  return new Promise<boolean | string[]>((resolve) => {
    const id = Math.floor(Math.random() * 10000)
    const promptData = {
      id,
      message,
      options,
    }

    prompts.update((all) => [promptData, ...all])

    const onConfirm = (result: boolean | string[]) => {
      prompts.update((all) => all.filter((p) => p.id !== id))
      resolve(result)
    }

    EventsOn("prompt:confirm", onConfirm)
  })
}

EventsOn("prompt", async (message: string, promptData: PromptOptions) => {
  await prompt(message, promptData)
})
