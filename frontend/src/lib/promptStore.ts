import { writable } from "svelte/store"
import { EventsEmit, EventsOn, EventsOff } from "../../wailsjs/runtime"

export const prompts = writable<Array<PromptOptions>>([])

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
      label: string
      id: number
    }
  | {
      type: "prompt"
      label: string
      id: number
      confirmText: string
      cancelText: string
      inputs: Array<{
        label: string
        type: "text" | "password"
      }>
    }

export const DEFAULT_PROMPT_OPTIONS: PromptOptions = {
  type: "confirm",
  confirmText: "OK",
  cancelText: "Cancel",
  id: 0,
  label: "Are you sure?",
} as const

// TODO: Move ID management to go side to support replying to events

export const prompt = (options: PromptOptions = DEFAULT_PROMPT_OPTIONS) => {
  return new Promise<boolean | string[]>((resolve) => {
    const id = Math.floor(Math.random() * 10000)
    const promptData = {
      id,
      ...options,
    }

    prompts.update((all) => [...all, promptData])

    const onConfirm = (result: boolean | string[]) => {
      prompts.update((all) => all.filter((p) => p.id !== id))
      resolve(result)
    }

    EventsOn("prompt:confirm", onConfirm)
  })
}

EventsOn("prompt", async (promptData: PromptOptions) => {
  await prompt(promptData)
})
