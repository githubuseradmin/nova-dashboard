export type Role = 'user' | 'admin'

export interface User {
  id: number
  email: string
  name: string
  role: Role
  createdAt: string
}

export interface Stat {
  label: string
  value: string
}

export interface DashboardData {
  greeting: string
  stats: Stat[]
  user: User
}
