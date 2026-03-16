<script lang="ts">
	import { Ellipsis, Pencil, ShieldBan, ShieldCheck, Trash2, Flag, MessageCircle, Paperclip } from '@lucide/svelte';
	import { Popover, Portal } from '@skeletonlabs/skeleton-svelte';
	import type { TaskResponse } from '$lib/api/types';
	import { priorityConfig } from '$lib/utils/task';

	interface Props {
		task: TaskResponse;
		isDragging?: boolean;
		onDragStart: (e: DragEvent) => void;
		onDragEnd: () => void;
		onDelete?: (taskId: string) => void;
		onToggleBlock?: (taskId: string, blocked: boolean) => void;
		onViewDetails?: (task: TaskResponse) => void;
		onViewUser?: (userId: string) => void;
	}

	let { task, isDragging = false, onDragStart, onDragEnd, onDelete, onToggleBlock, onViewDetails, onViewUser }: Props = $props();
	let isMenuOpen = $state(false);

	function formatDate(dateStr: string) {
		return new Intl.DateTimeFormat('en-GB', {
			day: '2-digit',
			month: 'short',
			year: 'numeric'
		}).format(new Date(dateStr));
	}

	function handleDragStart(e: DragEvent) {
		if (task.isBlocked) {
			e.preventDefault();
			return;
		}
		onDragStart(e);
	}

	const priority = $derived(priorityConfig[task.priority]);
</script>

<div
	role="button"
	tabindex="0"
	draggable={!task.isBlocked}
	ondragstart={handleDragStart}
	ondragend={onDragEnd}
	onclick={() => onViewDetails?.(task)}
	onkeydown={(e) => { if (e.key === 'Enter' || e.key === ' ') { e.preventDefault(); onViewDetails?.(task); } }}
	class="bg-white rounded-2xl p-5 border border-1 border-solid transition-all {task.isBlocked ? 'border-rose-200 cursor-not-allowed opacity-60 bg-surface-50' : 'border-gray-200 hover:shadow-md cursor-pointer active:cursor-grabbing'} {isDragging ? 'opacity-40 scale-95' : ''}"
>
	<!-- Header: Priority + Blocked Badge + Menu -->
	<div class="flex justify-between items-center mb-3">
		<div class="flex items-center gap-1.5">
			<span class="inline-flex items-center px-3 py-1 rounded-full text-[12px] font-semibold {priority.class}">
				{priority.label}
			</span>
			{#if task.isBlocked}
				<span class="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-[12px] font-semibold bg-rose-50 text-rose-600">
					<ShieldBan size={12} />
					Blocked
				</span>
			{/if}
		</div>
		<Popover
			open={isMenuOpen}
			onOpenChange={(e) => (isMenuOpen = e.open)}
			positioning={{ placement: 'bottom-end' }}
		>
			<Popover.Trigger 
				onclick={(e) => e.stopPropagation()}
				class="text-surface-400 hover:text-surface-600 transition-colors cursor-pointer p-1 rounded-md hover:bg-surface-100"
			>
				<Ellipsis size={18} />
			</Popover.Trigger>
			<Portal>
				<Popover.Positioner>
					<Popover.Content class="bg-white rounded-xl shadow-lg border border-surface-200 py-1 w-44 z-50">
						<button
							onclick={(e) => { e.stopPropagation(); isMenuOpen = false; }}
							class="w-full flex items-center gap-2.5 px-3 py-2 text-sm text-surface-700 hover:bg-surface-50 transition-colors cursor-pointer"
						>
							<Pencil size={15} class="text-surface-400" />
							<span>Edit</span>
						</button>
						<button
							onclick={(e) => {
								e.stopPropagation();
								onToggleBlock?.(task.id, !task.isBlocked);
								isMenuOpen = false;
							}}
							class="w-full flex items-center gap-2.5 px-3 py-2 text-sm text-surface-700 hover:bg-surface-50 transition-colors cursor-pointer"
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
							onclick={(e) => {
								e.stopPropagation();
								onDelete?.(task.id);
								isMenuOpen = false;
							}}
							class="w-full flex items-center gap-2.5 px-3 py-2 text-sm text-rose-600 hover:bg-rose-50 transition-colors cursor-pointer"
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
	<h4 class="font-light leading-snug mb-2">{task.title}</h4>

	<!-- Assignees -->
	<div class="flex items-center justify-between mb-4">
		<span class="text-xs font-medium text-surface-700">Assignees :</span>
		<div class="flex -space-x-2">
			{#if task.assignee.avatarUrl}
				<button 
					type="button"
					onclick={(e) => { e.stopPropagation(); onViewUser?.(task.assignee.id); }}
					class="size-7 rounded-full bg-surface-100 border-2 border-white flex items-center justify-center overflow-hidden cursor-pointer hover:ring-2 hover:ring-primary-500 transition-all"
					title="{task.assignee.firstName} {task.assignee.lastName}"
				>
					<img src={task.assignee.avatarUrl} alt="{task.assignee.firstName}" class="size-full object-cover" />
				</button>
			{:else}
				<button 
					type="button"
					onclick={(e) => { e.stopPropagation(); onViewUser?.(task.assignee.id); }}
					class="size-7 rounded-full bg-indigo-100 border-2 border-white flex items-center justify-center cursor-pointer hover:ring-2 hover:ring-primary-500 transition-all"
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
	<div class="flex items-center mb-4">
		{#if task.dueDate}
			<div class="flex items-center gap-2 text-xs font-medium text-surface-700">
				<Flag size={14} />
				<span>{formatDate(task.dueDate)}</span>
			</div>
		{/if}
	</div>

	<!-- Footer: Comments, Files -->
	<div class="flex items-center justify-between pt-3 border-t border-surface-100 text-surface-500 text-xs">
		<div class="flex items-center gap-1.5 ">
			<Paperclip size={14} />
			<span>{task.attachmentsCount} {task.attachmentsCount === 1 ? 'Attachment' : 'Attachments'}</span>
		</div>
		<div class="flex items-center gap-1.5">
			<MessageCircle size={14} />
			<span>{task.notesCount} {task.notesCount === 1 ? 'Comment' : 'Comments'}</span>
		</div>
	</div>
</div>
