import ky from 'ky';
import { PUBLIC_API_URL } from '$env/static/public';
import { storage, AUTH_KEYS } from '$lib/utils/storage';

export const api = ky.create({
	prefixUrl: PUBLIC_API_URL,
	timeout: 10000,
	hooks: {
		beforeRequest: [
			(request) => {
				const token = storage.get(AUTH_KEYS.ACCESS_TOKEN);
				if (token) {
					request.headers.set('Authorization', `Bearer ${token}`);
				}
			}
		]
	}
});
