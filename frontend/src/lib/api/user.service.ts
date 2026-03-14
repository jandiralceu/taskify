import { api } from './client';
import type { UserResponse } from './types';

export const userService = {
	async getProfile(): Promise<UserResponse> {
		return api.get('users/profile').json();
	}
};
