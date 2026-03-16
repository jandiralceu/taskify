<script lang="ts">
	import { X, Calendar, Flag, Clock, User, ShieldBan, CheckCircle, AlertCircle } from '@lucide/svelte';
	import { Dialog, Portal } from '@skeletonlabs/skeleton-svelte';
	import type { TaskResponse } from '$lib/api/types';
	import { priorityConfig, statusConfig } from '$lib/utils/task';

	interface Props {
		task: TaskResponse | null;
		isOpen: boolean;
		onClose: () => void;
	}

	let { task, isOpen, onClose }: Props = $props();

	function handleOpenChange(e: { open: boolean }) {
		if (!e.open) onClose();
	}

	function formatDate(dateStr: string) {
		return new Intl.DateTimeFormat('en-GB', {
			day: 'numeric',
			month: 'long',
			year: 'numeric',
			timeZone: 'UTC'
		}).format(new Date(dateStr));
	}

	function formatDateTime(dateStr: string) {
		return new Intl.DateTimeFormat('en-GB', {
			day: 'numeric',
			month: 'long',
			year: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		}).format(new Date(dateStr));
	}

	function formatHours(hours: number) {
		if (hours === 1) return '1 hour';
		return `${hours} hours`;
	}

	let priority = $derived(task ? priorityConfig[task.priority] : null);
	let status = $derived(task ? statusConfig[task.status] : null);
</script>

<Dialog open={isOpen} onOpenChange={handleOpenChange} closeOnInteractOutside={true}>
	<Portal>
		<Dialog.Backdrop class="fixed inset-0 z-50 bg-surface-950/30 backdrop-blur-[2px] transition-opacity" />
		<Dialog.Positioner class="fixed inset-y-0 right-0 z-50 flex items-stretch justify-end">
			<Dialog.Content class="w-full max-w-lg bg-white shadow-2xl overflow-y-auto animate-slide-in">
				{#if task}
					<!-- Header -->
					<div class="sticky top-0 bg-white/95 backdrop-blur-sm z-10 px-8 pt-8 pb-4 border-b border-surface-100">
						<div class="flex items-start justify-between gap-4">
							<div class="flex-1 min-w-0">
								<Dialog.Title class="text-2xl font-medium text-surface-900 tracking-tight leading-tight">
									{task.title}
								</Dialog.Title>
								<Dialog.Description class="text-sm text-surface-500 mt-1">
									Created {formatDateTime(task.createdAt)}
								</Dialog.Description>
							</div>
							<Dialog.CloseTrigger
								class="p-2 -mr-2 -mt-1 rounded-xl text-surface-400 hover:text-surface-900 hover:bg-surface-100 transition-all shrink-0"
							>
								<X size={20} />
							</Dialog.CloseTrigger>
						</div>
					</div>

					<!-- Body -->
					<div class="px-8 py-6 space-y-6">
						<!-- Metadata Grid -->
						<div class="space-y-4">
							<!-- Status -->
							<div class="flex items-center justify-between py-3 border-b border-surface-50">
								<div class="flex items-center gap-3 text-sm text-surface-600">
									<CheckCircle size={18} strokeWidth={2.5} />
									<span class="font-medium">Status</span>
								</div>
								{#if status}
									<span class="inline-flex items-center px-3 py-1 rounded-full text-xs font-semibold {status.class}">
										{status.label}
									</span>
								{/if}
							</div>

							<!-- Priority -->
							<div class="flex items-center justify-between py-3 border-b border-surface-50">
								<div class="flex items-center gap-3 text-sm text-surface-600">
									<Flag size={18} strokeWidth={2.5} />
									<span class="font-medium">Priority</span>
								</div>
								{#if priority}
									<span class="inline-flex items-center px-3 py-1 rounded-full text-xs font-semibold {priority.class}">
										{priority.label}
									</span>
								{/if}
							</div>

							<!-- Blocked -->
							{#if task.isBlocked}
								<div class="flex items-center justify-between py-3 border-b border-surface-50">
									<div class="flex items-center gap-2.5 text-sm text-surface-500">
										<ShieldBan size={16} />
										<span class="font-medium">Blocked</span>
									</div>
									<span class="inline-flex items-center gap-1 px-3 py-1 rounded-full text-xs font-semibold bg-rose-50 text-rose-600">
										<AlertCircle size={12} />
										Yes
									</span>
								</div>
							{/if}

							<!-- Due Date -->
							<div class="flex items-center justify-between py-3 border-b border-surface-50">
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
								<div class="flex items-center justify-between py-3 border-b border-surface-50">
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
								<div class="flex items-center justify-between py-3 border-b border-surface-50">
									<div class="flex items-center gap-2.5 text-sm text-surface-500">
										<Clock size={16} />
										<span class="font-medium">Actual</span>
									</div>
									<span class="text-sm font-medium text-surface-900">
										{formatHours(task.actualHours)}
									</span>
								</div>
							{/if}

							<!-- Assignee -->
							<div class="flex items-center justify-between py-3 border-b border-surface-50">
								<div class="flex items-center gap-2.5 text-sm text-surface-500">
									<User size={16} />
									<span class="font-medium">Assignee</span>
								</div>
								<div class="flex items-center gap-2">
									{#if task.assignee.avatarUrl}
										<div class="size-7 rounded-full overflow-hidden bg-surface-100 border-2 border-white shadow-sm">
											<img src={task.assignee.avatarUrl} alt="{task.assignee.firstName}" class="size-full object-cover" />
										</div>
									{:else}
										<div class="size-7 rounded-full bg-indigo-100 border-2 border-white shadow-sm flex items-center justify-center">
											<span class="text-[10px] font-bold text-indigo-700">
												{task.assignee.firstName[0]}{task.assignee.lastName[0]}
											</span>
										</div>
									{/if}
									<span class="text-sm font-medium text-surface-900">
										{task.assignee.firstName} {task.assignee.lastName}
									</span>
								</div>
							</div>
						</div>

						<!-- Description -->
						{#if task.description}
							<div class="space-y-3">
								<h4 class="text-sm font-semibold text-surface-900 uppercase tracking-wider">Description</h4>
								<div class="bg-surface-50 rounded-2xl p-5 border border-surface-100">
									<p class="text-sm text-surface-700 leading-relaxed whitespace-pre-wrap">{task.description}</p>
								</div>
							</div>
						{/if}

						<!-- Timestamps -->
						<div class="space-y-3 pt-4 border-t border-surface-100">
							<h4 class="text-sm font-semibold text-surface-900 uppercase tracking-wider">Timeline</h4>
							<div class="space-y-2">
								<div class="flex items-center justify-between text-xs">
									<span class="text-surface-500">Created</span>
									<span class="text-surface-700 font-medium">{formatDateTime(task.createdAt)}</span>
								</div>
								<div class="flex items-center justify-between text-xs">
									<span class="text-surface-500">Last Updated</span>
									<span class="text-surface-700 font-medium">{formatDateTime(task.updatedAt)}</span>
								</div>
								{#if task.completedAt}
									<div class="flex items-center justify-between text-xs">
										<span class="text-surface-500">Completed</span>
										<span class="text-emerald-600 font-medium">{formatDateTime(task.completedAt)}</span>
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
