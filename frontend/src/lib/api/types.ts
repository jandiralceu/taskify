export interface CreateUserRequest {
	firstName: string;
	lastName: string;
	email: string;
	password: string;
	role: 'admin' | 'employee';
}

export interface UserResponse {
	id: string;
	firstName: string;
	lastName: string;
	email: string;
	role: string;
	createdAt: string;
}

export interface SignInRequest {
	email: string;
	password: string;
}

export interface SignInResponse {
	accessToken: string;
	refreshToken: string;
}
