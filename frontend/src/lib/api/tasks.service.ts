import { api } from './client';
import type {
	CreateTaskRequest,
	TaskResponse,
	UpdateTaskRequest,
	PaginatedResponse,
	TaskStatus,
	TaskPriority
} from './types';

export interface GetTasksParams {
	page?: number;
	pageSize?: number;
	status?: TaskStatus;
	priority?: TaskPriority;
	search?: string;
}

class TasksService {
	async getTasks(params: GetTasksParams = {}) {
		const searchParams = new URLSearchParams();
		if (params.page) searchParams.set('page', params.page.toString());
		if (params.pageSize) searchParams.set('page_size', params.pageSize.toString());
		if (params.status) searchParams.set('status', params.status);
		if (params.priority) searchParams.set('priority', params.priority);
		if (params.search) searchParams.set('search', params.search);

		return api.get('tasks', { searchParams }).json<PaginatedResponse<TaskResponse>>();
	}

	async createTask(data: CreateTaskRequest) {
		return api.post('tasks', { json: data }).json<TaskResponse>();
	}

	async updateTask(id: string, data: UpdateTaskRequest) {
		return api.patch(`tasks/${id}`, { json: data }).json<TaskResponse>();
	}

	async deleteTask(id: string) {
		return api.delete(`tasks/${id}`).json<void>();
	}
}

export const tasksService = new TasksService();
