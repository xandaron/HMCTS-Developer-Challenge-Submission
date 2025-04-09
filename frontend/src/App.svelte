<script lang="ts">
  // import { currentRoute } from "./stores/route";
  import Home from "./routes/Home.svelte";
  import Login from "./routes/Login.svelte";
  import Signup from "./routes/Signup.svelte";

  import { writable } from "svelte/store";

  const currentRoute = writable(window.location.pathname);

  window.addEventListener("popstate", () => {
    currentRoute.set(window.location.pathname);
  });

  function navigate(path: string) {
    window.history.pushState({}, "", path);
    currentRoute.set(path);
  }
</script>

<main class="min-h-screen bg-indigo-950">
  <nav class="flex justify-center space-x-4 mt-4">
    <button on:click={() => navigate("/")} class="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-400">Home</button>
    <button on:click={() => navigate("/login")} class="px-4 py-2 bg-green-500 text-white rounded hover:bg-green-400">Login</button>
    <button on:click={() => navigate("/signup")} class="px-4 py-2 bg-gray-500 text-white rounded hover:bg-gray-400">Signup</button>
  </nav>

  <div>
    {#if $currentRoute === "/"}
      <Home />
    {:else if $currentRoute === "/login"}
      <Login />
    {:else if $currentRoute === "/signup"}
      <Signup />
    {:else}
      <h2>Page not found</h2>
    {/if}
  </div>
</main>

<style lang="postcss">
  @reference "tailwindcss";
</style>