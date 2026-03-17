import {
  createQuery,
  createMutation,
  useQueryClient,
} from '@tanstack/svelte-query'
import { userService } from '$lib/api/user.service'
import { storage, AUTH_KEYS } from '$lib/utils/storage'
import { TASKS_QUERY_KEY } from '$lib/state/tasks.svelte'
import type {
  GetUsersParams,
  UpdateUserRequest,
  UserResponse,
  ChangePasswordRequest,
} from '$lib/api/types'

export const PROFILE_QUERY_KEY = ['profile'] as const
export const PERMISSIONS_QUERY_KEY = ['permissions'] as const

// Criamos um estado reativo para o token
export const authState = $state({
  token: storage.get(AUTH_KEYS.ACCESS_TOKEN),
})

export function createProfileQuery() {
  return createQuery(() => ({
    queryKey: PROFILE_QUERY_KEY,
    queryFn: () => userService.getProfile(),
    // Agora o TanStack Query "vigia" o authState.token
    enabled: !!authState.token,
    staleTime: 1000 * 60 * 10,
    retry: 1,
  }))
}

export function createPermissionsQuery() {
  return createQuery(() => ({
    queryKey: PERMISSIONS_QUERY_KEY,
    queryFn: () => userService.getPermissions(),
    enabled: !!authState.token,
    staleTime: 1000 * 60 * 10,
    retry: 1,
  }))
}
export function createUserQuery(userId: () => string | undefined) {
  return createQuery(() => ({
    queryKey: ['user', userId()],
    queryFn: () => userService.getUserById(userId()!),
    enabled: !!authState.token && !!userId(),
    staleTime: 1000 * 60 * 5,
    retry: 1,
  }))
}

export const USERS_QUERY_KEY = ['users'] as const

export function uploadAvatarMutation() {
  const queryClient = useQueryClient()

  return createMutation<{ avatarUrl: string }, Error, File>(() => ({
    mutationFn: file => userService.uploadAvatar(file),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: PROFILE_QUERY_KEY })
    },
  }))
}

export function updateProfileMutation() {
  const queryClient = useQueryClient()

  return createMutation<UserResponse, Error, UpdateUserRequest>(() => ({
    mutationFn: data => userService.updateProfile(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: PROFILE_QUERY_KEY })
      queryClient.invalidateQueries({ queryKey: USERS_QUERY_KEY })
    },
  }))
}

export function changePasswordMutation() {
  return createMutation<void, Error, ChangePasswordRequest>(() => ({
    mutationFn: data => userService.changePassword(data),
  }))
}

export function deleteProfileMutation() {
  const queryClient = useQueryClient()

  return createMutation<void, Error, void>(() => ({
    mutationFn: () => userService.deleteProfile(),
    onSuccess: () => {
      queryClient.clear()
      authState.token = null
      storage.remove(AUTH_KEYS.ACCESS_TOKEN)
      storage.remove(AUTH_KEYS.REFRESH_TOKEN)
    },
  }))
}

export function updateUserMutation() {
  const queryClient = useQueryClient()

  return createMutation<
    UserResponse,
    Error,
    { id: string; data: UpdateUserRequest }
  >(() => ({
    mutationFn: ({ id, data }) => userService.updateUser(id, data),
    onSuccess: updated => {
      queryClient.invalidateQueries({ queryKey: USERS_QUERY_KEY })
      queryClient.invalidateQueries({ queryKey: ['user', updated.id] })
    },
  }))
}

export function deleteUserMutation() {
  const queryClient = useQueryClient()

  return createMutation<void, Error, string>(() => ({
    mutationFn: id => userService.deleteUser(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: USERS_QUERY_KEY })
      queryClient.invalidateQueries({ queryKey: TASKS_QUERY_KEY })
    },
  }))
}

export function getUsersQuery(paramsGetter: () => GetUsersParams = () => ({})) {
  return createQuery(() => {
    const params = paramsGetter()
    return {
      queryKey: [...USERS_QUERY_KEY, params],
      queryFn: () => userService.getUsers(params),
      enabled: !!authState.token,
      staleTime: 1000 * 60 * 2,
      retry: 1,
    }
  })
}
