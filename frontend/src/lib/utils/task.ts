import type { TaskPriority, TaskStatus } from '$lib/api/types';

interface LabelConfig {
	label: string;
	class: string;
}

export const priorityConfig: Record<TaskPriority, LabelConfig> = {
	low:      { label: 'Low',      class: 'bg-indigo-50 text-indigo-700' },
	medium:   { label: 'Medium',   class: 'bg-blue-50 text-blue-700' },
	high:     { label: 'High',     class: 'bg-orange-50 text-orange-700' },
	critical: { label: 'Critical', class: 'bg-rose-50 text-rose-700' }
};

export const statusConfig: Record<TaskStatus, LabelConfig> = {
	pending:     { label: 'Pending',     class: 'bg-slate-50 text-slate-700' },
	in_progress: { label: 'In Progress', class: 'bg-amber-50 text-amber-700' },
	completed:   { label: 'Completed',   class: 'bg-emerald-50 text-emerald-700' },
	cancelled:   { label: 'Cancelled',   class: 'bg-gray-50 text-gray-500' }
};
