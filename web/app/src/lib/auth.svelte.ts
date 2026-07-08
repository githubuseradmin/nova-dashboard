import { api } from './api'
import type { User } from './types'

// AuthStore is a small rune-based store holding the current user. A single
// instance is shared across the app via the `auth` export.
class AuthStore {
  user = $state<User | null>(null)
  loading = $state(true)

  // load restores the session on startup by calling /auth/me.
  async load(): Promise<void> {
    this.loading = true
    try {
      this.user = await api.get<User>('/auth/me')
    } catch {
      this.user = null
    } finally {
      this.loading = false
    }
  }

  async login(email: string, password: string): Promise<void> {
    this.user = await api.post<User>('/auth/login', { email, password })
  }

  async logout(): Promise<void> {
    try {
      await api.post('/auth/logout')
    } catch {
      // Ignore — we clear local state regardless.
    }
    this.user = null
  }
}

export const auth = new AuthStore()
