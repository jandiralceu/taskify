<script lang="ts">
  import {
    X,
    Calendar,
    Flag,
    Clock,
    User,
    ShieldBan,
    CheckCircle,
    AlertCircle,
    MessageCircle,
    Paperclip,
    FileText,
    Image,
    FileArchive,
    Download,
    LoaderCircle,
  } from '@lucide/svelte'
  import { Dialog, Portal } from '@skeletonlabs/skeleton-svelte'
  import type { TaskResponse } from '$lib/api/types'
  import { priorityConfig, statusConfig } from '$lib/utils/task'
  import { getTaskQuery } from '$lib/state/tasks.svelte'
  import { resolve } from '$app/paths'
  import { resolveAvatarUrl } from '$lib/utils/avatar'

  interface Props {
    task: TaskResponse | null
    isOpen: boolean
    onClose: () => void
  }

  let { task, isOpen, onClose }: Props = $props()

  // Fetch full task (with notes + attachments) whenever the drawer is open
  const taskId = $derived(isOpen && task ? task.id : '')
  const fullTaskQuery = getTaskQuery(() => taskId)
  const fullTask = $derived(fullTaskQuery.data ?? task)

  function handleOpenChange(e: { open: boolean }) {
    if (!e.open) onClose()
  }

  function formatDate(dateStr: string) {
    return new Intl.DateTimeFormat('en-GB', {
      day: 'numeric',
      month: 'long',
      year: 'numeric',
      timeZone: 'UTC',
    }).format(new Date(dateStr))
  }

  function formatDateTime(dateStr: string) {
    return new Intl.DateTimeFormat('en-GB', {
      day: 'numeric',
      month: 'long',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    }).format(new Date(dateStr))
  }

  function formatHours(hours: number) {
    if (hours === 1) return '1 hour'
    return `${hours} hours`
  }

  function formatFileSize(bytes: number) {
    if (bytes < 1024) return `${bytes} B`
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
    return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
  }

  function fileIcon(mimeType: string) {
    if (mimeType.startsWith('image/')) return Image
    if (mimeType === 'application/pdf' || mimeType.startsWith('text/')) return FileText
    if (mimeType.includes('zip') || mimeType.includes('tar') || mimeType.includes('gzip')) return FileArchive
    return Paperclip
  }

  let priority = $derived(fullTask ? priorityConfig[fullTask.priority] : null)
  let status = $derived(fullTask ? statusConfig[fullTask.status] : null)
</script>

<Dialog
  open={isOpen}
  onOpenChange={handleOpenChange}
  closeOnInteractOutside={true}
>
  <Portal>
    <Dialog.Backdrop
      class="fixed inset-0 z-50 bg-surface-950/30 backdrop-blur-[2px] transition-opacity"
    />
    <Dialog.Positioner
      class="fixed inset-y-0 right-0 z-50 flex items-stretch justify-end"
    >
      <Dialog.Content
        class="animate-slide-in w-full max-w-lg overflow-y-auto bg-white shadow-2xl"
      >
        {#if fullTask}
          <!-- Header -->
          <div
            class="sticky top-0 z-10 border-b border-surface-100 bg-white/95 px-8 pt-8 pb-4 backdrop-blur-sm"
          >
            <div class="flex items-start justify-between gap-4">
              <div class="min-w-0 flex-1">
                <Dialog.Title
                  class="text-2xl leading-tight font-medium tracking-tight text-surface-900"
                >
                  {fullTask.title}
                </Dialog.Title>
                <Dialog.Description class="mt-1 text-sm text-surface-500">
                  Created {formatDateTime(fullTask.createdAt)}
                </Dialog.Description>
              </div>
              <Dialog.CloseTrigger
                class="-mt-1 -mr-2 shrink-0 rounded-xl p-2 text-surface-400 transition-all hover:bg-surface-100 hover:text-surface-900"
              >
                <X size={20} />
              </Dialog.CloseTrigger>
            </div>
          </div>

          <!-- Body -->
          <div class="space-y-6 px-8 py-6">
            <!-- Metadata Grid -->
            <div class="space-y-4">
              <!-- Status -->
              <div class="flex items-center justify-between border-b border-surface-50 py-3">
                <div class="flex items-center gap-3 text-sm text-surface-600">
                  <CheckCircle size={18} strokeWidth={2.5} />
                  <span class="font-medium">Status</span>
                </div>
                {#if status}
                  <span class="inline-flex items-center rounded-full px-3 py-1 text-xs font-semibold {status.class}">
                    {status.label}
                  </span>
                {/if}
              </div>

              <!-- Priority -->
              <div class="flex items-center justify-between border-b border-surface-50 py-3">
                <div class="flex items-center gap-3 text-sm text-surface-600">
                  <Flag size={18} strokeWidth={2.5} />
                  <span class="font-medium">Priority</span>
                </div>
                {#if priority}
                  <span class="inline-flex items-center rounded-full px-3 py-1 text-xs font-semibold {priority.class}">
                    {priority.label}
                  </span>
                {/if}
              </div>

              <!-- Blocked -->
              {#if fullTask.isBlocked}
                <div class="flex items-center justify-between border-b border-surface-50 py-3">
                  <div class="flex items-center gap-2.5 text-sm text-surface-500">
                    <ShieldBan size={16} />
                    <span class="font-medium">Blocked</span>
                  </div>
                  <span class="inline-flex items-center gap-1 rounded-full bg-rose-50 px-3 py-1 text-xs font-semibold text-rose-600">
                    <AlertCircle size={12} />
                    Yes
                  </span>
                </div>
              {/if}

              <!-- Due Date -->
              <div class="flex items-center justify-between border-b border-surface-50 py-3">
                <div class="flex items-center gap-3 text-sm text-surface-600">
                  <Calendar size={18} strokeWidth={2.5} />
                  <span class="font-medium">Due Date</span>
                </div>
                <span class="text-sm font-medium text-surface-900">
                  {fullTask.dueDate ? formatDate(fullTask.dueDate) : '—'}
                </span>
              </div>

              <!-- Estimated Hours -->
              {#if fullTask.estimatedHours}
                <div class="flex items-center justify-between border-b border-surface-50 py-3">
                  <div class="flex items-center gap-3 text-sm text-surface-600">
                    <Clock size={18} strokeWidth={2.5} />
                    <span class="font-medium">Estimated</span>
                  </div>
                  <span class="text-sm font-medium text-surface-900">
                    {formatHours(fullTask.estimatedHours)}
                  </span>
                </div>
              {/if}

              <!-- Actual Hours -->
              {#if fullTask.actualHours}
                <div class="flex items-center justify-between border-b border-surface-50 py-3">
                  <div class="flex items-center gap-2.5 text-sm text-surface-500">
                    <Clock size={16} />
                    <span class="font-medium">Actual</span>
                  </div>
                  <span class="text-sm font-medium text-surface-900">
                    {formatHours(fullTask.actualHours)}
                  </span>
                </div>
              {/if}

              <!-- Assignee -->
              <div class="flex items-center justify-between border-b border-surface-50 py-3">
                <div class="flex items-center gap-2.5 text-sm text-surface-500">
                  <User size={16} />
                  <span class="font-medium">Assignee</span>
                </div>
                <div class="flex items-center gap-2">
                  {#if fullTask.assignee.avatarUrl}
                    <div class="size-7 overflow-hidden rounded-full border-2 border-white bg-surface-100 shadow-sm">
                      <img src={resolveAvatarUrl(fullTask.assignee.avatarUrl) ?? ''} alt={fullTask.assignee.firstName} class="size-full object-cover" />
                    </div>
                  {:else}
                    <div class="flex size-7 items-center justify-center rounded-full border-2 border-white bg-indigo-100 shadow-sm">
                      <span class="text-[10px] font-bold text-indigo-700">
                        {fullTask.assignee.firstName[0]}{fullTask.assignee.lastName[0]}
                      </span>
                    </div>
                  {/if}
                  <span class="text-sm font-medium text-surface-900">
                    {fullTask.assignee.firstName} {fullTask.assignee.lastName}
                  </span>
                </div>
              </div>
            </div>

            <!-- Description -->
            {#if fullTask.description}
              <div class="space-y-3">
                <h4 class="text-sm font-semibold tracking-wider text-surface-900 uppercase">Description</h4>
                <div class="rounded-2xl border border-surface-100 bg-surface-50 p-5">
                  <p class="text-sm leading-relaxed whitespace-pre-wrap text-surface-700">
                    {fullTask.description}
                  </p>
                </div>
              </div>
            {/if}

            <!-- Comments -->
            <div class="space-y-3 border-t border-surface-100 pt-5">
              <div class="flex items-center gap-2">
                <MessageCircle size={15} class="text-surface-500" />
                <h4 class="text-sm font-semibold tracking-wider text-surface-900 uppercase">Comments</h4>
                {#if fullTaskQuery.isFetching}
                  <LoaderCircle size={12} class="animate-spin text-surface-400" />
                {:else}
                  <span class="rounded-full bg-surface-100 px-2 py-0.5 text-xs font-medium text-surface-500">
                    {fullTask.notes?.length ?? fullTask.notesCount}
                  </span>
                {/if}
              </div>

              {#if fullTask.notes && fullTask.notes.length > 0}
                <div class="space-y-3">
                  {#each fullTask.notes as note (note.id)}
                    <div class="flex gap-3">
                      {#if note.user?.avatarUrl}
                        <div class="mt-0.5 size-7 overflow-hidden rounded-full border border-surface-100 bg-surface-100 shadow-sm">
                          <img src={resolveAvatarUrl(note.user.avatarUrl) ?? ''} alt={note.user.firstName} class="size-full object-cover" />
                        </div>
                      {:else if note.user}
                        <div class="mt-0.5 flex size-7 items-center justify-center rounded-full border border-surface-100 bg-indigo-100 text-[10px] font-bold text-indigo-700 shadow-sm">
                          {note.user.firstName[0]}{note.user.lastName[0]}
                        </div>
                      {:else}
                        <div class="mt-0.5 flex size-7 shrink-0 items-center justify-center rounded-full bg-primary-50 text-[10px] font-bold text-primary-600">
                          ?
                        </div>
                      {/if}
                      <div class="min-w-0 flex-1 rounded-xl border border-surface-100 bg-surface-50 px-4 py-3">
                        <div class="mb-1 flex items-center justify-between">
                          <span class="text-xs font-semibold text-surface-900">
                             {note.user ? `${note.user.firstName} ${note.user.lastName}` : 'Unknown User'}
                          </span>
                          <span class="text-[10px] text-surface-400">{formatDateTime(note.createdAt)}</span>
                        </div>
                        <p class="text-sm leading-relaxed text-surface-700">{note.content}</p>
                      </div>
                    </div>
                  {/each}
                </div>
              {:else if !fullTaskQuery.isFetching}
                <p class="text-sm text-surface-400">No comments yet.</p>
              {/if}
            </div>

            <!-- Attachments -->
            <div class="space-y-3 border-t border-surface-100 pt-5">
              <div class="flex items-center gap-2">
                <Paperclip size={15} class="text-surface-500" />
                <h4 class="text-sm font-semibold tracking-wider text-surface-900 uppercase">Attachments</h4>
                {#if fullTaskQuery.isFetching}
                  <LoaderCircle size={12} class="animate-spin text-surface-400" />
                {:else}
                  <span class="rounded-full bg-surface-100 px-2 py-0.5 text-xs font-medium text-surface-500">
                    {fullTask.attachments?.length ?? fullTask.attachmentsCount}
                  </span>
                {/if}
              </div>

              {#if fullTask.attachments && fullTask.attachments.length > 0}
                <div class="space-y-2">
                  {#each fullTask.attachments as att (att.id)}
                    {@const Icon = fileIcon(att.mimeType)}
                    <div class="flex items-center gap-3 rounded-xl border border-surface-100 bg-surface-50 px-4 py-3">
                      <Icon size={18} class="shrink-0 text-surface-400" />
                      <div class="min-w-0 flex-1">
                        <div class="flex items-center gap-2">
                          <p class="truncate text-sm font-medium text-surface-800">{att.fileName}</p>
                          {#if att.user}
                             <span class="text-[10px] text-surface-400">• by {att.user.firstName}</span>
                          {/if}
                        </div>
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
                    </div>
                  {/each}
                </div>
              {:else if !fullTaskQuery.isFetching}
                <p class="text-sm text-surface-400">No attachments yet.</p>
              {/if}
            </div>

            <!-- Timestamps -->
            <div class="space-y-3 border-t border-surface-100 pt-4">
              <h4 class="text-sm font-semibold tracking-wider text-surface-900 uppercase">Timeline</h4>
              <div class="space-y-2">
                <div class="flex items-center justify-between text-xs">
                  <span class="text-surface-500">Created</span>
                  <span class="font-medium text-surface-700">{formatDateTime(fullTask.createdAt)}</span>
                </div>
                <div class="flex items-center justify-between text-xs">
                  <span class="text-surface-500">Last Updated</span>
                  <span class="font-medium text-surface-700">{formatDateTime(fullTask.updatedAt)}</span>
                </div>
                {#if fullTask.completedAt}
                  <div class="flex items-center justify-between text-xs">
                    <span class="text-surface-500">Completed</span>
                    <span class="font-medium text-emerald-600">{formatDateTime(fullTask.completedAt)}</span>
                  </div>
                {/if}
              </div>
            </div>
          </div>
        {/if}
      </Dialog.Content>
    </Dialog.Positioner>
  </Portal>
</Dialog>

<style>
  @keyframes slide-in {
    from {
      transform: translateX(100%);
      opacity: 0.8;
    }
    to {
      transform: translateX(0);
      opacity: 1;
    }
  }

  :global(.animate-slide-in) {
    animation: slide-in 0.3s cubic-bezier(0.16, 1, 0.3, 1);
  }
</style>
