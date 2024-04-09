<script lang="ts">
  import Toast from "./Toast.svelte"

  import { dismissToast, toasts } from "../../lib/toastStore"
  import { fly } from "svelte/transition"
  import { flip } from "svelte/animate"
</script>

<section class="toast toast-center toast-bottom">
  {#each $toasts as toast (toast.id)}
    <div
      in:fly={{ y: 200, duration: 200 }}
      out:fly={{ x: -100, duration: 200 }}
      animate:flip={{ duration: 200 }}
    >
      <Toast
        type={toast.type}
        dismissible={toast.dismissible}
        on:dismiss={() => dismissToast(toast.id)}
      >
        {toast.message}
      </Toast>
    </div>
  {/each}
</section>
