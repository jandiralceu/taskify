import { redirect } from '@sveltejs/kit'
import { storage, AUTH_KEYS } from '$lib/utils/storage'

export const load = () => {
  const token = storage.get(AUTH_KEYS.ACCESS_TOKEN)
  if (token) {
    throw redirect(302, '/')
  }
}
