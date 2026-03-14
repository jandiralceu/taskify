import { createQuery } from '@tanstack/svelte-query';
import { userService } from '$lib/api/user.service';
import { storage, AUTH_KEYS } from '$lib/utils/storage';

export const PROFILE_QUERY_KEY = ['profile'] as const;

// Criamos um estado reativo para o token
export const authState = $state({
	token: storage.get(AUTH_KEYS.ACCESS_TOKEN)
});

export function createProfileQuery() {
	return createQuery(() => ({
		queryKey: PROFILE_QUERY_KEY,
		queryFn: () => userService.getProfile(),
		// Agora o TanStack Query "vigia" o authState.token
		enabled: !!authState.token,
		staleTime: 1000 * 60 * 10,
		retry: 1
	}));
}
