import ky from 'ky';
import { PUBLIC_API_URL } from '$env/static/public';
import { storage, AUTH_KEYS } from '$lib/utils/storage';

const PUBLIC_ENDPOINTS = ['auth/signin', 'auth/register'];

export const api = ky.create({
	prefixUrl: PUBLIC_API_URL,
	hooks: {
		beforeRequest: [
			(request) => {
				const isPublic = PUBLIC_ENDPOINTS.some((endpoint) => request.url.includes(endpoint));
				if (isPublic) return;

				const token = storage.get(AUTH_KEYS.ACCESS_TOKEN);
				if (token) {
					request.headers.set('Authorization', `Bearer ${token}`);
				}
			}
		]
	}
});
