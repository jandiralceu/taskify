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
  } from '@lucide/svelte'
  import { Dialog, Portal } from '@skeletonlabs/skeleton-svelte'
  import type { TaskResponse } from '$lib/api/types'
  import { priorityConfig, statusConfig } from '$lib/utils/task'

  interface Props {
    task: TaskResponse | null
    isOpen: boolean
    onClose: () => void
  }

  let { task, isOpen, onClose }: Props = $props()

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

  let priority = $derived(task ? priorityConfig[task.priority] : null)
  let status = $derived(task ? statusConfig[task.status] : null)
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
        {#if task}
          <!-- Header -->
          <div
            class="sticky top-0 z-10 border-b border-surface-100 bg-white/95 px-8 pt-8 pb-4 backdrop-blur-sm"
          >
            <div class="flex items-start justify-between gap-4">
              <div class="min-w-0 flex-1">
                <Dialog.Title
                  class="text-2xl leading-tight font-medium tracking-tight text-surface-900"
                >
                  {task.title}
                </Dialog.Title>
                <Dialog.Description class="mt-1 text-sm text-surface-500">
                  Created {formatDateTime(task.createdAt)}
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
              <div
                class="flex items-center justify-between border-b border-surface-50 py-3"
              >
                <div class="flex items-center gap-3 text-sm text-surface-600">
                  <CheckCircle size={18} strokeWidth={2.5} />
                  <span class="font-medium">Status</span>
                </div>
                {#if status}
                  <span
                    class="inline-flex items-center rounded-full px-3 py-1 text-xs font-semibold {status.class}"
                  >
                    {status.label}
                  </span>
                {/if}
              </div>

              <!-- Priority -->
              <div
                class="flex items-center justify-between border-b border-surface-50 py-3"
              >
                <div class="flex items-center gap-3 text-sm text-surface-600">
                  <Flag size={18} strokeWidth={2.5} />
                  <span class="font-medium">Priority</span>
                </div>
                {#if priority}
                  <span
                    class="inline-flex items-center rounded-full px-3 py-1 text-xs font-semibold {priority.class}"
                  >
                    {priority.label}
                  </span>
                {/if}
              </div>

              <!-- Blocked -->
              {#if task.isBlocked}
                <div
                  class="flex items-center justify-between border-b border-surface-50 py-3"
                >
                  <div
                    class="flex items-center gap-2.5 text-sm text-surface-500"
                  >
                    <ShieldBan size={16} />
                    <span class="font-medium">Blocked</span>
                  </div>
                  <span
                    class="inline-flex items-center gap-1 rounded-full bg-rose-50 px-3 py-1 text-xs font-semibold text-rose-600"
                  >
                    <AlertCircle size={12} />
                    Yes
                  </span>
                </div>
              {/if}

              <!-- Due Date -->
              <div
                class="flex items-center justify-between border-b border-surface-50 py-3"
              >
                <div class="flex items-center gap-3 text-sm text-surface-600">
                  <Calendar size={18} strokeWidth={2.5} />
                  <span class="font-medium">Due Date</span>
                </div>
                <span class="text-sm font-medium text-surface-900">
                  {task.dueDate ? formatDate(task.dueDate) : '—'}
                </span>
              </div>

              <!-- Estimated Hours -->
              {#if task.estimatedHours}
                <div
                  class="flex items-center justify-between border-b border-surface-50 py-3"
                >
                  <div class="flex items-center gap-3 text-sm text-surface-600">
                    <Clock size={18} strokeWidth={2.5} />
                    <span class="font-medium">Estimated</span>
                  </div>
                  <span class="text-sm font-medium text-surface-900">
                    {formatHours(task.estimatedHours)}
                  </span>
                </div>
              {/if}

              <!-- Actual Hours -->
              {#if task.actualHours}
                <div
                  class="flex items-center justify-between border-b border-surface-50 py-3"
                >
                  <div
                    class="flex items-center gap-2.5 text-sm text-surface-500"
                  >
                    <Clock size={16} />
                    <span class="font-medium">Actual</span>
                  </div>
                  <span class="text-sm font-medium text-surface-900">
                    {formatHours(task.actualHours)}
                  </span>
                </div>
              {/if}

              <!-- Assignee -->
              <div
                class="flex items-center justify-between border-b border-surface-50 py-3"
              >
                <div class="flex items-center gap-2.5 text-sm text-surface-500">
                  <User size={16} />
                  <span class="font-medium">Assignee</span>
                </div>
                <div class="flex items-center gap-2">
                  {#if task.assignee.avatarUrl}
                    <div
                      class="size-7 overflow-hidden rounded-full border-2 border-white bg-surface-100 shadow-sm"
                    >
                      <img
                        src={task.assignee.avatarUrl}
                        alt={task.assignee.firstName}
                        class="size-full object-cover"
                      />
                    </div>
                  {:else}
                    <div
                      class="flex size-7 items-center justify-center rounded-full border-2 border-white bg-indigo-100 shadow-sm"
                    >
                      <span class="text-[10px] font-bold text-indigo-700">
                        {task.assignee.firstName[0]}{task.assignee.lastName[0]}
                      </span>
                    </div>
                  {/if}
                  <span class="text-sm font-medium text-surface-900">
                    {task.assignee.firstName}
                    {task.assignee.lastName}
                  </span>
                </div>
              </div>
            </div>

            <!-- Description -->
            {#if task.description}
              <div class="space-y-3">
                <h4
                  class="text-sm font-semibold tracking-wider text-surface-900 uppercase"
                >
                  Description
                </h4>
                <div
                  class="rounded-2xl border border-surface-100 bg-surface-50 p-5"
                >
                  <p
                    class="text-sm leading-relaxed whitespace-pre-wrap text-surface-700"
                  >
                    {task.description}
                  </p>
                </div>
              </div>
            {/if}

            <!-- Timestamps -->
            <div class="space-y-3 border-t border-surface-100 pt-4">
              <h4
                class="text-sm font-semibold tracking-wider text-surface-900 uppercase"
              >
                Timeline
              </h4>
              <div class="space-y-2">
                <div class="flex items-center justify-between text-xs">
                  <span class="text-surface-500">Created</span>
                  <span class="font-medium text-surface-700"
                    >{formatDateTime(task.createdAt)}</span
                  >
                </div>
                <div class="flex items-center justify-between text-xs">
                  <span class="text-surface-500">Last Updated</span>
                  <span class="font-medium text-surface-700"
                    >{formatDateTime(task.updatedAt)}</span
                  >
                </div>
                {#if task.completedAt}
                  <div class="flex items-center justify-between text-xs">
                    <span class="text-surface-500">Completed</span>
                    <span class="font-medium text-emerald-600"
                      >{formatDateTime(task.completedAt)}</span
                    >
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
