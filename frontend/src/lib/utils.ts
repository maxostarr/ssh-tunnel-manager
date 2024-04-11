import { Connect } from "../../wailsjs/go/main/App"
import { addToast } from "./toastStore"

export const openRemote = (id: string) => async () => {
  await Connect(id).catch((err) => {
    console.error(err)
    addToast({
      message: "Failed to connect to remote",
      type: "error",
      dismissible: true,
      timeout: 5000,
    })
  })
}
