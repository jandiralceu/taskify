<script lang="ts">
  import {
    Ellipsis,
    Pencil,
    ShieldBan,
    ShieldCheck,
    Trash2,
    Flag,
    MessageCircle,
    Paperclip,
    TriangleAlert,
  } from '@lucide/svelte'
  import { Popover, Portal, Dialog } from '@skeletonlabs/skeleton-svelte'
  import { goto } from '$app/navigation'
  import { resolve } from '$app/paths'
  import type { TaskResponse } from '$lib/api/types'
  import { priorityConfig } from '$lib/utils/task'
  import Button from '$lib/components/Button.svelte'

  interface Props {
    task: TaskResponse
    isDragging?: boolean
    onDragStart: (e: DragEvent) => void
    onDragEnd: () => void
    onDelete?: (taskId: string, title: string) => void
    onToggleBlock?: (taskId: string, blocked: boolean) => void
    onViewDetails?: (task: TaskResponse) => void
    onViewUser?: (userId: string) => void
  }

  let {
    task,
    isDragging = false,
    onDragStart,
    onDragEnd,
    onDelete,
    onToggleBlock,
    onViewDetails,
    onViewUser,
  }: Props = $props()
  let isMenuOpen = $state(false)
  let isDeleteConfirmOpen = $state(false)

  function openDeleteConfirm() {
    isMenuOpen = false
    isDeleteConfirmOpen = true
  }

  function confirmDelete() {
    isDeleteConfirmOpen = false
    onDelete?.(task.id, task.title)
  }

  function formatDate(dateStr: string) {
    return new Intl.DateTimeFormat('en-GB', {
      day: '2-digit',
      month: 'short',
      year: 'numeric',
      timeZone: 'UTC',
    }).format(new Date(dateStr))
  }

  function handleDragStart(e: DragEvent) {
    if (task.isBlocked) {
      e.preventDefault()
      return
    }
    onDragStart(e)
  }

  const priority = $derived(priorityConfig[task.priority])
</script>

<div
  role="button"
  tabindex="0"
  draggable={!task.isBlocked}
  ondragstart={handleDragStart}
  ondragend={onDragEnd}
  onclick={() => onViewDetails?.(task)}
  onkeydown={e => {
    if (e.key === 'Enter' || e.key === ' ') {
      e.preventDefault()
      onViewDetails?.(task)
    }
  }}
  class="rounded-2xl border border-1 border-solid bg-white p-5 transition-all {task.isBlocked
    ? 'cursor-not-allowed border-rose-200 bg-surface-50 opacity-60'
    : 'cursor-pointer border-gray-200 hover:shadow-md active:cursor-grabbing'} {isDragging
    ? 'scale-95 opacity-40'
    : ''}"
>
  <!-- Header: Priority + Blocked Badge + Menu -->
  <div class="mb-3 flex items-center justify-between">
    <div class="flex items-center gap-1.5">
      <span
        class="inline-flex items-center rounded-full px-3 py-1 text-[12px] font-semibold {priority.class}"
      >
        {priority.label}
      </span>
      {#if task.isBlocked}
        <span
          class="inline-flex items-center gap-1 rounded-full bg-rose-50 px-2.5 py-1 text-[12px] font-semibold text-rose-600"
        >
          <ShieldBan size={12} />
          Blocked
        </span>
      {/if}
    </div>
    <Popover
      open={isMenuOpen}
      onOpenChange={e => (isMenuOpen = e.open)}
      positioning={{ placement: 'bottom-end' }}
    >
      <Popover.Trigger
        onclick={e => e.stopPropagation()}
        class="cursor-pointer rounded-md p-1 text-surface-400 transition-colors hover:bg-surface-100 hover:text-surface-600"
      >
        <Ellipsis size={18} />
      </Popover.Trigger>
      <Portal>
        <Popover.Positioner>
          <Popover.Content
            class="z-50 w-44 rounded-xl border border-surface-200 bg-white py-1 shadow-lg"
          >
            <button
              onclick={e => {
                e.stopPropagation()
                isMenuOpen = false
                goto(resolve(`/tasks/${task.id}/edit`))
              }}
              class="flex w-full cursor-pointer items-center gap-2.5 px-3 py-2 text-sm text-surface-700 transition-colors hover:bg-surface-50"
            >
              <Pencil size={15} class="text-surface-400" />
              <span>Edit</span>
            </button>
            <button
              onclick={e => {
                e.stopPropagation()
                onToggleBlock?.(task.id, !task.isBlocked)
                isMenuOpen = false
              }}
              class="flex w-full cursor-pointer items-center gap-2.5 px-3 py-2 text-sm text-surface-700 transition-colors hover:bg-surface-50"
            >
              {#if task.isBlocked}
                <ShieldCheck size={15} class="text-emerald-500" />
                <span>Unblock</span>
              {:else}
                <ShieldBan size={15} class="text-surface-400" />
                <span>Block</span>
              {/if}
            </button>
            <hr class="my-1 border-surface-100" />
            <button
              onclick={e => {
                e.stopPropagation()
                openDeleteConfirm()
              }}
              class="flex w-full cursor-pointer items-center gap-2.5 px-3 py-2 text-sm text-rose-600 transition-colors hover:bg-rose-50"
            >
              <Trash2 size={15} />
              <span>Delete</span>
            </button>
          </Popover.Content>
        </Popover.Positioner>
      </Portal>
    </Popover>
  </div>

  <!-- Title -->
  <h4 class="mb-2 leading-snug font-light">{task.title}</h4>

  <!-- Assignees -->
  <div class="mb-4 flex items-center justify-between">
    <span class="text-xs font-medium text-surface-700">Assignees :</span>
    <div class="flex -space-x-2">
      {#if task.assignee.avatarUrl}
        <button
          type="button"
          onclick={e => {
            e.stopPropagation()
            onViewUser?.(task.assignee.id)
          }}
          class="flex size-7 cursor-pointer items-center justify-center overflow-hidden rounded-full border-2 border-white bg-surface-100 transition-all hover:ring-2 hover:ring-primary-500"
          title="{task.assignee.firstName} {task.assignee.lastName}"
        >
          <img
            src={task.assignee.avatarUrl}
            alt={task.assignee.firstName}
            class="size-full object-cover"
          />
        </button>
      {:else}
        <button
          type="button"
          onclick={e => {
            e.stopPropagation()
            onViewUser?.(task.assignee.id)
          }}
          class="flex size-7 cursor-pointer items-center justify-center rounded-full border-2 border-white bg-indigo-100 transition-all hover:ring-2 hover:ring-primary-500"
          title="{task.assignee.firstName} {task.assignee.lastName}"
        >
          <span class="text-[10px] font-bold text-indigo-700">
            {task.assignee.firstName[0]}{task.assignee.lastName[0]}
          </span>
        </button>
      {/if}
    </div>
  </div>

  <!-- Date -->
  <div class="mb-4 flex items-center">
    {#if task.dueDate}
      <div class="flex items-center gap-2 text-xs font-medium text-surface-700">
        <Flag size={14} />
        <span>{formatDate(task.dueDate)}</span>
      </div>
    {/if}
  </div>

  <!-- Footer: Comments, Files -->
  <div
    class="flex items-center justify-between border-t border-surface-100 pt-3 text-xs text-surface-500"
  >
    <div class="flex items-center gap-1.5">
      <Paperclip size={14} />
      <span
        >{task.attachmentsCount}
        {task.attachmentsCount === 1 ? 'Attachment' : 'Attachments'}</span
      >
    </div>
    <div class="flex items-center gap-1.5">
      <MessageCircle size={14} />
      <span
        >{task.notesCount}
        {task.notesCount === 1 ? 'Comment' : 'Comments'}</span
      >
    </div>
  </div>
</div>

<!-- Delete Confirmation Dialog -->
<Dialog
  role="alertdialog"
  open={isDeleteConfirmOpen}
  onOpenChange={e => {
    if (!e.open) isDeleteConfirmOpen = false
  }}
>
  <Portal>
    <Dialog.Backdrop
      class="fixed inset-0 z-50 bg-surface-950/40 backdrop-blur-sm"
    />
    <Dialog.Positioner
      class="fixed inset-0 z-50 flex items-center justify-center p-4"
    >
      <Dialog.Content
        onclick={e => e.stopPropagation()}
        class="w-full max-w-md space-y-4 rounded-2xl border border-surface-100 bg-white p-6 shadow-2xl"
      >
        <div class="flex items-start gap-4">
          <div
            class="flex size-10 shrink-0 items-center justify-center rounded-xl bg-rose-50 text-rose-500"
          >
            <TriangleAlert size={20} />
          </div>
          <div>
            <Dialog.Title class="text-base font-bold text-surface-900">
              Delete Task
            </Dialog.Title>
            <Dialog.Description class="mt-1 text-sm text-surface-500">
              Are you sure you want to delete
              <span class="font-semibold text-surface-700">"{task.title}"</span
              >? This action cannot be undone.
            </Dialog.Description>
          </div>
        </div>

        <div class="flex items-center justify-end gap-3 pt-2">
          <Dialog.CloseTrigger
            type="button"
            class="rounded-xl px-5 py-2.5 text-sm font-medium text-surface-500 transition-all hover:bg-surface-50 hover:text-surface-900"
          >
            Cancel
          </Dialog.CloseTrigger>
          <Button variant="danger" onclick={confirmDelete}>
            <Trash2 size={15} />
            Delete
          </Button>
        </div>
      </Dialog.Content>
    </Dialog.Positioner>
  </Portal>
</Dialog>
