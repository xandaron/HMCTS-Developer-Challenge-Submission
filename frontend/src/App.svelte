<script lang="ts">
  import "./app.css";
  import { onMount } from "svelte";

  let message: string = "Loading...";

  async function fetchMessage(): Promise<void> {
    try {
      const res = await fetch("http://localhost:8080/api/message");
      const data: { text: string } = await res.json();
      message = data.text;
    } catch (error) {
      message = "Error fetching data";
      console.error(error);
    }
  }

  onMount(fetchMessage);
</script>

<main>
  <h1 class="text-3xl font-bold text-center">Hello world!</h1>
  <h1 class="text-3xl font-bold text-center text-blue-600">{message}</h1>

  <style lang="postcss">
    @reference "tailwindcss";

    :global(html) {
      background-color: theme(--color-gray-100);
    }
  </style>
</main>
