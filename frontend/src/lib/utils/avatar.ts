/**
 * Normalises an avatar URL so it is always root-relative or absolute.
 *
 * Handles three cases:
 *  - null / undefined  → null
 *  - already absolute  → returned as-is   (e.g. "http://...")
 *  - root-relative     → returned as-is   (e.g. "/uploads/...")
 *  - bare-relative     → prepends "/"     (e.g. "uploads/..." → "/uploads/...")
 */
export function resolveAvatarUrl(
  avatarUrl: string | null | undefined
): string | null {
  if (!avatarUrl) return null
  if (avatarUrl.startsWith('http') || avatarUrl.startsWith('/')) return avatarUrl
  return `/${avatarUrl}`
}
