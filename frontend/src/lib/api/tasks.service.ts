import { api } from './client';
import type {
	CreateTaskRequest,
	TaskResponse,
	UpdateTaskRequest,
	TaskStatus,
	TaskPriority
} from './types';

export interface GetTasksParams {
	status?: TaskStatus;
	priority?: TaskPriority;
	search?: string;
	sort?: string;
	order?: 'asc' | 'desc';
}

class TasksService {
	async getTasks(params: GetTasksParams = {}) {
		const searchParams = new URLSearchParams();
		if (params.status) searchParams.set('status', params.status);
		if (params.priority) searchParams.set('priority', params.priority);
		if (params.search) searchParams.set('search', params.search);
		if (params.sort) searchParams.set('sort', params.sort);
		if (params.order) searchParams.set('order', params.order);

		return api.get('tasks', { searchParams }).json<TaskResponse[]>();
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
