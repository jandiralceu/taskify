import { createQuery, createMutation, useQueryClient } from '@tanstack/svelte-query';
import { userService } from '$lib/api/user.service';
import { storage, AUTH_KEYS } from '$lib/utils/storage';
import { TASKS_QUERY_KEY } from '$lib/state/tasks.svelte';
import type { GetUsersParams, UpdateUserRequest, UserResponse } from '$lib/api/types';

export const PROFILE_QUERY_KEY = ['profile'] as const;
export const PERMISSIONS_QUERY_KEY = ['permissions'] as const;

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

export function createPermissionsQuery() {
	return createQuery(() => ({
		queryKey: PERMISSIONS_QUERY_KEY,
		queryFn: () => userService.getPermissions(),
		enabled: !!authState.token,
		staleTime: 1000 * 60 * 10,
		retry: 1
	}));
}
export function createUserQuery(userId: () => string | undefined) {
	return createQuery(() => ({
		queryKey: ['user', userId()],
		queryFn: () => userService.getUserById(userId()!),
		enabled: !!authState.token && !!userId(),
		staleTime: 1000 * 60 * 5,
		retry: 1
	}));
}

export const USERS_QUERY_KEY = ['users'] as const;

export function updateUserMutation() {
	const queryClient = useQueryClient();

	return createMutation<UserResponse, Error, { id: string; data: UpdateUserRequest }>(() => ({
		mutationFn: ({ id, data }) => userService.updateUser(id, data),
		onSuccess: (updated) => {
			queryClient.invalidateQueries({ queryKey: USERS_QUERY_KEY });
			queryClient.invalidateQueries({ queryKey: ['user', updated.id] });
		}
	}));
}

export function deleteUserMutation() {
	const queryClient = useQueryClient();

	return createMutation<void, Error, string>(() => ({
		mutationFn: (id) => userService.deleteUser(id),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: USERS_QUERY_KEY });
			queryClient.invalidateQueries({ queryKey: TASKS_QUERY_KEY });
		}
	}));
}

export function getUsersQuery(paramsGetter: () => GetUsersParams = () => ({})) {
	return createQuery(() => {
		const params = paramsGetter();
		return {
			queryKey: [...USERS_QUERY_KEY, params],
			queryFn: () => userService.getUsers(params),
			enabled: !!authState.token,
			staleTime: 1000 * 60 * 2,
			retry: 1
		};
	});
}
