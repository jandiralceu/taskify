import { createQuery, createMutation, useQueryClient } from '@tanstack/svelte-query';
import { tasksService, type GetTasksParams } from '$lib/api/tasks.service';
import type { CreateTaskRequest, UpdateTaskRequest, TaskResponse } from '$lib/api/types';

export const TASKS_QUERY_KEY = ['tasks'];

export function createTasksQuery(params: GetTasksParams = {}) {
	return createQuery(() => ({
		queryKey: [...TASKS_QUERY_KEY, params],
		queryFn: () => tasksService.getTasks(params)
	}));
}

/**
 * Mutation for updating an existing task.
 *
 * Used by the drag-and-drop system to persist status changes when a card
 * is dropped into a different kanban column. After a successful update,
 * the tasks cache is invalidated so the board reflects the new state
 * without requiring a manual refresh.
 *
 * @example
 * const updateTask = createUpdateTaskMutation();
 * await updateTask.mutateAsync({ id: task.id, data: { status: 'in_progress' } });
 */
export function createUpdateTaskMutation() {
	const queryClient = useQueryClient();

	return createMutation<TaskResponse, Error, { id: string; data: UpdateTaskRequest }>(() => ({
		mutationFn: ({ id, data }) => tasksService.updateTask(id, data),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: TASKS_QUERY_KEY });
		}
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
