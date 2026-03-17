export const storage = {
  get(key: string): string | null {
    return localStorage.getItem(key)
  },

  set(key: string, value: string): void {
    localStorage.setItem(key, value)
  },

  remove(key: string): void {
    localStorage.removeItem(key)
  },

  clear(): void {
    localStorage.clear()
  },
}

export const AUTH_KEYS = {
  ACCESS_TOKEN: 'accessToken',
  REFRESH_TOKEN: 'refreshToken',
} as const
