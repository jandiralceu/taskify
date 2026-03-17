import { api } from './client';
import type { UserResponse, PermissionsResponse, PaginatedResponse, GetUsersParams, UpdateUserRequest, ChangePasswordRequest } from './types';

export const userService = {
	async getProfile(): Promise<UserResponse> {
		return api.get('users/profile').json();
	},
	async getPermissions(): Promise<PermissionsResponse> {
		return api.get('users/permissions').json();
	},
	async getUserById(id: string): Promise<UserResponse> {
		return api.get(`users/${id}`).json();
	},
	async getUsers(params: GetUsersParams = {}): Promise<PaginatedResponse<UserResponse>> {
		const searchParams = new URLSearchParams();
		if (params.page) searchParams.set('page', String(params.page));
		if (params.limit) searchParams.set('limit', String(params.limit));
		if (params.firstName) searchParams.set('firstName', params.firstName);
		if (params.lastName) searchParams.set('lastName', params.lastName);
		if (params.email) searchParams.set('email', params.email);
		if (params.role) searchParams.set('role', params.role);
		if (params.sort) searchParams.set('sort', params.sort);
		if (params.order) searchParams.set('order', params.order);
		return api.get('users', { searchParams }).json();
	},
	async updateProfile(data: UpdateUserRequest): Promise<UserResponse> {
		return api.patch('users/profile', { json: data }).json();
	},
	async changePassword(data: ChangePasswordRequest): Promise<void> {
		await api.patch('users/change-password', { json: data });
	},
	async deleteProfile(): Promise<void> {
		await api.delete('users/profile');
	},
	async updateUser(id: string, data: UpdateUserRequest): Promise<UserResponse> {
		return api.patch(`users/${id}`, { json: data }).json();
	},
	async deleteUser(id: string): Promise<void> {
		await api.delete(`users/${id}`);
	}
};
