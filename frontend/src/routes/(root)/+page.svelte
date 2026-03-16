<script lang="ts">
	import { Plus, Ellipsis, LoaderCircle } from '@lucide/svelte';
	import TaskCard from '$lib/components/TaskCard.svelte';
	import { createProfileQuery } from '$lib/state/user.svelte';
	import { getTasksQuery, updateTaskMutation, deleteTaskMutation } from '$lib/state/tasks.svelte';
	import AddTaskModal from '$lib/components/AddTaskModal.svelte';
	import type { TaskStatus, UserRole } from '$lib/api/types';

	const ADMIN: UserRole = 'admin';

	const profileQuery = createProfileQuery();
	const tasksQuery = getTasksQuery();
	const updateTask = updateTaskMutation();
	const deleteTask = deleteTaskMutation();

	const columns: { id: TaskStatus; title: string }[] = [
		{ id: 'pending', title: 'Pending' },
		{ id: 'in_progress', title: 'In Progress' },
		{ id: 'completed', title: 'Completed' },
		{ id: 'cancelled', title: 'Cancelled' }
	];

	let isModalOpen = $state(false);

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

	async function handleDeleteTask(taskId: string) {
		await deleteTask.mutateAsync(taskId);
	}

	async function handleToggleBlock(taskId: string, blocked: boolean) {
		await updateTask.mutateAsync({ id: taskId, data: { isBlocked: blocked } });
	}
</script>

<div class="h-full flex flex-col pt-8">
	<!-- Project Header -->
	<header class="px-8 pb-12">
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

	<div class="px-8 flex items-center justify-between mb-16">
		<h3 class="text-3xl font-light text-surface-900 tracking-tight">Tasks</h3>
		<button 
			onclick={handleAddTask}
			class="bg-primary-500 hover:bg-primary-700 text-white px-6 py-2.5 rounded-xl font-medium transition-all active:scale-95 flex items-center gap-2"
		>
			<Plus size={18} />
			Add Task
		</button>
	</div>

	<!-- Board Content -->
	<div class="flex-1 overflow-x-auto custom-scrollbar-h">
		<div class="inline-flex h-full pl-3 pr-8 pb-8">
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
					<div
						class="flex-1 overflow-y-auto pr-2 space-y-4 custom-scrollbar rounded-xl transition-colors {dragOverColumn === column.id ? 'bg-primary-500/5 ring-2 ring-primary-500/20' : ''}"
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
								/>
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
				</div>
			{/each}
		</div>
	</div>
</div>

<AddTaskModal isOpen={isModalOpen} onClose={() => isModalOpen = false} />

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
</style>
