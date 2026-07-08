<script lang="ts">
  import { onMount } from 'svelte'
  import { auth } from './lib/auth.svelte'
  import { router } from './lib/router.svelte'
  import Login from './routes/Login.svelte'
  import Dashboard from './routes/Dashboard.svelte'
  import Admin from './routes/Admin.svelte'

  onMount(() => {
    void auth.load()
  })

  const onDashboard = $derived(router.route !== '/admin')
</script>

{#if auth.loading}
  <div class="center">
    <div class="spinner" aria-label="Loading"></div>
  </div>
{:else if !auth.user}
  <Login />
{:else}
  <div class="shell">
    <header class="topbar">
      <div class="brand">nova</div>
      <nav class="nav">
        <a class:active={onDashboard} href="#/dashboard">Dashboard</a>
        {#if auth.user.role === 'admin'}
          <a class:active={router.route === '/admin'} href="#/admin">Admin</a>
        {/if}
      </nav>
      <div class="account">
        <span class="account-email">{auth.user.email}</span>
        <button class="btn ghost sm" onclick={() => auth.logout()}>Log out</button>
      </div>
    </header>

    <main class="content">
      {#if router.route === '/admin' && auth.user.role === 'admin'}
        <Admin />
      {:else}
        <Dashboard />
      {/if}
    </main>
  </div>
{/if}
