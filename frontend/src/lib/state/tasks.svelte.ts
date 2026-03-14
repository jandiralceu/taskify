import { createQuery, createMutation, useQueryClient } from '@tanstack/svelte-query';
import { tasksService, type GetTasksParams } from '$lib/api/tasks.service';
import type { CreateTaskRequest, TaskResponse } from '$lib/api/types';

export const TASKS_QUERY_KEY = ['tasks'];

export function createTasksQuery(params: GetTasksParams = {}) {
	return createQuery(() => ({
		queryKey: [...TASKS_QUERY_KEY, params],
		queryFn: () => tasksService.getTasks(params)
	}));
}

export function createCreateTaskMutation() {
	const queryClient = useQueryClient();
	
	return createMutation<TaskResponse, Error, CreateTaskRequest>(() => ({
		mutationFn: (data) => tasksService.createTask(data),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: TASKS_QUERY_KEY });
		}
	}));
}
