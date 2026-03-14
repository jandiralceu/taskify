import { api } from './client';
import type { CreateUserRequest, SignInRequest, SignInResponse } from './types';

export const authService = {
	async signup(data: CreateUserRequest): Promise<void> {
		await api.post('auth/register', { json: data });
	},

	async signin(data: SignInRequest): Promise<SignInResponse> {
		return api.post('auth/signin', { json: data }).json();
	}
};
