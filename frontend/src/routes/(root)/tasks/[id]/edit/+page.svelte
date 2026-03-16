<script lang="ts">
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { ArrowLeft, LoaderCircle, Flag, Clock, User } from '@lucide/svelte';
	import { getTaskQuery, updateTaskMutation } from '$lib/state/tasks.svelte';
	import { getUsersQuery } from '$lib/state/user.svelte';
	import { toaster } from '$lib/state/toast.svelte';
	import Input from '$lib/components/Input.svelte';
	import type { TaskStatus, TaskPriority } from '$lib/api/types';

	const taskId = $derived(page.params.id ?? '');

	const taskQuery = getTaskQuery(() => taskId);
	const updateTask = updateTaskMutation();

	const usersQuery = getUsersQuery(() => ({ limit: 100 }));

	// Editable fields
	let title = $state('');
	let description = $state('');
	let status = $state<TaskStatus>('pending');
	let priority = $state<TaskPriority>('medium');
	let assignedTo = $state('');
	let dueDate = $state('');
	let estimatedHours = $state<number | ''>('');
	let actualHours = $state<number | ''>('');
	let isBlocked = $state(false);
	let isArchived = $state(false);

	$effect(() => {
		if (taskQuery.data) {
			const t = taskQuery.data;
			title = t.title;
			description = t.description ?? '';
			status = t.status;
			priority = t.priority;
			assignedTo = t.assignedTo ?? '';
			dueDate = t.dueDate ? t.dueDate.split('T')[0] : '';
			estimatedHours = t.estimatedHours ?? '';
			actualHours = t.actualHours ?? '';
			isBlocked = t.isBlocked;
			isArchived = t.isArchived;
		}
	});

	const statusOptions: { value: TaskStatus; label: string }[] = [
		{ value: 'pending', label: 'Pending' },
		{ value: 'in_progress', label: 'In Progress' },
		{ value: 'completed', label: 'Completed' },
		{ value: 'cancelled', label: 'Cancelled' }
	];

	const priorityOptions: { value: TaskPriority; label: string }[] = [
		{ value: 'low', label: 'Low' },
		{ value: 'medium', label: 'Medium' },
		{ value: 'high', label: 'High' },
		{ value: 'critical', label: 'Critical' }
	];

	function formatDate(dateStr: string) {
		return new Intl.DateTimeFormat('en-GB', {
			day: '2-digit',
			month: 'long',
			year: 'numeric'
		}).format(new Date(dateStr));
	}

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault();

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
					estimatedHours: estimatedHours !== '' ? Number(estimatedHours) : undefined,
					actualHours: actualHours !== '' ? Number(actualHours) : undefined,
					isBlocked,
					isArchived
				}
			});

			toaster.success({
				title: 'Task Updated',
				description: `"${title}" has been updated successfully.`
			});

			goto(resolve('/'));
		} catch {
			toaster.error({
				title: 'Update Failed',
				description: 'Could not update the task. Please try again.'
			});
		}
	}
</script>

<svelte:head>
	<title>Edit Task - Taskify</title>
</svelte:head>

<div class="h-full flex flex-col pt-8">
	<header class="px-8 pb-8">
		<button
			onclick={() => goto(resolve('/'))}
			class="flex items-center gap-2 text-sm text-surface-500 hover:text-surface-800 transition-colors mb-6"
		>
			<ArrowLeft size={16} />
			Back to Board
		</button>

		<h2 class="text-4xl text-surface-900 tracking-tight leading-tight">
			<span class="font-light">Edit</span>
			<span class="font-normal"> Task</span>
		</h2>
	</header>

	<div class="px-8 flex-1 overflow-y-auto pb-8">
		{#if taskQuery.isPending}
			<div class="flex flex-col items-center justify-center py-24 text-surface-400">
				<LoaderCircle size={28} class="animate-spin mb-3" />
				<span class="text-sm font-medium">Loading task...</span>
			</div>
		{:else if taskQuery.isError}
			<div class="flex flex-col items-center justify-center py-24 text-rose-500">
				<p class="text-sm font-medium">Failed to load task.</p>
			</div>
		{:else if taskQuery.data}
			{@const task = taskQuery.data}
			<div class="max-w-2xl space-y-6">

				<!-- Read-only meta card -->
				<div class="bg-white rounded-2xl border border-surface-200 p-6 grid grid-cols-1 sm:grid-cols-3 gap-4">
					<!-- Created by -->
					<div class="space-y-1">
						<p class="text-xs font-semibold text-surface-400 uppercase tracking-wide">Created by</p>
						<div class="flex items-center gap-2">
							{#if task.assignee.avatarUrl}
								<img src={task.assignee.avatarUrl} alt="" class="size-6 rounded-full object-cover" />
							{:else}
								<div class="size-6 rounded-full bg-indigo-100 flex items-center justify-center">
									<span class="text-[9px] font-bold text-indigo-700">
										{task.assignee.firstName[0]}{task.assignee.lastName[0]}
									</span>
								</div>
							{/if}
							<span class="text-sm text-surface-700 font-medium">
								{task.assignee.firstName} {task.assignee.lastName}
							</span>
						</div>
					</div>

					<!-- Created at -->
					<div class="space-y-1">
						<p class="text-xs font-semibold text-surface-400 uppercase tracking-wide">Created</p>
						<p class="text-sm text-surface-700">{formatDate(task.createdAt)}</p>
					</div>

					<!-- Updated at -->
					<div class="space-y-1">
						<p class="text-xs font-semibold text-surface-400 uppercase tracking-wide">Last updated</p>
						<p class="text-sm text-surface-700">{formatDate(task.updatedAt)}</p>
					</div>
				</div>

				<!-- Edit form -->
				<form onsubmit={handleSubmit} class="bg-white rounded-2xl border border-surface-200 p-6 space-y-5">
					<h3 class="text-sm font-semibold text-surface-900">Task Details</h3>

					<!-- Title -->
					<Input id="title" label="Title" placeholder="Task title" bind:value={title} required />

					<!-- Description -->
					<div class="space-y-1">
						<label for="description" class="block text-sm font-medium text-surface-700">
							Description
						</label>
						<textarea
							id="description"
							rows={4}
							placeholder="Task description..."
							bind:value={description}
							class="block w-full rounded-xl border border-surface-300 bg-surface-50 px-4 py-3 text-sm text-surface-900 placeholder-surface-500 focus:outline-none focus:border-primary-500 focus:ring-2 focus:ring-primary-500/20 transition-all resize-none"
						></textarea>
					</div>

					<!-- Status + Priority -->
					<div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
						<div class="space-y-1">
							<label for="status" class="block text-sm font-medium text-surface-700">Status</label>
							<select
								id="status"
								bind:value={status}
								class="block h-12 w-full rounded-xl border border-surface-300 bg-surface-50 px-4 text-sm text-surface-900 focus:outline-none focus:border-primary-500 focus:ring-2 focus:ring-primary-500/20 transition-all"
							>
								{#each statusOptions as opt}
									<option value={opt.value}>{opt.label}</option>
								{/each}
							</select>
						</div>

						<div class="space-y-1">
							<label for="priority" class="block text-sm font-medium text-surface-700">Priority</label>
							<select
								id="priority"
								bind:value={priority}
								class="block h-12 w-full rounded-xl border border-surface-300 bg-surface-50 px-4 text-sm text-surface-900 focus:outline-none focus:border-primary-500 focus:ring-2 focus:ring-primary-500/20 transition-all"
							>
								{#each priorityOptions as opt}
									<option value={opt.value}>{opt.label}</option>
								{/each}
							</select>
						</div>
					</div>

					<!-- Assigned To -->
					<div class="space-y-1">
						<label for="assignedTo" class="block text-sm font-medium text-surface-700">
							<span class="flex items-center gap-1.5"><User size={14} /> Assigned To</span>
						</label>
						<select
							id="assignedTo"
							bind:value={assignedTo}
							class="block h-12 w-full rounded-xl border border-surface-300 bg-surface-50 px-4 text-sm text-surface-900 focus:outline-none focus:border-primary-500 focus:ring-2 focus:ring-primary-500/20 transition-all"
						>
							<option value="">Unassigned</option>
							{#if usersQuery.data?.data}
								{#each usersQuery.data.data as user}
									<option value={user.id}>{user.firstName} {user.lastName}</option>
								{/each}
							{/if}
						</select>
					</div>

					<!-- Due Date + Estimated Hours -->
					<div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
						<div class="space-y-1">
							<label for="dueDate" class="block text-sm font-medium text-surface-700">
								<span class="flex items-center gap-1.5"><Flag size={14} /> Due Date</span>
							</label>
							<input
								id="dueDate"
								type="date"
								bind:value={dueDate}
								class="block h-12 w-full rounded-xl border border-surface-300 bg-surface-50 px-4 text-sm text-surface-900 focus:outline-none focus:border-primary-500 focus:ring-2 focus:ring-primary-500/20 transition-all"
							/>
						</div>

						<div class="space-y-1">
							<label for="estimatedHours" class="block text-sm font-medium text-surface-700">
								<span class="flex items-center gap-1.5"><Clock size={14} /> Estimated Hours</span>
							</label>
							<input
								id="estimatedHours"
								type="number"
								min="0"
								step="0.5"
								placeholder="0"
								bind:value={estimatedHours}
								class="block h-12 w-full rounded-xl border border-surface-300 bg-surface-50 px-4 text-sm text-surface-900 focus:outline-none focus:border-primary-500 focus:ring-2 focus:ring-primary-500/20 transition-all"
							/>
						</div>
					</div>

					<!-- Actual Hours -->
					<div class="space-y-1">
						<label for="actualHours" class="block text-sm font-medium text-surface-700">
							<span class="flex items-center gap-1.5"><Clock size={14} /> Actual Hours</span>
						</label>
						<input
							id="actualHours"
							type="number"
							min="0"
							step="0.5"
							placeholder="0"
							bind:value={actualHours}
							class="block h-12 w-full rounded-xl border border-surface-300 bg-surface-50 px-4 text-sm text-surface-900 focus:outline-none focus:border-primary-500 focus:ring-2 focus:ring-primary-500/20 transition-all"
						/>
					</div>

					<hr class="border-surface-100" />

					<!-- Toggles: Blocked + Archived -->
					<div class="space-y-3">
						<!-- Blocked -->
						<div class="flex items-center justify-between p-4 rounded-xl border border-surface-200 bg-surface-50">
							<div>
								<p class="text-sm font-medium text-surface-800">Blocked</p>
								<p class="text-xs text-surface-500 mt-0.5">Blocked tasks cannot be moved on the board</p>
							</div>
							<button
								type="button"
								aria-label="Toggle blocked status"
								role="switch"
								aria-checked={isBlocked}
								onclick={() => (isBlocked = !isBlocked)}
								class="relative inline-flex h-6 w-11 shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 focus:outline-none {isBlocked ? 'bg-rose-500' : 'bg-surface-300'}"
							>
								<span
									class="pointer-events-none inline-block size-5 rounded-full bg-white shadow-sm transform transition duration-200 {isBlocked ? 'translate-x-5' : 'translate-x-0'}"
								></span>
							</button>
						</div>

						<!-- Archived -->
						<div class="flex items-center justify-between p-4 rounded-xl border border-surface-200 bg-surface-50">
							<div>
								<p class="text-sm font-medium text-surface-800">Archived</p>
								<p class="text-xs text-surface-500 mt-0.5">Archived tasks are hidden from the board</p>
							</div>
							<button
								type="button"
								aria-label="Toggle archived status"
								role="switch"
								aria-checked={isArchived}
								onclick={() => (isArchived = !isArchived)}
								class="relative inline-flex h-6 w-11 shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 focus:outline-none {isArchived ? 'bg-primary-500' : 'bg-surface-300'}"
							>
								<span
									class="pointer-events-none inline-block size-5 rounded-full bg-white shadow-sm transform transition duration-200 {isArchived ? 'translate-x-5' : 'translate-x-0'}"
								></span>
							</button>
						</div>
					</div>

					<!-- Footer -->
					<div class="flex items-center justify-end gap-3 pt-2">
						<button
							type="button"
							onclick={() => goto(resolve('/'))}
							class="px-6 py-2.5 rounded-xl font-medium text-surface-500 hover:text-surface-900 hover:bg-surface-50 transition-all text-sm"
						>
							Cancel
						</button>
						<button
							type="submit"
							disabled={updateTask.isPending}
							class="bg-primary-500 hover:bg-primary-700 text-white px-8 py-2.5 rounded-xl font-bold shadow-lg shadow-primary-500/20 transition-all active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2 text-sm"
						>
							{#if updateTask.isPending}
								<LoaderCircle size={16} class="animate-spin" />
								Saving...
							{:else}
								Save Changes
							{/if}
						</button>
					</div>
				</form>
			</div>
		{/if}
	</div>
</div>
