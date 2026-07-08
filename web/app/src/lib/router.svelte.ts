// A tiny hash-based router. Hash routing needs no server rewrites, which keeps
// the SPA trivially deployable under any base path.

function currentRoute(): string {
  return window.location.hash.slice(1) || '/'
}

class Router {
  route = $state(currentRoute())

  constructor() {
    window.addEventListener('hashchange', () => {
      this.route = currentRoute()
    })
  }

  navigate(to: string): void {
    window.location.hash = to
  }
}

export const router = new Router()
