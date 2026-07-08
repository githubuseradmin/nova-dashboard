<script lang="ts">
  import { auth } from '../lib/auth.svelte'
  import { ApiError } from '../lib/api'

  let email = $state('admin@nova.local')
  let password = $state('')
  let error = $state('')
  let busy = $state(false)

  async function submit(event: Event) {
    event.preventDefault()
    error = ''
    busy = true
    try {
      await auth.login(email, password)
    } catch (err) {
      error = err instanceof ApiError ? err.message : 'Something went wrong'
    } finally {
      busy = false
    }
  }
</script>

<div class="auth-wrap">
  <form class="card auth-card" onsubmit={submit}>
    <div class="auth-brand">nova</div>
    <p class="auth-sub">Sign in to your dashboard</p>

    {#if error}
      <div class="alert" role="alert">{error}</div>
    {/if}

    <label class="field">
      <span>Email</span>
      <input type="email" bind:value={email} autocomplete="username" required />
    </label>

    <label class="field">
      <span>Password</span>
      <input type="password" bind:value={password} autocomplete="current-password" required />
    </label>

    <button class="btn primary" type="submit" disabled={busy}>
      {busy ? 'Signing in…' : 'Sign in'}
    </button>

    <p class="auth-hint">Demo admin is pre-filled — set a password via the seed on first run.</p>
  </form>
</div>
