import { api } from './client';
import type { UserResponse, PermissionsResponse } from './types';

export const userService = {
	async getProfile(): Promise<UserResponse> {
		return api.get('users/profile').json();
	},
	async getPermissions(): Promise<PermissionsResponse> {
		return api.get('users/permissions').json();
	},
	async getUserById(id: string): Promise<UserResponse> {
		return api.get(`users/${id}`).json();
	}
};
