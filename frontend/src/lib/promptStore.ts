import { writable } from "svelte/store"
import { EventsEmit, EventsOn, EventsOff } from "../../wailsjs/runtime"

export const prompts = writable<Array<PromptData | ConfirmData>>([])

export interface Prompt {
  id: number
  message: string
  options: PromptOptions
}

type ConfirmOptions = {
  type: "confirm"
  label: string
  id: string
  confirmText: string
  cancelText: string
}

type PromptOptions = Omit<ConfirmOptions, "type"> & {
  type: "prompt"
  inputs: Array<{
    label: string
    key: string
    type: "text" | "password"
  }>
}

// export type PromptOptions = ConfirmOptions | NewType
export type UIPromptOptions = Omit<PromptOptions, "id">
export type UIConfirmOptions = Omit<ConfirmOptions, "id">

export type CapitalizeKeys<T> = {
  [K in keyof T as Capitalize<string & K>]: T[K]
}

export type BackendPromptOptions = {
  ID: string
  Data: CapitalizeKeys<
    Omit<UIPromptOptions, "inputs"> & {
      inputs: Array<CapitalizeKeys<PromptOptions["inputs"][0]>>
    }
  >
}

export type PromptData = PromptOptions & {
  resolve: (result: string[]) => void
  reject: (reason: any) => void
}

export type ConfirmData = ConfirmOptions & {
  resolve: (result: boolean) => void
  reject: (reason: any) => void
}

const handleConfirm = async (options: ConfirmOptions) => {
  const result = await new Promise<boolean>((resolve, reject) => {
    const confirmData = {
      ...options,
      resolve,
      reject,
    }

    prompts.update((all) => [...all, confirmData])

    // const onConfirm = (result: boolean | string[]) => {
    //   prompts.update((all) => all.filter((p) => p.id !== id))
    //   resolve(result)
    // }
  }).finally(() => {
    prompts.update((all) => all.filter((p) => p.id !== options.id))
  })
  console.log("🚀 ~ result ~ result:", result)

  return result
}

export const confirm = (options: UIConfirmOptions) => {
  const id = Math.random().toString(36).substr(2, 9)
  return handleConfirm({ ...options, id })
}

const handlePrompt = async (options: PromptOptions) => {
  const result = await new Promise<string[]>((resolve, reject) => {
    const promptData = {
      ...options,
      resolve,
      reject,
    }

    prompts.update((all) => [...all, promptData])

    // const onConfirm = (result: boolean | string[]) => {
    //   prompts.update((all) => all.filter((p) => p.id !== id))
    //   resolve(result)
    // }
  }).finally(() => {
    prompts.update((all) => all.filter((p) => p.id !== options.id))
  })

  return result
}

export const prompt = (options: UIPromptOptions) => {
  const id = Math.random().toString(36).substr(2, 9)
  return handlePrompt({ ...options, id })
}

const backendToPromptOptions = (
  options: BackendPromptOptions,
): PromptOptions => {
  return {
    id: options.ID,
    type: options.Data.Type,
    label: options.Data.Label,
    confirmText: options.Data.ConfirmText,
    cancelText: options.Data.CancelText,
    inputs:
      options.Data.Type === "prompt" &&
      options.Data.Inputs.map((input) => ({
        label: input.Label,
        type: input.Type,
        key: input.Key,
      })),
  }
}

EventsOn("prompt", async (promptData: BackendPromptOptions) => {
  console.log("🚀 ~ EventsOn ~ promptData:", promptData)
  const result = await handlePrompt(backendToPromptOptions(promptData))
  console.log("🚀 ~ EventsOn ~ result:", result)

  EventsEmit("prompt" + promptData.ID, "success", result)
})
