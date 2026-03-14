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

export interface SignOutRequest {
	refreshToken: string;
}

export interface RefreshTokenRequest {
	refreshToken: string;
}

export interface RefreshTokenResponse {
	accessToken: string;
	refreshToken: string;
}

export type TaskStatus = 'pending' | 'in_progress' | 'completed' | 'cancelled';
export type TaskPriority = 'low' | 'medium' | 'high' | 'critical';

