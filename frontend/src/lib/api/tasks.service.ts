import { api } from './client'
import type {
  CreateTaskRequest,
  TaskResponse,
  TaskNoteResponse,
  TaskAttachmentResponse,
  UpdateTaskRequest,
  TaskStatus,
  TaskPriority,
} from './types'

export interface GetTasksParams {
  status?: TaskStatus
  priority?: TaskPriority
  search?: string
  sort?: string
  order?: 'asc' | 'desc'
  isBlocked?: boolean
}

class TasksService {
  async getTasks(params: GetTasksParams = {}) {
    const searchParams = new URLSearchParams()
    if (params.status) searchParams.set('status', params.status)
    if (params.priority) searchParams.set('priority', params.priority)
    if (params.search) searchParams.set('search', params.search)
    if (params.sort) searchParams.set('sort', params.sort)
    if (params.order) searchParams.set('order', params.order)
    if (params.isBlocked !== undefined)
      searchParams.set('isBlocked', String(params.isBlocked))

    return api.get('tasks', { searchParams }).json<TaskResponse[]>()
  }

  async getTask(id: string) {
    return api.get(`tasks/${id}`).json<TaskResponse>()
  }

  async createTask(data: CreateTaskRequest) {
    return api.post('tasks', { json: data }).json<TaskResponse>()
  }

  async updateTask(id: string, data: UpdateTaskRequest) {
    return api.patch(`tasks/${id}`, { json: data }).json<TaskResponse>()
  }

  async deleteTask(id: string) {
    return api.delete(`tasks/${id}`).json<void>()
  }

  // Notes
  async addNote(taskId: string, content: string) {
    return api.post(`tasks/${taskId}/notes`, { json: { content } }).json<TaskNoteResponse>()
  }

  async deleteNote(noteId: string) {
    return api.delete(`tasks/notes/${noteId}`).json<void>()
  }

  // Attachments
  async addAttachment(taskId: string, file: File) {
    const formData = new FormData()
    formData.append('file', file)
    return api.post(`tasks/${taskId}/attachments`, { body: formData }).json<TaskAttachmentResponse>()
  }

  async deleteAttachment(attachmentId: string) {
    return api.delete(`tasks/attachments/${attachmentId}`).json<void>()
  }
}

export const tasksService = new TasksService()
