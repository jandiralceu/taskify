import { api } from './client';
import type { CreateUserRequest, SignInRequest, SignInResponse, SignOutRequest, RefreshTokenRequest, RefreshTokenResponse } from './types';

export const authService = {
	async signup(data: CreateUserRequest): Promise<void> {
		await api.post('auth/register', { json: data });
	},

	async signin(data: SignInRequest): Promise<SignInResponse> {
		return api.post('auth/signin', { json: data }).json();
	},

	async signout(data: SignOutRequest): Promise<void> {
		await api.post('auth/signout', { json: data });
	},

	async refresh(data: RefreshTokenRequest): Promise<RefreshTokenResponse> {
		return api.post('auth/refresh', { json: data }).json();
	}
};
