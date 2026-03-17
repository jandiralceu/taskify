<script lang="ts">
  import { page } from '$app/state'
  import { goto } from '$app/navigation'
  import { resolve } from '$app/paths'
  import { resolveAvatarUrl } from '$lib/utils/avatar'
  import {
    ArrowLeft,
    LoaderCircle,
    Flag,
    Clock,
    User,
    MessageCircle,
    Paperclip,
    Trash2,
    Send,
    Upload,
    Download,
    FileText,
    Image,
    FileArchive,
  } from '@lucide/svelte'
  import {
    getTaskQuery,
    updateTaskMutation,
    addNoteMutation,
    deleteNoteMutation,
    addAttachmentMutation,
    deleteAttachmentMutation,
  } from '$lib/state/tasks.svelte'
  import Button from '$lib/components/Button.svelte'
  import { getUsersQuery } from '$lib/state/user.svelte'
  import { toaster } from '$lib/state/toast.svelte'
  import Input from '$lib/components/Input.svelte'
  import type { TaskStatus, TaskPriority } from '$lib/api/types'

  const taskId = $derived(page.params.id ?? '')

  const taskQuery = getTaskQuery(() => taskId)
  const updateTask = updateTaskMutation()
  const addNote = addNoteMutation()
  const deleteNote = deleteNoteMutation()
  const addAttachment = addAttachmentMutation()
  const deleteAttachment = deleteAttachmentMutation()

  const usersQuery = getUsersQuery(() => ({ limit: 100 }))

  // Editable fields
  let title = $state('')
  let description = $state('')
  let status = $state<TaskStatus>('pending')
  let priority = $state<TaskPriority>('medium')
  let assignedTo = $state('')
  let dueDate = $state('')
  let estimatedHours = $state<number | ''>('')
  let actualHours = $state<number | ''>('')
  let isBlocked = $state(false)
  let isArchived = $state(false)

  // Comments state
  let newComment = $state('')

  // Attachments state
  let attachmentInput = $state<HTMLInputElement>(null!)

  function fileIcon(mimeType: string) {
    if (mimeType.startsWith('image/')) return Image
    if (mimeType === 'application/pdf' || mimeType.startsWith('text/')) return FileText
    if (mimeType.includes('zip') || mimeType.includes('tar') || mimeType.includes('gzip')) return FileArchive
    return Paperclip
  }

  function formatFileSize(bytes: number) {
    if (bytes < 1024) return `${bytes} B`
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
    return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
  }

  async function handleAddComment() {
    const content = newComment.trim()
    if (!content || !taskId) return
    try {
      await addNote.mutateAsync({ taskId, content })
      newComment = ''
    } catch {
      toaster.error({ title: 'Error', description: 'Failed to add comment.' })
    }
  }

  async function handleDeleteComment(noteId: string) {
    try {
      await deleteNote.mutateAsync({ noteId, taskId })
    } catch {
      toaster.error({ title: 'Error', description: 'Failed to delete comment.' })
    }
  }

  async function onAttachmentSelected(e: Event) {
    const file = (e.target as HTMLInputElement).files?.[0]
    if (!file || !taskId) return
    try {
      await addAttachment.mutateAsync({ taskId, file })
      toaster.success({ title: 'Uploaded', description: `${file.name} attached successfully.` })
    } catch {
      toaster.error({ title: 'Upload Failed', description: 'Could not attach the file.' })
    } finally {
      attachmentInput.value = ''
    }
  }

  async function handleDeleteAttachment(attachmentId: string) {
    try {
      await deleteAttachment.mutateAsync({ attachmentId, taskId })
    } catch {
      toaster.error({ title: 'Error', description: 'Failed to remove attachment.' })
    }
  }

  $effect(() => {
    if (taskQuery.data) {
      const t = taskQuery.data
      title = t.title
      description = t.description ?? ''
      status = t.status
      priority = t.priority
      assignedTo = t.assignedTo ?? ''
      dueDate = t.dueDate ? t.dueDate.split('T')[0] : ''
      estimatedHours = t.estimatedHours ?? ''
      actualHours = t.actualHours ?? ''
      isBlocked = t.isBlocked
      isArchived = t.isArchived
    }
  })

  const statusOptions: { value: TaskStatus; label: string }[] = [
    { value: 'pending', label: 'Pending' },
    { value: 'in_progress', label: 'In Progress' },
    { value: 'completed', label: 'Completed' },
    { value: 'cancelled', label: 'Cancelled' },
  ]

  const priorityOptions: { value: TaskPriority; label: string }[] = [
    { value: 'low', label: 'Low' },
    { value: 'medium', label: 'Medium' },
    { value: 'high', label: 'High' },
    { value: 'critical', label: 'Critical' },
  ]

  function formatDate(dateStr: string) {
    return new Intl.DateTimeFormat('en-GB', {
      day: '2-digit',
      month: 'long',
      year: 'numeric',
    }).format(new Date(dateStr))
  }

  async function handleSubmit(e: SubmitEvent) {
    e.preventDefault()

    try {
      await updateTask.mutateAsync({
        id: taskId,
        data: {
          title,
          description: description || undefined,
          status,
          priority,
          assignedTo: assignedTo || undefined,
          dueDate: dueDate ? `${dueDate}T00:00:00Z` : undefined,
          estimatedHours:
            estimatedHours !== '' ? Number(estimatedHours) : undefined,
          actualHours: actualHours !== '' ? Number(actualHours) : undefined,
          isBlocked,
          isArchived,
        },
      })

      toaster.success({
        title: 'Task Updated',
        description: `"${title}" has been updated successfully.`,
      })

      goto(resolve('/'))
    } catch {
      toaster.error({
        title: 'Update Failed',
        description: 'Could not update the task. Please try again.',
      })
    }
  }
</script>

<svelte:head>
  <title>Edit Task - Taskify</title>
</svelte:head>

<div class="flex h-full flex-col pt-8">
  <header class="px-8 pb-8">
    <button
      onclick={() => goto(resolve('/'))}
      class="mb-6 flex items-center gap-2 text-sm text-surface-500 transition-colors hover:text-surface-800"
    >
      <ArrowLeft size={16} />
      Back to Board
    </button>

    <h2 class="text-4xl leading-tight tracking-tight text-surface-900">
      <span class="font-light">Edit</span>
      <span class="font-normal"> Task</span>
    </h2>
  </header>

  <div class="flex-1 overflow-y-auto px-8 pb-8">
    {#if taskQuery.isPending}
      <div
        class="flex flex-col items-center justify-center py-24 text-surface-400"
      >
        <LoaderCircle size={28} class="mb-3 animate-spin" />
        <span class="text-sm font-medium">Loading task...</span>
      </div>
    {:else if taskQuery.isError}
      <div
        class="flex flex-col items-center justify-center py-24 text-rose-500"
      >
        <p class="text-sm font-medium">Failed to load task.</p>
      </div>
    {:else if taskQuery.data}
      {@const task = taskQuery.data}
      <div class="max-w-2xl space-y-6">
        <!-- Read-only meta card -->
        <div
          class="grid grid-cols-1 gap-4 rounded-2xl border border-surface-200 bg-white p-6 sm:grid-cols-3"
        >
          <!-- Created by -->
          <div class="space-y-1">
            <p
              class="text-xs font-semibold tracking-wide text-surface-400 uppercase"
            >
              Created by
            </p>
            <div class="flex items-center gap-2">
              {#if task.assignee.avatarUrl}
                <img
                  src={resolveAvatarUrl(task.assignee.avatarUrl) ?? ''}
                  alt=""
                  class="size-6 rounded-full object-cover"
                />
              {:else}
                <div
                  class="flex size-6 items-center justify-center rounded-full bg-indigo-100"
                >
                  <span class="text-[9px] font-bold text-indigo-700">
                    {task.assignee.firstName[0]}{task.assignee.lastName[0]}
                  </span>
                </div>
              {/if}
              <span class="text-sm font-medium text-surface-700">
                {task.assignee.firstName}
                {task.assignee.lastName}
              </span>
            </div>
          </div>

          <!-- Created at -->
          <div class="space-y-1">
            <p
              class="text-xs font-semibold tracking-wide text-surface-400 uppercase"
            >
              Created
            </p>
            <p class="text-sm text-surface-700">{formatDate(task.createdAt)}</p>
          </div>

          <!-- Updated at -->
          <div class="space-y-1">
            <p
              class="text-xs font-semibold tracking-wide text-surface-400 uppercase"
            >
              Last updated
            </p>
            <p class="text-sm text-surface-700">{formatDate(task.updatedAt)}</p>
          </div>
        </div>

        <!-- Edit form -->
        <form
          onsubmit={handleSubmit}
          class="space-y-5 rounded-2xl border border-surface-200 bg-white p-6"
        >
          <h3 class="text-sm font-semibold text-surface-900">Task Details</h3>

          <!-- Title -->
          <Input
            id="title"
            label="Title"
            placeholder="Task title"
            bind:value={title}
            required
          />

          <!-- Description -->
          <div class="space-y-1">
            <label
              for="description"
              class="block text-sm font-medium text-surface-700"
            >
              Description
            </label>
            <textarea
              id="description"
              rows={4}
              placeholder="Task description..."
              bind:value={description}
              class="block w-full resize-none rounded-xl border border-surface-300 bg-surface-50 px-4 py-3 text-sm text-surface-900 placeholder-surface-500 transition-all focus:border-primary-500 focus:ring-2 focus:ring-primary-500/20 focus:outline-none"
            ></textarea>
          </div>

          <!-- Status + Priority -->
          <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <div class="space-y-1">
              <label
                for="status"
                class="block text-sm font-medium text-surface-700">Status</label
              >
              <select
                id="status"
                bind:value={status}
                class="block h-12 w-full rounded-xl border border-surface-300 bg-surface-50 px-4 text-sm text-surface-900 transition-all focus:border-primary-500 focus:ring-2 focus:ring-primary-500/20 focus:outline-none"
              >
                {#each statusOptions as opt (opt.value)}
                  <option value={opt.value}>{opt.label}</option>
                {/each}
              </select>
            </div>

            <div class="space-y-1">
              <label
                for="priority"
                class="block text-sm font-medium text-surface-700"
                >Priority</label
              >
              <select
                id="priority"
                bind:value={priority}
                class="block h-12 w-full rounded-xl border border-surface-300 bg-surface-50 px-4 text-sm text-surface-900 transition-all focus:border-primary-500 focus:ring-2 focus:ring-primary-500/20 focus:outline-none"
              >
                {#each priorityOptions as opt (opt.value)}
                  <option value={opt.value}>{opt.label}</option>
                {/each}
              </select>
            </div>
          </div>

          <!-- Assigned To -->
          <div class="space-y-1">
            <label
              for="assignedTo"
              class="block text-sm font-medium text-surface-700"
            >
              <span class="flex items-center gap-1.5"
                ><User size={14} /> Assigned To</span
              >
            </label>
            <select
              id="assignedTo"
              bind:value={assignedTo}
              class="block h-12 w-full rounded-xl border border-surface-300 bg-surface-50 px-4 text-sm text-surface-900 transition-all focus:border-primary-500 focus:ring-2 focus:ring-primary-500/20 focus:outline-none"
            >
              <option value="">Unassigned</option>
              {#if usersQuery.data?.data}
                {#each usersQuery.data.data as user (user.id)}
                  <option value={user.id}
                    >{user.firstName} {user.lastName}</option
                  >
                {/each}
              {/if}
            </select>
          </div>

          <!-- Due Date + Estimated Hours -->
          <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <div class="space-y-1">
              <label
                for="dueDate"
                class="block text-sm font-medium text-surface-700"
              >
                <span class="flex items-center gap-1.5"
                  ><Flag size={14} /> Due Date</span
                >
              </label>
              <input
                id="dueDate"
                type="date"
                bind:value={dueDate}
                class="block h-12 w-full rounded-xl border border-surface-300 bg-surface-50 px-4 text-sm text-surface-900 transition-all focus:border-primary-500 focus:ring-2 focus:ring-primary-500/20 focus:outline-none"
              />
            </div>

            <div class="space-y-1">
              <label
                for="estimatedHours"
                class="block text-sm font-medium text-surface-700"
              >
                <span class="flex items-center gap-1.5"
                  ><Clock size={14} /> Estimated Hours</span
                >
              </label>
              <input
                id="estimatedHours"
                type="number"
                min="0"
                step="0.5"
                placeholder="0"
                bind:value={estimatedHours}
                class="block h-12 w-full rounded-xl border border-surface-300 bg-surface-50 px-4 text-sm text-surface-900 transition-all focus:border-primary-500 focus:ring-2 focus:ring-primary-500/20 focus:outline-none"
              />
            </div>
          </div>

          <!-- Actual Hours -->
          <div class="space-y-1">
            <label
              for="actualHours"
              class="block text-sm font-medium text-surface-700"
            >
              <span class="flex items-center gap-1.5"
                ><Clock size={14} /> Actual Hours</span
              >
            </label>
            <input
              id="actualHours"
              type="number"
              min="0"
              step="0.5"
              placeholder="0"
              bind:value={actualHours}
              class="block h-12 w-full rounded-xl border border-surface-300 bg-surface-50 px-4 text-sm text-surface-900 transition-all focus:border-primary-500 focus:ring-2 focus:ring-primary-500/20 focus:outline-none"
            />
          </div>

          <hr class="border-surface-100" />

          <!-- Toggles: Blocked + Archived -->
          <div class="space-y-3">
            <!-- Blocked -->
            <div
              class="flex items-center justify-between rounded-xl border border-surface-200 bg-surface-50 p-4"
            >
              <div>
                <p class="text-sm font-medium text-surface-800">Blocked</p>
                <p class="mt-0.5 text-xs text-surface-500">
                  Blocked tasks cannot be moved on the board
                </p>
              </div>
              <button
                type="button"
                aria-label="Toggle blocked status"
                role="switch"
                aria-checked={isBlocked}
                onclick={() => (isBlocked = !isBlocked)}
                class="relative inline-flex h-6 w-11 shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 focus:outline-none {isBlocked
                  ? 'bg-rose-500'
                  : 'bg-surface-300'}"
              >
                <span
                  class="pointer-events-none inline-block size-5 transform rounded-full bg-white shadow-sm transition duration-200 {isBlocked
                    ? 'translate-x-5'
                    : 'translate-x-0'}"
                ></span>
              </button>
            </div>

            <!-- Archived -->
            <div
              class="flex items-center justify-between rounded-xl border border-surface-200 bg-surface-50 p-4"
            >
              <div>
                <p class="text-sm font-medium text-surface-800">Archived</p>
                <p class="mt-0.5 text-xs text-surface-500">
                  Archived tasks are hidden from the board
                </p>
              </div>
              <button
                type="button"
                aria-label="Toggle archived status"
                role="switch"
                aria-checked={isArchived}
                onclick={() => (isArchived = !isArchived)}
                class="relative inline-flex h-6 w-11 shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 focus:outline-none {isArchived
                  ? 'bg-primary-500'
                  : 'bg-surface-300'}"
              >
                <span
                  class="pointer-events-none inline-block size-5 transform rounded-full bg-white shadow-sm transition duration-200 {isArchived
                    ? 'translate-x-5'
                    : 'translate-x-0'}"
                ></span>
              </button>
            </div>
          </div>

          <!-- Footer -->
          <div class="flex items-center justify-end gap-3 pt-2">
            <Button
              variant="ghost"
              onclick={() => goto(resolve('/'))}
            >
              Cancel
            </Button>
            <Button
              type="submit"
              loading={updateTask.isPending}
              loadingText="Saving..."
            >
              Save Changes
            </Button>
          </div>
        </form>

        <!-- Comments -->
        <div class="max-w-2xl space-y-4 rounded-2xl border border-surface-200 bg-white p-6">
          <div class="flex items-center gap-2">
            <MessageCircle size={16} class="text-surface-500" />
            <h3 class="text-sm font-semibold text-surface-900">Comments</h3>
            <span class="rounded-full bg-surface-100 px-2 py-0.5 text-xs font-medium text-surface-500">
              {taskQuery.data?.notes?.length ?? taskQuery.data?.notesCount ?? 0}
            </span>
          </div>

          <!-- Add comment -->
          <div class="space-y-2">
            <textarea
              rows={3}
              placeholder="Write a comment..."
              bind:value={newComment}
              class="block w-full resize-none rounded-xl border border-surface-300 bg-surface-50 px-4 py-3 text-sm text-surface-900 placeholder-surface-400 transition-all focus:border-primary-500 focus:ring-2 focus:ring-primary-500/20 focus:outline-none"
            ></textarea>
            <div class="flex justify-end">
              <Button
                onclick={handleAddComment}
                disabled={!newComment.trim()}
                loading={addNote.isPending}
                loadingText="Posting..."
              >
                <Send size={14} />
                Post Comment
              </Button>
            </div>
          </div>

          <!-- Comment list -->
          {#if taskQuery.data?.notes && taskQuery.data.notes.length > 0}
            <div class="space-y-3 border-t border-surface-100 pt-4">
              {#each taskQuery.data.notes as note (note.id)}
                <div class="group flex gap-3">
                  <div class="mt-0.5 flex size-7 shrink-0 items-center justify-center rounded-full bg-primary-50 text-[10px] font-bold text-primary-600">
                    ?
                  </div>
                  <div class="min-w-0 flex-1 rounded-xl border border-surface-100 bg-surface-50 px-4 py-3">
                    <p class="text-sm leading-relaxed text-surface-700">{note.content}</p>
                    <div class="mt-1.5 flex items-center justify-between">
                      <p class="text-[11px] text-surface-400">
                        {new Intl.DateTimeFormat('en-GB', { day: 'numeric', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit' }).format(new Date(note.createdAt))}
                      </p>
                      <button
                        type="button"
                        onclick={() => handleDeleteComment(note.id)}
                        class="invisible rounded-md p-1 text-surface-300 transition-colors hover:bg-rose-50 hover:text-rose-500 group-hover:visible"
                        title="Delete comment"
                      >
                        <Trash2 size={13} />
                      </button>
                    </div>
                  </div>
                </div>
              {/each}
            </div>
          {:else}
            <p class="text-sm text-surface-400">No comments yet.</p>
          {/if}
        </div>

        <!-- Attachments -->
        <div class="max-w-2xl space-y-4 rounded-2xl border border-surface-200 bg-white p-6">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <Paperclip size={16} class="text-surface-500" />
              <h3 class="text-sm font-semibold text-surface-900">Attachments</h3>
              <span class="rounded-full bg-surface-100 px-2 py-0.5 text-xs font-medium text-surface-500">
                {taskQuery.data?.attachments?.length ?? taskQuery.data?.attachmentsCount ?? 0}
              </span>
            </div>
            <div>
              <input
                bind:this={attachmentInput}
                type="file"
                class="hidden"
                onchange={onAttachmentSelected}
              />
              <Button
                variant="ghost"
                onclick={() => attachmentInput.click()}
                loading={addAttachment.isPending}
                loadingText="Uploading..."
                class="!px-3"
              >
                <Upload size={14} />
                Upload File
              </Button>
            </div>
          </div>

          {#if taskQuery.data?.attachments && taskQuery.data.attachments.length > 0}
            <div class="space-y-2">
              {#each taskQuery.data.attachments as att (att.id)}
                {@const Icon = fileIcon(att.mimeType)}
                <div class="group flex items-center gap-3 rounded-xl border border-surface-100 bg-surface-50 px-4 py-3">
                  <Icon size={18} class="shrink-0 text-surface-400" />
                  <div class="min-w-0 flex-1">
                    <p class="truncate text-sm font-medium text-surface-800">{att.fileName}</p>
                    <p class="text-xs text-surface-400">{formatFileSize(att.fileSize)}</p>
                  </div>
                  <a
                    href={resolve(att.filePath)}
                    download={att.fileName}
                    class="shrink-0 rounded-lg p-1.5 text-surface-400 transition-colors hover:bg-surface-100 hover:text-surface-700"
                    title="Download"
                  >
                    <Download size={15} />
                  </a>
                  <button
                    type="button"
                    onclick={() => handleDeleteAttachment(att.id)}
                    class="shrink-0 rounded-lg p-1.5 text-surface-300 transition-colors hover:bg-rose-50 hover:text-rose-500"
                    title="Remove attachment"
                  >
                    <Trash2 size={14} />
                  </button>
                </div>
              {/each}
            </div>
          {:else}
            <p class="text-sm text-surface-400">No attachments yet.</p>
          {/if}
        </div>
      </div>
    {/if}
  </div>
</div>
