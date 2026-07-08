<script lang="ts">
  import { onMount } from 'svelte'
  import { api } from '../lib/api'
  import type { User, Role } from '../lib/types'

  let users = $state<User[]>([])
  let error = $state('')
  let pending = $state<number | null>(null)

  async function load() {
    try {
      users = await api.get<User[]>('/admin/users')
    } catch {
      error = 'Failed to load users.'
    }
  }

  onMount(load)

  async function setRole(user: User, role: Role) {
    error = ''
    pending = user.id
    try {
      const updated = await api.put<User>(`/admin/users/${user.id}/role`, { role })
      users = users.map((u) => (u.id === updated.id ? updated : u))
    } catch (err) {
      error = err instanceof Error ? err.message : 'Update failed.'
    } finally {
      pending = null
    }
  }
</script>

<section class="page">
  <h1 class="page-title">Users</h1>

  {#if error}
    <div class="alert" role="alert">{error}</div>
  {/if}

  <div class="card table-card">
    <table class="table">
      <thead>
        <tr>
          <th>ID</th>
          <th>Email</th>
          <th>Name</th>
          <th>Role</th>
          <th class="right">Actions</th>
        </tr>
      </thead>
      <tbody>
        {#each users as user (user.id)}
          <tr>
            <td class="muted">{user.id}</td>
            <td>{user.email}</td>
            <td>{user.name}</td>
            <td><span class="badge {user.role}">{user.role}</span></td>
            <td class="right">
              {#if user.role === 'admin'}
                <button class="btn ghost sm" disabled={pending === user.id} onclick={() => setRole(user, 'user')}>
                  Make user
                </button>
              {:else}
                <button class="btn ghost sm" disabled={pending === user.id} onclick={() => setRole(user, 'admin')}>
                  Make admin
                </button>
              {/if}
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
</section>
