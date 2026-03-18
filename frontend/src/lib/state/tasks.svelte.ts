import {
  createQuery,
  createMutation,
  useQueryClient,
} from '@tanstack/svelte-query'
import { tasksService, type GetTasksParams } from '$lib/api/tasks.service'
import type {
  CreateTaskRequest,
  UpdateTaskRequest,
  TaskResponse,
  TaskNoteResponse,
  TaskAttachmentResponse,
} from '$lib/api/types'

export const TASKS_QUERY_KEY = ['tasks']

export function getTasksQuery(paramsGetter: () => GetTasksParams = () => ({})) {
  return createQuery(() => {
    const params = paramsGetter()
    return {
      queryKey: [...TASKS_QUERY_KEY, 'list', params],
      queryFn: () => tasksService.getTasks(params),
      staleTime: 0,
      gcTime: 0,
    }
  })
}

export function getArchivedTasksQuery(paramsGetter: () => GetTasksParams = () => ({})) {
  return createQuery(() => {
    const params = paramsGetter()
    return {
      queryKey: [...TASKS_QUERY_KEY, 'archived', params],
      queryFn: () => tasksService.getArchivedTasks(params),
      staleTime: 0,
      gcTime: 0,
    }
  })
}

export function getTaskQuery(idGetter: () => string) {
  return createQuery(() => {
    const id = idGetter()
    return {
      queryKey: [...TASKS_QUERY_KEY, id],
      queryFn: () => tasksService.getTask(id),
      enabled: !!id,
      staleTime: 0,
      gcTime: 0,
    }
  })
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
 * const updateTask = updateTaskMutation();
 * await updateTask.mutateAsync({ id: task.id, data: { status: 'in_progress' } });
 */
export function updateTaskMutation() {
  const queryClient = useQueryClient()

  return createMutation<
    TaskResponse,
    Error,
    { id: string; data: UpdateTaskRequest }
  >(() => ({
    mutationFn: ({ id, data }) => tasksService.updateTask(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: TASKS_QUERY_KEY })
    },
  }))
}

export function createTaskMutation() {
  const queryClient = useQueryClient()

  return createMutation<TaskResponse, Error, CreateTaskRequest>(() => ({
    mutationFn: data => tasksService.createTask(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: TASKS_QUERY_KEY })
    },
  }))
}

export function deleteTaskMutation() {
  const queryClient = useQueryClient()

  return createMutation<void, Error, string>(() => ({
    mutationFn: id => tasksService.deleteTask(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: TASKS_QUERY_KEY })
    },
  }))
}

export function addNoteMutation() {
  const queryClient = useQueryClient()

  return createMutation<TaskNoteResponse, Error, { taskId: string; content: string }>(() => ({
    mutationFn: ({ taskId, content }) => tasksService.addNote(taskId, content),
    onSuccess: (_, { taskId }) => {
      queryClient.invalidateQueries({ queryKey: [...TASKS_QUERY_KEY, taskId] })
    },
  }))
}

export function deleteNoteMutation() {
  const queryClient = useQueryClient()

  return createMutation<void, Error, { noteId: string; taskId: string }>(() => ({
    mutationFn: ({ noteId }) => tasksService.deleteNote(noteId),
    onSuccess: (_, { taskId }) => {
      queryClient.invalidateQueries({ queryKey: [...TASKS_QUERY_KEY, taskId] })
    },
  }))
}

export function addAttachmentMutation() {
  const queryClient = useQueryClient()

  return createMutation<TaskAttachmentResponse, Error, { taskId: string; file: File }>(() => ({
    mutationFn: ({ taskId, file }) => tasksService.addAttachment(taskId, file),
    onSuccess: (_, { taskId }) => {
      queryClient.invalidateQueries({ queryKey: [...TASKS_QUERY_KEY, taskId] })
    },
  }))
}

export function deleteAttachmentMutation() {
  const queryClient = useQueryClient()

  return createMutation<void, Error, { attachmentId: string; taskId: string }>(() => ({
    mutationFn: ({ attachmentId }) => tasksService.deleteAttachment(attachmentId),
    onSuccess: (_, { taskId }) => {
      queryClient.invalidateQueries({ queryKey: [...TASKS_QUERY_KEY, taskId] })
    },
  }))
}
