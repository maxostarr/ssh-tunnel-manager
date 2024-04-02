<script lang="ts">
  import { AddRemote } from "../../wailsjs/go/main/App"
  import { loadRemotes, remotesStore } from "../lib/store"

  import { onMount } from "svelte"

  let remotes: any[] = []
  let debugRemotes: string = ""

  const saveNewRemote = async (event: Event) => {
    console.log("🚀 ~ saveNewRemote ~ event:", event)
    event.preventDefault()
    const form = event.target as HTMLFormElement
    const formData = new FormData(form)
    const data = Object.fromEntries(formData.entries())

    await AddRemote(
      data.name.toString(),
      data.host.toString(),
      parseInt(data.port.toString()),
      data.username.toString(),
      data.password.toString(),
    )

    await loadRemotes()
  }

  onMount(async () => {
    await loadRemotes()
  })

  remotesStore.subscribe((value) => {
    remotes = value
    debugRemotes = JSON.stringify(value, null, 2)
  })
</script>

<div>
  {{ debugRemotes }}

  <form action="saveNewRemote" on:submit={saveNewRemote}>
    <input type="text" name="name" placeholder="Name" />
    <input type="text" name="host" placeholder="Host" />
    <input type="number" name="port" placeholder="Port" />
    <input type="text" name="username" placeholder="Username" />
    <input type="password" name="password" placeholder="Password" />
    <button type="submit">Save</button>
  </form>
</div>

<style></style>
