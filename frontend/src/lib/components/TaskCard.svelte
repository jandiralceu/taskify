<script lang="ts">
	import { Ellipsis, Eye, Pencil, ShieldBan, Trash2, Flag, MessageCircle, Paperclip, User } from '@lucide/svelte';
	import { Popover, Portal } from '@skeletonlabs/skeleton-svelte';
	import type { TaskResponse, TaskPriority } from '$lib/api/types';

	interface Props {
		task: TaskResponse;
		isDragging?: boolean;
		onDragStart: (e: DragEvent) => void;
		onDragEnd: () => void;
		onDelete?: (taskId: string) => void;
	}

	let { task, isDragging = false, onDragStart, onDragEnd, onDelete }: Props = $props();

	const priorityConfig: Record<TaskPriority, { label: string; class: string }> = {
		low:      { label: 'Low',      class: 'bg-indigo-50 text-indigo-700' },
		medium:   { label: 'Medium',   class: 'bg-blue-50 text-blue-700' },
		high:     { label: 'High',     class: 'bg-orange-50 text-orange-700' },
		critical: { label: 'Critical', class: 'bg-rose-50 text-rose-700' }
	};

	function formatDate(dateStr: string) {
		return new Intl.DateTimeFormat('en-GB', {
			day: '2-digit',
			month: 'short',
			year: 'numeric'
		}).format(new Date(dateStr));
	}

	const priority = $derived(priorityConfig[task.priority]);
</script>

<div
	role="listitem"
	draggable="true"
	ondragstart={onDragStart}
	ondragend={onDragEnd}
	class="bg-white rounded-2xl p-5 border border-gray-200 border-1 border-solid hover:shadow-md transition-all cursor-grab active:cursor-grabbing {isDragging ? 'opacity-40 scale-95' : ''}"
>
	<!-- Header: Priority + Menu -->
	<div class="flex justify-between items-center mb-3">
		<span class="inline-flex items-center px-3 py-1 rounded-full text-[12px] font-semibold {priority.class}">
			{priority.label}
		</span>
		<Popover positioning={{ placement: 'bottom-end' }}>
			<Popover.Trigger class="text-surface-400 hover:text-surface-600 transition-colors cursor-pointer p-1 rounded-md hover:bg-surface-100">
				<Ellipsis size={18} />
			</Popover.Trigger>
			<Portal>
				<Popover.Positioner>
					<Popover.Content class="bg-white rounded-xl shadow-lg border border-surface-200 py-1 w-44 z-50">
						<button class="w-full flex items-center gap-2.5 px-3 py-2 text-sm text-surface-700 hover:bg-surface-50 transition-colors cursor-pointer">
							<Eye size={15} class="text-surface-400" />
							<span>View Details</span>
						</button>
						<button class="w-full flex items-center gap-2.5 px-3 py-2 text-sm text-surface-700 hover:bg-surface-50 transition-colors cursor-pointer">
							<Pencil size={15} class="text-surface-400" />
							<span>Edit</span>
						</button>
						<button class="w-full flex items-center gap-2.5 px-3 py-2 text-sm text-surface-700 hover:bg-surface-50 transition-colors cursor-pointer">
							<ShieldBan size={15} class="text-surface-400" />
							<span>Block</span>
						</button>
						<hr class="my-1 border-surface-100" />
						<button
							onclick={() => onDelete?.(task.id)}
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
	<h4 class="font-medium leading-snug mb-2">{task.title}</h4>

	<!-- Assignees -->
	<div class="flex items-center justify-between mb-4">
		<span class="text-[13px] font-medium text-surface-500">Assignees :</span>
		<div class="flex -space-x-2">
			{#if task.assignee}
				{#if task.assignee.avatarUrl}
					<div 
						class="size-7 rounded-full bg-surface-100 border-2 border-white flex items-center justify-center overflow-hidden"
						title="{task.assignee.firstName} {task.assignee.lastName}"
					>
						<img src={task.assignee.avatarUrl} alt="{task.assignee.firstName}" class="size-full object-cover" />
					</div>
				{:else}
					<div 
						class="size-7 rounded-full bg-indigo-100 border-2 border-white flex items-center justify-center"
						title="{task.assignee.firstName} {task.assignee.lastName}"
					>
						<span class="text-[10px] font-bold text-indigo-700">
							{task.assignee.firstName[0]}{task.assignee.lastName[0]}
						</span>
					</div>
				{/if}
			{:else}
				<div class="size-7 rounded-full bg-surface-100 border-2 border-white flex items-center justify-center text-surface-400">
					<User size={14} />
				</div>
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

	<!-- Footer: Comments, Links, Files -->
	<div class="flex items-center justify-between pt-3 border-t border-surface-100 text-surface-700 text-xs">
		<div class="flex items-center gap-1.5">
			<Paperclip size={14} />
			<span>1 Attachments</span>
		</div>
		<div class="flex items-center gap-1.5">
			<MessageCircle size={14} />
			<span>12 Comments</span>
		</div>
	</div>
</div>
