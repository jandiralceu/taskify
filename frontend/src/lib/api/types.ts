export type UserRole = 'admin' | 'employee'

export interface CreateUserRequest {
  firstName: string
  lastName: string
  email: string
  password: string
  role: UserRole
}

export interface UserResponse {
  id: string
  firstName: string
  lastName: string
  email: string
  role: UserRole
  isActive: boolean
  avatarUrl?: string
  createdAt: string
}

export interface PaginatedResponse<T> {
  data: T[]
  total: number
  page: number
  limit: number
  totalPages: number
}

export interface GetUsersParams {
  page?: number
  limit?: number
  firstName?: string
  lastName?: string
  email?: string
  role?: UserRole
  sort?: string
  order?: 'asc' | 'desc'
}

export interface PermissionsResponse {
  role: UserRole
  permissions: {
    tasks: string[]
    users: string[]
    admin_area: boolean
  }
}

export interface SignInRequest {
  email: string
  password: string
}

export interface SignInResponse {
  accessToken: string
  refreshToken: string
}

export interface SignOutRequest {
  refreshToken: string
}

export interface RefreshTokenRequest {
  refreshToken: string
}

export interface RefreshTokenResponse {
  accessToken: string
  refreshToken: string
}

export type TaskStatus = 'pending' | 'in_progress' | 'completed' | 'cancelled'
export type TaskPriority = 'low' | 'medium' | 'high' | 'critical'

export interface TaskNoteResponse {
  id: string
  taskId: string
  userId: string
  content: string
  user?: UserResponse
  createdAt: string
  updatedAt: string
}

export interface TaskAttachmentResponse {
  id: string
  taskId: string
  userId: string
  fileName: string
  fileSize: number
  mimeType: string
  filePath: string
  user?: UserResponse
  createdAt: string
}

export interface TaskResponse {
  id: string
  title: string
  description: string
  status: TaskStatus
  priority: TaskPriority
  isBlocked: boolean
  createdBy: string
  assignedTo?: string
  assignee: UserResponse
  dueDate?: string
  completedAt?: string
  estimatedHours?: number
  actualHours?: number
  isArchived: boolean
  notesCount: number
  attachmentsCount: number
  notes?: TaskNoteResponse[]
  attachments?: TaskAttachmentResponse[]
  createdAt: string
  updatedAt: string
}

export interface CreateTaskRequest {
  title: string
  description?: string
  status?: TaskStatus
  priority?: TaskPriority
  isBlocked?: boolean
  assignedTo?: string
  dueDate?: string
  estimatedHours?: number
}

export interface UpdateUserRequest {
  firstName?: string
  lastName?: string
  isActive?: boolean
}

export interface ChangePasswordRequest {
  oldPassword: string
  newPassword: string
}

export interface UpdateTaskRequest {
  title?: string
  description?: string
  status?: TaskStatus
  priority?: TaskPriority
  isBlocked?: boolean
  assignedTo?: string
  dueDate?: string
  estimatedHours?: number
  actualHours?: number
  isArchived?: boolean
}
