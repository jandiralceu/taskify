
/**
 * Checks if a user has a specific permission.
 * Format: "resource:action" (e.g. "/api/v1/users:GET")
 * Supports wildcards like "*:*" or "/api/v1/tasks*:*"
 */
export function hasPermission(
  permissions: string[] | undefined,
  resource: string,
  action: string = '*'
): boolean {
  if (!permissions) return false

  return permissions.some(p => {
    const [pResource, pAction] = p.split(':')
    
    const resourceMatch = matches(pResource, resource)
    const actionMatch = matches(pAction, action)
    
    return resourceMatch && actionMatch
  })
}

function matches(pattern: string, value: string): boolean {
  if (pattern === '*') return true
  
  // Simple wildcard support: /api/v1/tasks*
  if (pattern.endsWith('*')) {
    const prefix = pattern.slice(0, -1)
    return value.startsWith(prefix)
  }
  
  return pattern === value
}

/**
 * Helper to check if a user is an admin based on role or full access permission
 */
export function isAdmin(
  role: string | undefined, 
  permissions: string[] | undefined
): boolean {
  if (!role) return false
  if (role === 'admin') return true
  
  return hasPermission(permissions, '*', '*')
}
