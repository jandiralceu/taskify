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
	avatarUrl?: string;
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

export interface TaskResponse {
	id: string;
	title: string;
	description: string;
	status: TaskStatus;
	priority: TaskPriority;
	isBlocked: boolean;
	createdBy: string;
	assignedTo?: string;
	dueDate?: string;
	completedAt?: string;
	estimatedHours?: number;
	actualHours?: number;
	isArchived: boolean;
	createdAt: string;
	updatedAt: string;
}

export interface CreateTaskRequest {
	title: string;
	description?: string;
	status?: TaskStatus;
	priority?: TaskPriority;
	isBlocked?: boolean;
	assignedTo?: string;
	dueDate?: string;
	estimatedHours?: number;
}

export interface UpdateTaskRequest {
	title?: string;
	description?: string;
	status?: TaskStatus;
	priority?: TaskPriority;
	isBlocked?: boolean;
	assignedTo?: string;
	dueDate?: string;
	estimatedHours?: number;
	actualHours?: number;
	isArchived?: boolean;
}

export interface PaginatedResponse<T> {
	data: T[];
	total: number;
	page: number;
	pageSize: number;
	totalPages: number;
}

