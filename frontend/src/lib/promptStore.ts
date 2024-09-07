import { writable } from "svelte/store"
import { EventsEmit, EventsOn, EventsOff } from "../../wailsjs/runtime"

export const prompts = writable<Array<PromptData<any> | ConfirmData>>([])


type ConfirmOptions = {
  type: "confirm"
  label: string
  id: string
  confirmText: string
  cancelText: string
}

type PromptOptions<T extends string> = Omit<ConfirmOptions, "type"> & {
  type: "prompt"
  inputs: Readonly<Array<{
    label: string
    key: T
    type: "text" | "password"
  }>>
}

// export type PromptOptions = ConfirmOptions | NewType
export type UIPromptOptions<T extends string> = Omit<PromptOptions<T>, "id">
export type UIConfirmOptions = Omit<ConfirmOptions, "id">

export type CapitalizeKeys<T> = {
  [K in keyof T as Capitalize<string & K>]: T[K]
}

export type BackendPromptOptions<T extends string> = {
  ID: string
  Data: CapitalizeKeys<
    Omit<UIPromptOptions<T>, "inputs"> & {
      inputs: Array<CapitalizeKeys<PromptOptions<T>["inputs"][0]>>
    }
  >
}

export type PromptData<T extends string> = PromptOptions<T> & {
  resolve: (result: {
    [key in T]: string
  }) => void
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

  return result
}

export const confirm = (options: UIConfirmOptions) => {
  const id = Math.random().toString(36).substr(2, 9)
  return handleConfirm({ ...options, id })
}

const handlePrompt = async <T extends string>(options: PromptOptions<T>) => {
  const result = await new Promise<{
    [key in T]: string
  }>((resolve, reject) => {
    const promptData: PromptData<T> = {
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

export const prompt = <T extends string>(options: UIPromptOptions<T>) => {
  const id = Math.random().toString(36).substr(2, 9)
  return handlePrompt({ ...options, id })
}

const backendToPromptOptions = <T extends string>(
  options: BackendPromptOptions<T>,
): PromptOptions<T> => {
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

EventsOn("prompt", async <T extends string>(promptData: BackendPromptOptions<T>) => {
  const result = await handlePrompt(backendToPromptOptions(promptData))
  EventsEmit("prompt" + promptData.ID, "success", result)
})
