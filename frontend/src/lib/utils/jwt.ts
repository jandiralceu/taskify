/**
 * Decodes a JWT token without verifying the signature (standard browser-side decoding).
 */
export function decodeToken<T = unknown>(token: string | null): T | null {
  if (!token) return null
  
  try {
    const base64Url = token.split('.')[1]
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/')
    const jsonPayload = decodeURIComponent(
      atob(base64)
        .split('')
        .map(c => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2))
        .join('')
    )
    
    return JSON.parse(jsonPayload)
  } catch (e) {
    console.error('Failed to decode JWT token:', e)
    return null
  }
}

export interface DecodedToken {
  sub: string
  role: string
  permissions: string[]
  type: 'access' | 'refresh'
  iat: number
  exp: number
  jti: string
}
