<script lang="ts">
	import { Ellipsis, Flag, File, MessageCircle, Link } from '@lucide/svelte';
	import type { TaskResponse, TaskPriority } from '$lib/api/types';

	interface Props {
		task: TaskResponse;
		isDragging?: boolean;
		onDragStart: (e: DragEvent) => void;
		onDragEnd: () => void;
	}

	let { task, isDragging = false, onDragStart, onDragEnd }: Props = $props();

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
	class="bg-white rounded-2xl p-5 border border-surface-200 shadow-sm hover:shadow-md transition-all cursor-grab active:cursor-grabbing {isDragging ? 'opacity-40 scale-95' : ''}"
>
	<!-- Header: Priority + Menu -->
	<div class="flex justify-between items-center mb-3">
		<span class="inline-flex items-center px-3 py-1 rounded-full text-[12px] font-semibold {priority.class}">
			{priority.label}
		</span>
		<button class="text-surface-400 hover:text-surface-600 transition-colors">
			<Ellipsis size={18} />
		</button>
	</div>

	<!-- Title -->
	<h4 class="text-[16px] font-bold text-surface-900 leading-snug mb-2">{task.title}</h4>

	<!-- Description -->
	{#if task.description}
		<p class="text-[13px] text-surface-500 leading-relaxed mb-4 line-clamp-2">
			{task.description}
		</p>
	{/if}

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
					<span class="text-[10px] font-bold">?</span>
				</div>
			{/if}
		</div>
	</div>

	<!-- Date -->
	<div class="flex items-center mb-4">
		{#if task.dueDate}
			<div class="flex items-center gap-2 text-[13px] font-medium text-surface-500">
				<Flag size={14} />
				<span>{formatDate(task.dueDate)}</span>
			</div>
		{/if}
	</div>

	<!-- Footer: Comments, Links, Files -->
	<div class="flex items-center gap-4 pt-3 border-t border-surface-100 text-surface-500">
		<div class="flex items-center gap-1.5 text-[13px]">
			<MessageCircle size={14} />
			<span>12 Comments</span>
		</div>
		<div class="flex items-center gap-1.5 text-[13px]">
			<Link size={14} />
			<span>1 Links</span>
		</div>
		<div class="flex items-center gap-1.5 text-[13px]">
			<File size={14} />
			<span>0/3</span>
		</div>
	</div>
</div>
