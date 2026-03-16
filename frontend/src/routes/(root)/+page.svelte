<script lang="ts">
	import { Plus, Ellipsis, LoaderCircle, Search, ArrowUp, ArrowDown, Ban, X, FilterX } from '@lucide/svelte';
	import TaskCard from '$lib/components/TaskCard.svelte';
	import TaskDetailDrawer from '$lib/components/TaskDetailDrawer.svelte';
	import { createProfileQuery } from '$lib/state/user.svelte';
	import { getTasksQuery, updateTaskMutation, deleteTaskMutation } from '$lib/state/tasks.svelte';
	import AddTaskModal from '$lib/components/AddTaskModal.svelte';
	import UserDetailDrawer from '$lib/components/UserDetailDrawer.svelte';
	import { toaster } from '$lib/state/toast.svelte';
	import type { TaskResponse, TaskStatus, TaskPriority, UserRole } from '$lib/api/types';

	const ADMIN: UserRole = 'admin';

	// Filter & sort state
	let filterSearch = $state('');
	let debouncedSearch = $state('');
	let filterPriority = $state<TaskPriority | ''>('');
	let filterBlocked = $state(false);
	let sortField = $state('createdAt');
	let sortOrder = $state<'asc' | 'desc'>('desc');

	const profileQuery = createProfileQuery();
	
	const tasksQuery = getTasksQuery(() => ({
		search: debouncedSearch || undefined,
		priority: (filterPriority as TaskPriority) || undefined,
		isBlocked: filterBlocked || undefined,
		sort: sortField,
		order: sortOrder
	}));

	const updateTask = updateTaskMutation();
	const deleteTask = deleteTaskMutation();

	const columns: { id: TaskStatus; title: string }[] = [
		{ id: 'pending', title: 'Pending' },
		{ id: 'in_progress', title: 'In Progress' },
		{ id: 'completed', title: 'Completed' },
		{ id: 'cancelled', title: 'Cancelled' }
	];

	let isModalOpen = $state(false);
	let selectedTask = $state<TaskResponse | null>(null);
	let isDrawerOpen = $state(false);
	let selectedUserId = $state<string | undefined>(undefined);
	let isUserDrawerOpen = $state(false);

	const hasActiveFilters = $derived(filterSearch !== '' || filterPriority !== '' || filterBlocked);

	function clearFilters() {
		filterSearch = '';
		debouncedSearch = '';
		filterPriority = '';
		filterBlocked = false;
	}

	$effect(() => {
		const query = filterSearch; // Read synchronously to track as dependency
		const handler = setTimeout(() => {
			debouncedSearch = query;
		}, 300);

		return () => clearTimeout(handler);
	});

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
	let draggingTaskId = $state<string | null>(null);
	let dragOverColumn = $state<string | null>(null);

	let formattedDate = $derived.by(() => {
		return new Intl.DateTimeFormat('en-GB', {
			day: 'numeric',
			month: 'long',
			year: 'numeric'
		}).format(new Date());
	});

	function handleAddTask() {
		isModalOpen = true;
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
		draggingTaskId = taskId;
		e.dataTransfer!.effectAllowed = 'move';
		e.dataTransfer!.setData('text/plain', taskId);
	}

	/**
	 * Fired when the drag operation ends, regardless of whether a drop occurred.
	 * Resets all drag-related visual state to ensure no column stays highlighted.
	 */
	function onDragEnd() {
		draggingTaskId = null;
		dragOverColumn = null;
	}

	/**
	 * Fired continuously while a dragged card hovers over a column.
	 * preventDefault() is required to allow the drop event to fire.
	 */
	function onDragOver(e: DragEvent, columnId: string) {
		e.preventDefault();
		e.dataTransfer!.dropEffect = 'move';
		dragOverColumn = columnId;
	}

	/**
	 * Fired when the dragged card leaves a column's bounding box.
	 * The relatedTarget check prevents false positives caused by the cursor
	 * briefly leaving and re-entering child elements (cards, buttons, etc.)
	 * inside the same column.
	 */
	function onDragLeave(e: DragEvent, columnId: string) {
		const related = e.relatedTarget as Node | null;
		const target = e.currentTarget as HTMLElement;
		if (!target.contains(related)) {
			if (dragOverColumn === columnId) dragOverColumn = null;
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
		e.preventDefault();
		const taskId = e.dataTransfer!.getData('text/plain');
		dragOverColumn = null;
		draggingTaskId = null;

		const task = tasksQuery.data?.find(t => t.id === taskId);
		if (!task || task.status === targetStatus) return;

		await updateTask.mutateAsync({ id: taskId, data: { status: targetStatus } });
	}

	async function handleDeleteTask(taskId: string, title?: string) {
		try {
			await deleteTask.mutateAsync(taskId);
			toaster.success({
				title: 'Task Deleted',
				description: title ? `"${title}" has been deleted.` : 'Task has been deleted.'
			});
		} catch {
			toaster.error({
				title: 'Delete Failed',
				description: 'Could not delete the task. Please try again.'
			});
		}
	}

	async function handleToggleBlock(taskId: string, blocked: boolean) {
		await updateTask.mutateAsync({ id: taskId, data: { isBlocked: blocked } });
	}

	function handleViewDetails(task: TaskResponse) {
		selectedTask = task;
		isDrawerOpen = true;
	}

	function handleViewUser(userId: string) {
		selectedUserId = userId;
		isUserDrawerOpen = true;
	}
</script>

<div class="min-h-full flex flex-col">
	<!-- Project Header -->
	<header class="px-16 pt-12 pb-10">
		<div class="flex items-end justify-between">
			<div class="space-y-2">
				<h2 class="text-4xl text-surface-900 tracking-tight leading-tight">
					<span class="font-light">Welcome</span>
					<span class="font-normal">{profileQuery.data?.firstName || 'User'}</span>,
					{#if profileQuery.data?.role === ADMIN}
						<span class="font-light">here's an overview of</span> <br />
						<span class="font-normal">all tasks across the team!</span>
					{:else}
						<span class="font-light">here's a</span> <br />
						<span class="font-normal">look at your tasks!</span>
					{/if}
				</h2>
				<p class="text-surface-800 text-xl font-light">
					Today is {formattedDate}
				</p>
			</div>
		</div>
	</header>

	<div class="sticky top-0 z-20 bg-[#F7F3F9]/95 backdrop-blur-sm px-16 py-5 mb-4 border-b border-surface-100/50 transition-all space-y-4">
		<div class="flex items-center justify-between">
			<h3 class="text-3xl font-light text-surface-900 tracking-tight">Tasks</h3>
			<button
				onclick={handleAddTask}
				type="button"
				class="bg-primary-500 hover:bg-primary-700 text-white px-5 py-2 rounded-xl font-medium transition-all active:scale-95 flex items-center gap-2"
			>
				<Plus size={18} />
				Add Task
			</button>
		</div>

		<!-- Filter Bar -->
		<div class="flex items-center gap-3 flex-wrap">
			<!-- Search -->
			<div class="relative flex-1 max-w-xs">
				<Search size={16} class="absolute left-3 top-1/2 -translate-y-1/2 text-surface-400 pointer-events-none" />
				<input
					type="text"
					placeholder="Search tasks..."
					bind:value={filterSearch}
					class="filter-input w-full pl-9 pr-3 py-2"
				/>
			</div>

			<!-- Priority Filter -->
			<select bind:value={filterPriority} class="filter-input filter-select">
				<option value="">All Priorities</option>
				<option value="low">🟢 Low</option>
				<option value="medium">🟡 Medium</option>
				<option value="high">🟠 High</option>
				<option value="critical">🔴 Critical</option>
			</select>

			<!-- Blocked Toggle -->
			<button
				onclick={() => (filterBlocked = !filterBlocked)}
				type="button"
				class="inline-flex items-center gap-1.5 px-3 py-2 text-sm font-medium rounded-xl border transition-all {filterBlocked
					? 'bg-red-50 border-red-300 text-red-700 shadow-sm shadow-red-100'
					: 'bg-white border-surface-200 text-surface-600 hover:border-surface-300'}"
			>
				<Ban size={14} />
				Blocked
			</button>

			<!-- Divider -->
			<div class="w-px h-6 bg-surface-200"></div>

			<!-- Sort Field -->
			<select bind:value={sortField} class="filter-input filter-select">
				<option value="createdAt">Sort: Date Created</option>
				<option value="title">Sort: Title</option>
				<option value="priority">Sort: Priority</option>
				<option value="dueDate">Sort: Due Date</option>
			</select>

			<!-- Sort Order Toggle -->
			<button
				onclick={() => (sortOrder = sortOrder === 'asc' ? 'desc' : 'asc')}
				type="button"
				class="p-2 bg-white border border-surface-200 rounded-xl hover:border-surface-300 text-surface-600 transition-all"
				title={sortOrder === 'asc' ? 'Ascending' : 'Descending'}
			>
				{#if sortOrder === 'asc'}
					<ArrowUp size={16} />
				{:else}
					<ArrowDown size={16} />
				{/if}
			</button>

			<!-- Clear Filters -->
			{#if hasActiveFilters}
				<span class="text-xs text-surface-400 tabular-nums">
					{tasksQuery.data?.length ?? 0} matches
				</span>
				<button
					onclick={clearFilters}
					type="button"
					class="inline-flex items-center gap-1 px-3 py-2 text-sm font-medium text-primary-600 hover:text-primary-800 hover:bg-primary-50 rounded-xl transition-all"
				>
					<X size={14} />
					Clear
				</button>
			{/if}
		</div>
	</div>

	<!-- Board Content -->
	<div class="flex-1 overflow-x-auto overflow-y-hidden custom-scrollbar-h">
		<div class="inline-flex h-[calc(100vh-100px)] pl-11 pr-16 pb-8">
			{#each columns as column (column.id)}
				<div class="w-[340px] flex flex-col gap-6 shrink-0 border-r border-slate-300/50 last:border-r-0 px-5">
					<!-- Column Header -->
					<div class="flex items-center justify-between">
						<h3 class="text-sm font-medium text-surface-900">{column.title}</h3>
						<button class="text-surface-300 hover:text-surface-500 transition-colors">
							<Ellipsis size={18} />
						</button>
					</div>

					<!-- Cards Area (drop zone) -->
					<div class="relative flex-1 flex flex-col min-h-0">
						<div
							class="flex-1 overflow-y-auto pr-2 space-y-4 custom-scrollbar rounded-xl transition-colors pb-16 {dragOverColumn === column.id ? 'bg-primary-500/5 ring-2 ring-primary-500/20' : ''}"
							ondragover={(e) => onDragOver(e, column.id)}
							ondragleave={(e) => onDragLeave(e, column.id)}
							ondrop={(e) => onDrop(e, column.id)}
							role="list"
						>
							{#if tasksQuery.isPending}
								<div class="flex flex-col items-center justify-center py-12 text-surface-400">
									<LoaderCircle size={24} class="animate-spin mb-2" />
									<span class="text-xs font-medium">Loading tasks...</span>
								</div>
							{:else if tasksQuery.isError}
								<div class="p-4 bg-red-50 text-red-600 rounded-xl text-xs font-medium text-center">
									Failed to load tasks
								</div>
							{:else if tasksQuery.data}
								{#each tasksQuery.data.filter(t => t.status === column.id) as task (task.id)}
									<TaskCard
										{task}
										isDragging={draggingTaskId === task.id}
										onDragStart={(e) => onDragStart(e, task.id)}
										{onDragEnd}
										onDelete={handleDeleteTask}
										onToggleBlock={handleToggleBlock}
										onViewDetails={handleViewDetails}
										onViewUser={handleViewUser}
									/>
								{:else}
									{#if hasActiveFilters}
										<div class="flex flex-col items-center justify-center py-12 px-4 text-center bg-white/50 rounded-2xl border-2 border-dashed border-surface-200 animate-in fade-in zoom-in duration-300">
											<div class="p-3 bg-surface-100 rounded-full mb-3 text-surface-400">
												<FilterX size={20} />
											</div>
											<p class="text-sm font-medium text-surface-600">No matches found</p>
											<p class="text-xs text-surface-400 mt-1">Refine your search parameters</p>
										</div>
									{/if}
								{/each}
							{/if}

							<!-- Plus Button only on Pending column -->
							{#if column.id === 'pending'}
								<button
									onclick={handleAddTask}
									class="w-full py-3 rounded-xl border-2 border-dashed border-primary-500/30 text-primary-500/60 hover:text-primary-500 hover:border-primary-500 hover:bg-primary-500/5 transition-all flex items-center justify-center gap-2 font-medium text-sm"
								>
									<Plus size={16} />
									Add task
								</button>
							{/if}
						</div>
						<!-- Bottom Fade Overlay -->
						<div class="absolute bottom-0 left-0 right-0 h-12 bg-gradient-to-t from-[#F7F3F9] to-transparent pointer-events-none z-10 rounded-b-xl"></div>
					</div>
				</div>
			{/each}
		</div>
	</div>
</div>

<AddTaskModal isOpen={isModalOpen} onClose={() => isModalOpen = false} />
<TaskDetailDrawer task={selectedTask} isOpen={isDrawerOpen} onClose={() => { isDrawerOpen = false; selectedTask = null; }} />
<UserDetailDrawer userId={selectedUserId} open={isUserDrawerOpen} onOpenChange={(open) => isUserDrawerOpen = open} />

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
		font-size: 0.875rem;
		background: white;
		border: 1px solid var(--color-surface-200, #e2e8f0);
		border-radius: 0.75rem;
		outline: none;
		transition: all 0.15s ease;
		color: var(--color-surface-700, #334155);
	}
	.filter-input::placeholder {
		color: var(--color-surface-400, #94a3b8);
	}
	.filter-input:focus {
		border-color: var(--color-primary-400, #a78bfa);
		box-shadow: 0 0 0 3px rgba(139, 92, 246, 0.15);
	}
	.filter-select {
		padding: 0.5rem 2rem 0.5rem 0.75rem;
		cursor: pointer;
		appearance: none;
		background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='16' height='16' viewBox='0 0 24 24' fill='none' stroke='%236b7280' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpath d='m6 9 6 6 6-6'/%3E%3C/svg%3E");
		background-repeat: no-repeat;
		background-position: right 0.5rem center;
		background-size: 1rem;
	}
</style>
