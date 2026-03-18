<script lang="ts">
  import {
    Plus,
    Ellipsis,
    LoaderCircle,
    Search,
    ArrowUp,
    ArrowDown,
    Ban,
    X,
  } from '@lucide/svelte'
  import TaskCard from '$lib/components/TaskCard.svelte'
  import TaskDetailDrawer from '$lib/components/TaskDetailDrawer.svelte'
  import { createProfileQuery } from '$lib/state/user.svelte'
  import {
    getTasksQuery,
    getArchivedTasksQuery,
    updateTaskMutation,
    deleteTaskMutation,
  } from '$lib/state/tasks.svelte'
  import AddTaskModal from '$lib/components/AddTaskModal.svelte'
  import UserDetailDrawer from '$lib/components/UserDetailDrawer.svelte'
  import ArchivedTaskTable from '$lib/components/ArchivedTaskTable.svelte'
  import { toaster } from '$lib/state/toast.svelte'
  import Button from '$lib/components/Button.svelte'
  import type {
    TaskResponse,
    TaskStatus,
    TaskPriority,
    UserRole,
  } from '$lib/api/types'

  const ADMIN: UserRole = 'admin'

  // Filter & sort state
  let filterSearch = $state('')
  let debouncedSearch = $state('')
  let filterPriority = $state<TaskPriority | ''>('')
  let filterBlocked = $state(false)
  let filterArchived = $state(false)
  let sortField = $state('createdAt')
  let sortOrder = $state<'asc' | 'desc'>('desc')

  const profileQuery = createProfileQuery()

  const commonParams = () => ({
    search: debouncedSearch || undefined,
    priority: (filterPriority as TaskPriority) || undefined,
    isBlocked: filterBlocked || undefined,
    sort: sortField,
    order: sortOrder,
  })

  const activeTasksQuery = getTasksQuery(commonParams)
  const archivedTasksQuery = getArchivedTasksQuery(commonParams)

  const tasksQuery = $derived(
    filterArchived ? archivedTasksQuery : activeTasksQuery
  )

  const updateTask = updateTaskMutation()
  const deleteTask = deleteTaskMutation()

  const columns: { id: TaskStatus; title: string }[] = [
    { id: 'pending', title: 'Pending' },
    { id: 'in_progress', title: 'In Progress' },
    { id: 'completed', title: 'Completed' },
    { id: 'cancelled', title: 'Cancelled' },
  ]

  let isModalOpen = $state(false)
  let selectedTask = $state<TaskResponse | null>(null)
  let isDrawerOpen = $state(false)
  let selectedUserId = $state<string | undefined>(undefined)
  let isUserDrawerOpen = $state(false)

  const hasActiveFilters = $derived(
    filterSearch !== '' ||
      filterPriority !== '' ||
      filterBlocked ||
      filterArchived
  )

  function clearFilters() {
    filterSearch = ''
    debouncedSearch = ''
    filterPriority = ''
    filterBlocked = false
    filterArchived = false
  }

  $effect(() => {
    const query = filterSearch // Read synchronously to track as dependency
    const handler = setTimeout(() => {
      debouncedSearch = query
    }, 300)

    return () => clearTimeout(handler)
  })

  /**
   * Drag-and-drop state.
   *
   * - draggingTaskId: the id of the task card currently being dragged.
   *   Used to apply a visual "ghost" style (reduced opacity + scale) to
   *   the source card while it is in flight.
   *
   * - dragOverColumn: the id of the column currently under the dragged card.
   *   Used to highlight the target drop zone with a purple tint and ring.
   */
  let draggingTaskId = $state<string | null>(null)
  let dragOverColumn = $state<string | null>(null)

  let formattedDate = $derived.by(() => {
    return new Intl.DateTimeFormat('en-GB', {
      day: 'numeric',
      month: 'long',
      year: 'numeric',
    }).format(new Date())
  })

  function handleAddTask() {
    isModalOpen = true
  }

  // ---------------------------------------------------------------------------
  // Drag-and-drop handlers
  //
  // Flow overview:
  //   1. User starts dragging a card  → onDragStart
  //   2. Card enters a column         → onDragOver  (fires repeatedly)
  //   3. Card leaves a column         → onDragLeave
  //   4. User releases the card       → onDrop  (or onDragEnd on cancel)
  //
  // The task id is passed through the native DataTransfer API so it survives
  // across all drag events without relying on closure state.
  // ---------------------------------------------------------------------------

  /**
   * Fired when the user begins dragging a task card.
   * Stores the task id in both reactive state (for visual feedback) and in
   * the DataTransfer object (for retrieval on drop).
   */
  function onDragStart(e: DragEvent, taskId: string) {
    draggingTaskId = taskId
    e.dataTransfer!.effectAllowed = 'move'
    e.dataTransfer!.setData('text/plain', taskId)
  }

  /**
   * Fired when the drag operation ends, regardless of whether a drop occurred.
   * Resets all drag-related visual state to ensure no column stays highlighted.
   */
  function onDragEnd() {
    draggingTaskId = null
    dragOverColumn = null
  }

  /**
   * Fired continuously while a dragged card hovers over a column.
   * preventDefault() is required to allow the drop event to fire.
   */
  function onDragOver(e: DragEvent, columnId: string) {
    e.preventDefault()
    e.dataTransfer!.dropEffect = 'move'
    dragOverColumn = columnId
  }

  /**
   * Fired when the dragged card leaves a column's bounding box.
   * The relatedTarget check prevents false positives caused by the cursor
   * briefly leaving and re-entering child elements (cards, buttons, etc.)
   * inside the same column.
   */
  function onDragLeave(e: DragEvent, columnId: string) {
    const related = e.relatedTarget as Node | null
    const target = e.currentTarget as HTMLElement
    if (!target.contains(related)) {
      if (dragOverColumn === columnId) dragOverColumn = null
    }
  }

  /**
   * Fired when the user drops a card onto a column.
   * Reads the task id from DataTransfer, guards against no-op moves
   * (dropping onto the same column), then calls the update mutation to
   * persist the new status on the backend. TanStack Query's cache
   * invalidation in the mutation's onSuccess triggers an automatic re-fetch,
   * moving the card to the correct column in the UI.
   */
  async function onDrop(e: DragEvent, targetStatus: TaskStatus) {
    e.preventDefault()
    const taskId = e.dataTransfer!.getData('text/plain')
    dragOverColumn = null
    draggingTaskId = null

    const task = tasksQuery.data?.find(t => t.id === taskId)
    if (!task || task.status === targetStatus) return

    await updateTask.mutateAsync({ id: taskId, data: { status: targetStatus } })
  }

  async function handleDeleteTask(taskId: string, title?: string) {
    try {
      await deleteTask.mutateAsync(taskId)
      toaster.success({
        title: 'Task Deleted',
        description: title
          ? `"${title}" has been deleted.`
          : 'Task has been deleted.',
      })
    } catch {
      toaster.error({
        title: 'Delete Failed',
        description: 'Could not delete the task. Please try again.',
      })
    }
  }

  async function handleToggleBlock(taskId: string, blocked: boolean) {
    await updateTask.mutateAsync({ id: taskId, data: { isBlocked: blocked } })
  }

  function handleViewDetails(task: TaskResponse) {
    selectedTask = task
    isDrawerOpen = true
  }

  function handleViewUser(userId: string) {
    selectedUserId = userId
    isUserDrawerOpen = true
  }
</script>

<svelte:head>
  <title>Tasks - Taskify</title>
</svelte:head>

<div class="flex min-h-full flex-col">
  <!-- Project Header -->
  <header class="px-16 pt-12 pb-10">
    <div class="flex items-end justify-between">
      <div class="space-y-2">
        <h2 class="text-4xl leading-tight tracking-tight text-surface-900">
          <span class="font-light">Welcome</span>
          <span class="font-normal"
            >{profileQuery.data?.firstName || 'User'}</span
          >,
          {#if profileQuery.data?.role === ADMIN}
            <span class="font-light">here's an overview of</span> <br />
            <span class="font-normal">all tasks across the team!</span>
          {:else}
            <span class="font-light">here's a</span> <br />
            <span class="font-normal">look at your tasks!</span>
          {/if}
        </h2>
        <p class="text-xl font-light text-surface-800">
          Today is {formattedDate}
        </p>
      </div>
    </div>
  </header>

  <div
    class="sticky top-0 z-20 mb-4 border-b border-surface-100/50 bg-[#F7F3F9]/95 px-16 py-6 backdrop-blur-sm transition-all"
  >
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-10">
        <button
          onclick={() => (filterArchived = false)}
          class="group relative transition-all"
        >
          <h3
            class="text-3xl tracking-tight transition-colors {!filterArchived
              ? 'font-normal text-surface-900'
              : 'font-light text-surface-400 hover:text-surface-600'}"
          >
            Tasks
          </h3>
        </button>

        <button
          onclick={() => (filterArchived = true)}
          class="group relative transition-all"
        >
          <h3
            class="text-3xl tracking-tight transition-colors {filterArchived
              ? 'font-normal text-surface-900'
              : 'font-light text-surface-400 hover:text-surface-600'}"
          >
            Archived
          </h3>
        </button>
      </div>
    </div>

    <!-- Filter Bar -->
    <div class="mt-8 flex flex-wrap items-center justify-between gap-4">
      <div class="flex flex-wrap items-center gap-3">
        <!-- Search -->
        <div class="relative w-72">
          <Search
            size={16}
            class="pointer-events-none absolute top-1/2 left-3.5 -translate-y-1/2 text-surface-400"
          />
          <input
            type="text"
            placeholder="Search tasks..."
            bind:value={filterSearch}
            class="filter-input h-11 w-full py-2.5 pr-4 pl-10"
          />
        </div>

        <div class="mx-1 h-6 w-px bg-surface-200"></div>

        <!-- Priority Filter -->
        <select
          bind:value={filterPriority}
          class="filter-input filter-select h-11 min-w-[150px] px-4"
        >
          <option value="">All Priorities</option>
          <option value="low">Low</option>
          <option value="medium">Medium</option>
          <option value="high">High</option>
          <option value="critical">Critical</option>
        </select>

        <!-- Blocked Toggle -->
        <button
          onclick={() => (filterBlocked = !filterBlocked)}
          type="button"
          class="inline-flex h-11 items-center gap-2 rounded-xl border px-4 text-sm font-medium transition-all {filterBlocked
            ? 'border-red-300 bg-red-50 text-red-700 shadow-sm shadow-red-100'
            : 'border-surface-200 bg-white text-surface-600 hover:border-surface-300'}"
        >
          <Ban size={15} />
          Blocked
        </button>

        {#if hasActiveFilters}
          <div class="ml-2 flex items-center gap-3">
            <span class="text-xs font-medium text-surface-400 tabular-nums">
              {tasksQuery.data?.length ?? 0} matches
            </span>
            <button
              onclick={clearFilters}
              type="button"
              class="inline-flex items-center gap-1 rounded-xl bg-surface-100 px-3 py-1.5 text-xs font-bold text-surface-700 transition-all hover:bg-surface-200"
            >
              <X size={12} strokeWidth={3} />
              RESET
            </button>
          </div>
        {/if}
      </div>

      <div class="flex items-center gap-4">
        <!-- Sort Field -->
        <div
          class="flex h-11 items-center overflow-hidden rounded-xl border border-surface-200 bg-white"
        >
          <select
            bind:value={sortField}
            class="h-full cursor-pointer appearance-none border-none bg-transparent bg-[url('data:image/svg+xml,%3Csvg%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20width%3D%2216%22%20height%3D%2216%22%20viewBox%3D%220%200%2024%2024%22%20fill%3D%22none%22%20stroke%3D%22%2364748b%22%20stroke-width%3D%222%22%20stroke-linecap%3D%22round%22%20stroke-linejoin%3D%22round%22%3E%3Cpath%20d%3D%22m6%209%206%206%206-6%22%2F%3E%3C%2Fsvg%3E')] bg-[position:right_0.75rem_center] bg-no-repeat px-4 pr-8 text-sm font-medium text-surface-700 outline-none focus:ring-0"
          >
            <option value="createdAt">Date Created</option>
            <option value="title">Title</option>
            <option value="priority">Priority</option>
            <option value="dueDate">Due Date</option>
          </select>
          <div class="h-full w-px bg-surface-100"></div>
          <!-- Sort Order Toggle -->
          <button
            onclick={() => (sortOrder = sortOrder === 'asc' ? 'desc' : 'asc')}
            type="button"
            class="h-full px-3 text-surface-500 transition-all hover:bg-surface-50 hover:text-surface-900"
            title={sortOrder === 'asc' ? 'Ascending' : 'Descending'}
          >
            {#if sortOrder === 'asc'}
              <ArrowUp size={16} strokeWidth={2.5} />
            {:else}
              <ArrowDown size={16} strokeWidth={2.5} />
            {/if}
          </button>
        </div>

        <button
          onclick={handleAddTask}
          type="button"
          class="flex items-center gap-2 rounded-xl bg-primary-500 px-6 py-2.5 font-semibold text-white shadow-lg shadow-primary-500/20 transition-all hover:bg-primary-600 hover:shadow-primary-500/30 active:scale-95"
        >
          <Plus size={18} strokeWidth={2.5} />
          Add Task
        </button>
      </div>
    </div>
  </div>

  <!-- Main Content -->
  {#if !filterArchived}
    <!-- Board View (Kanban) -->
    <div class="custom-scrollbar-h flex-1 overflow-x-auto overflow-y-hidden">
      <div class="inline-flex h-[calc(100vh-100px)] pr-16 pb-8 pl-11">
        {#each columns as column (column.id)}
          <div
            class="flex w-[340px] shrink-0 flex-col gap-6 border-r border-slate-300/50 px-5 last:border-r-0"
          >
            <!-- Column Header -->
            <div class="flex items-center justify-between">
              <h3 class="text-sm font-medium text-surface-900">
                {column.title}
              </h3>
              <button
                class="text-surface-300 transition-colors hover:text-surface-500"
              >
                <Ellipsis size={18} />
              </button>
            </div>

            <!-- Cards Area (drop zone) -->
            <div class="relative flex min-h-0 flex-1 flex-col">
              <div
                class="custom-scrollbar flex-1 space-y-4 overflow-y-auto rounded-xl pr-2 pb-16 transition-colors {dragOverColumn ===
                column.id
                  ? 'bg-primary-500/5 ring-2 ring-primary-500/20'
                  : ''}"
                ondragover={e => onDragOver(e, column.id)}
                ondragleave={e => onDragLeave(e, column.id)}
                ondrop={e => onDrop(e, column.id)}
                role="list"
              >
                {#if tasksQuery.isPending}
                  {#each Array(3) as _, i (i)}
                    <div
                      class="h-32 w-full animate-pulse rounded-2xl bg-surface-50"
                    ></div>
                  {/each}
                {:else if tasksQuery.data}
                  {#each tasksQuery.data.filter(t => t.status === column.id) as task (task.id)}
                    <TaskCard
                      {task}
                      isDragging={draggingTaskId === task.id}
                      onDragStart={e => onDragStart(e, task.id)}
                      {onDragEnd}
                      onDelete={handleDeleteTask}
                      onToggleBlock={handleToggleBlock}
                      onViewDetails={handleViewDetails}
                      onViewUser={handleViewUser}
                    />
                  {/each}
                {/if}

                <!-- Plus Button only on Pending column -->
                {#if column.id === 'pending'}
                  <button
                    onclick={handleAddTask}
                    class="flex w-full items-center justify-center gap-2 rounded-xl border-2 border-dashed border-primary-500/30 py-3 text-sm font-medium text-primary-500/60 transition-all hover:border-primary-500 hover:bg-primary-500/5 hover:text-primary-500"
                  >
                    <Plus size={16} />
                    Add task
                  </button>
                {/if}
              </div>
              <!-- Bottom Fade Overlay -->
              <div
                class="pointer-events-none absolute bottom-0 left-0 h-16 w-full bg-gradient-to-t from-[#F7F3F9] to-transparent"
              ></div>
            </div>
          </div>
        {/each}
      </div>
    </div>
  {:else}
    <!-- Archived Table View -->
    <div class="flex-1 overflow-y-auto px-16 pt-4 pb-12">
      {#if tasksQuery.isPending}
        <div class="flex h-64 items-center justify-center">
          <LoaderCircle class="animate-spin text-primary-500" size={32} />
        </div>
      {:else if tasksQuery.isError}
        <div class="rounded-xl bg-rose-50 p-8 text-center">
          <p class="text-sm font-medium text-rose-600">
            Failed to load archived tasks.
          </p>
          <Button
            variant="ghost"
            onclick={() => tasksQuery.refetch()}
            class="mt-2">Try Again</Button
          >
        </div>
      {:else if tasksQuery.data}
        <ArchivedTaskTable
          tasks={tasksQuery.data}
          onViewDetails={handleViewDetails}
        />
      {/if}
    </div>
  {/if}
</div>

<AddTaskModal isOpen={isModalOpen} onClose={() => (isModalOpen = false)} />
<TaskDetailDrawer
  task={selectedTask}
  isOpen={isDrawerOpen}
  onClose={() => {
    isDrawerOpen = false
    selectedTask = null
  }}
/>
<UserDetailDrawer
  userId={selectedUserId}
  open={isUserDrawerOpen}
  onOpenChange={open => (isUserDrawerOpen = open)}
/>

<style>
  /* Vertical Scrollbar */
  .custom-scrollbar::-webkit-scrollbar {
    width: 4px;
  }
  .custom-scrollbar::-webkit-scrollbar-track {
    background: transparent;
  }
  .custom-scrollbar::-webkit-scrollbar-thumb {
    background: #e2e8f0;
    border-radius: 20px;
  }

  /* Horizontal Scrollbar */
  .custom-scrollbar-h::-webkit-scrollbar {
    height: 6px;
  }
  .custom-scrollbar-h::-webkit-scrollbar-track {
    background: transparent;
  }
  .custom-scrollbar-h::-webkit-scrollbar-thumb {
    background: #e2e8f0;
    border-radius: 20px;
  }

  /* Filter inputs */
  .filter-input {
    font-size: 0.8125rem;
    font-weight: 500;
    background: white;
    border: 1px solid #e2e8f0;
    border-radius: 0.75rem;
    outline: none;
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
    color: #334155;
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.03);
  }
  .filter-input::placeholder {
    color: #94a3b8;
    font-weight: 400;
  }
  .filter-input:hover {
    border-color: #cbd5e1;
    background: #f8fafc;
  }
  .filter-input:focus {
    border-color: #8b5cf6;
    background: white;
    box-shadow: 0 0 0 4px rgba(139, 92, 246, 0.12);
  }
  .filter-select {
    padding: 0.5rem 2.25rem 0.5rem 1rem;
    cursor: pointer;
    appearance: none;
    background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='16' height='16' viewBox='0 0 24 24' fill='none' stroke='%2364748b' stroke-width='2.5' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpath d='m6 9 6 6 6-6'/%3E%3C/svg%3E");
    background-repeat: no-repeat;
    background-position: right 0.75rem center;
    background-size: 0.9rem;
  }
</style>
