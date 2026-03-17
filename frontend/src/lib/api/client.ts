import ky from 'ky'
import { PUBLIC_API_URL } from '$env/static/public'
import { storage, AUTH_KEYS } from '$lib/utils/storage'

const PUBLIC_ENDPOINTS = ['auth/signin', 'auth/register']

export const api = ky.create({
  prefixUrl: PUBLIC_API_URL,
  hooks: {
    beforeRequest: [
      request => {
        const isPublic = PUBLIC_ENDPOINTS.some(endpoint =>
          request.url.includes(endpoint)
        )
        if (isPublic) return

        const token = storage.get(AUTH_KEYS.ACCESS_TOKEN)
        if (token) {
          request.headers.set('Authorization', `Bearer ${token}`)
        }
      },
    ],
    afterResponse: [
      async (request, options, response) => {
        // We only care about 401 Unauthorized errors
        if (response.status !== 401) return

        // Don't attempt refresh for public endpoints (like signin itself)
        const isPublic = PUBLIC_ENDPOINTS.some(endpoint =>
          request.url.includes(endpoint)
        )
        if (isPublic) return

        const refreshToken = storage.get(AUTH_KEYS.REFRESH_TOKEN)
        if (!refreshToken) {
          handleUnauthorized()
          return
        }

        try {
          // Call refresh endpoint directly to avoid circular dependency
          const newTokens = await ky
            .post(`${PUBLIC_API_URL}/auth/refresh`, {
              json: { refreshToken },
            })
            .json<{ accessToken: string; refreshToken: string }>()

          // Save new tokens
          storage.set(AUTH_KEYS.ACCESS_TOKEN, newTokens.accessToken)
          storage.set(AUTH_KEYS.REFRESH_TOKEN, newTokens.refreshToken)

          // Retry original request with new token
          request.headers.set(
            'Authorization',
            `Bearer ${newTokens.accessToken}`
          )
          return ky(request, options)
        } catch (refreshError) {
          console.error('Session expired. Please log in again.', refreshError)
          handleUnauthorized()
        }
      },
    ],
  },
})

function handleUnauthorized() {
  storage.remove(AUTH_KEYS.ACCESS_TOKEN)
  storage.remove(AUTH_KEYS.REFRESH_TOKEN)
  if (typeof window !== 'undefined') {
    window.location.href = '/signin'
  }
}
