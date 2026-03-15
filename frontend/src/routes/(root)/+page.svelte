<script lang="ts">
	import { Plus, Ellipsis, SquareCheck, LoaderCircle } from '@lucide/svelte';
	import { createProfileQuery } from '$lib/state/user.svelte';
	import { createTasksQuery, createUpdateTaskMutation } from '$lib/state/tasks.svelte';
	import AddTaskModal from '$lib/components/AddTaskModal.svelte';
	import type { TaskStatus } from '$lib/api/types';

	const profileQuery = createProfileQuery();
	const tasksQuery = createTasksQuery({ pageSize: 100 });
	const updateTaskMutation = createUpdateTaskMutation();

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

		const task = tasksQuery.data?.data.find(t => t.id === taskId);
		if (!task || task.status === targetStatus) return;

		await updateTaskMutation.mutateAsync({ id: taskId, data: { status: targetStatus } });
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
					<span class="font-light">here's a look at</span> <br />
					<span class="font-normal">your tasks for today!</span>
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
			class="bg-[#820AD1] hover:bg-[#6a08aa] text-white px-6 py-2.5 rounded-xl font-medium shadow-lg shadow-[#820AD1]/20 transition-all active:scale-95 flex items-center gap-2"
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
						class="flex-1 overflow-y-auto pr-2 space-y-4 custom-scrollbar rounded-xl transition-colors {dragOverColumn === column.id ? 'bg-[#820AD1]/5 ring-2 ring-[#820AD1]/20' : ''}"
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
						{:else if tasksQuery.data?.data}
							{#each tasksQuery.data.data.filter(t => t.status === column.id) as task (task.id)}
								<div
									role="listitem"
									draggable="true"
									ondragstart={(e) => onDragStart(e, task.id)}
									ondragend={onDragEnd}
									class="bg-white rounded-xl p-5 shadow-[0_2px_15px_-3px_rgba(0,0,0,0.07),0_10px_20px_-2px_rgba(0,0,0,0.04)] hover:shadow-lg transition-all cursor-grab active:cursor-grabbing group border border-surface-50/50 {draggingTaskId === task.id ? 'opacity-40 scale-95' : ''}"
								>
									<div class="flex justify-between items-start mb-3">
										<h4 class="text-[15px] font-bold text-surface-800 leading-snug">{task.title}</h4>
									</div>

									{#if task.description}
										<p class="text-xs text-surface-400 leading-relaxed mb-4 line-clamp-3">
											{task.description}
										</p>
									{/if}

									<div class="flex items-center justify-between mt-auto pt-2">
										<div class="flex gap-2">
											<span class="px-2.5 py-1 bg-surface-50 text-surface-600 rounded text-[10px] font-bold border border-surface-100 uppercase tracking-wider">
												{task.priority}
											</span>
										</div>

										<div class="flex items-center gap-3 text-surface-300">
											{#if task.dueDate}
												<div class="flex items-center gap-1">
													<SquareCheck size={13} strokeWidth={2.5} />
													<span class="text-[10px] font-bold">1</span>
												</div>
											{/if}
										</div>
									</div>
								</div>
							{/each}
						{/if}

						<!-- Plus Button only on Pending column -->
						{#if column.id === 'pending'}
							<button
								onclick={handleAddTask}
								class="w-full py-3 rounded-xl border-2 border-dashed border-[#820AD1]/30 text-[#820AD1]/60 hover:text-[#820AD1] hover:border-[#820AD1] hover:bg-[#820AD1]/5 transition-all flex items-center justify-center gap-2 font-medium text-sm"
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
