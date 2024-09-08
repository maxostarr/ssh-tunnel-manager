import { Connect } from "../../wailsjs/go/app/App"
import { addToast } from "./toastStore"

export const openRemote = async (id: string) => {
  console.log("🚀 ~ openRemote ~ id:", id)
  await Connect(id).catch((err) => {
    console.error(err)
    addToast({
      message: `Failed to connect to remote: ${err}`,
      type: "error",
      dismissible: true,
      timeout: 5000,
    })
  })
}
