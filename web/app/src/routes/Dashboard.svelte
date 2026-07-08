<script lang="ts">
  import { onMount } from 'svelte'
  import { api } from '../lib/api'
  import type { DashboardData } from '../lib/types'

  let data = $state<DashboardData | null>(null)
  let error = $state('')

  onMount(async () => {
    try {
      data = await api.get<DashboardData>('/dashboard')
    } catch {
      error = 'Failed to load the dashboard.'
    }
  })
</script>

<section class="page">
  <h1 class="page-title">Dashboard</h1>

  {#if error}
    <div class="alert" role="alert">{error}</div>
  {:else if data}
    <p class="lead">{data.greeting}</p>
    <div class="cards">
      {#each data.stats as stat (stat.label)}
        <div class="card stat">
          <div class="stat-label">{stat.label}</div>
          <div class="stat-value">{stat.value}</div>
        </div>
      {/each}
    </div>
  {:else}
    <div class="spinner"></div>
  {/if}
</section>
