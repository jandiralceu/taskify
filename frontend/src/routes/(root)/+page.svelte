<script lang="ts">
	import { Plus, ChevronDown, Ellipsis, Link2, Paperclip, MessageSquare } from '@lucide/svelte';

	const columns = [
		{ id: 'pending', title: 'Pending' },
		{ id: 'in_progress', title: 'In Progress' },
		{ id: 'completed', title: 'Completed' },
		{ id: 'cancelled', title: 'Cancelled' }
	];

	const tasks = [
		{ 
			id: 1, 
			col: 'pending', 
			title: 'Lorem Ipsum', 
			date: '04/09/2020', 
			priority: 'low',
			tags: ['Development', 'Legal', 'Tag Name']
		},
		{ 
			id: 2, 
			col: 'pending', 
			title: 'Cras sed sem lacinia', 
			date: '04/09/2020', 
			priority: 'medium',
			tags: ['Development', 'Legal']
		},
		{ 
			id: 3, 
			col: 'in_progress', 
			title: 'Curabitur Varius', 
			date: '04/09/2020', 
			priority: 'high',
			tags: ['Development', 'Legal', 'Tag Name']
		},
		{ 
			id: 4, 
			col: 'completed', 
			title: 'Morbi quis venenatis', 
			date: '04/09/2020', 
			priority: 'critical',
			tags: ['Development', 'Legal']
		}
	];

	function getPriorityClass(p: string) {
		switch (p.toLowerCase()) {
			case 'low': return 'bg-sky-50 text-sky-500';
			case 'medium': return 'bg-orange-50 text-orange-500';
			case 'high': return 'bg-rose-50 text-rose-500';
			case 'critical': return 'bg-red-100 text-red-700 border border-red-200';
			default: return 'bg-slate-50 text-slate-400';
		}
	}
</script>

<div class="space-y-16 h-full flex flex-col pt-16">
	<!-- Project Header -->
	<header class="space-y-6">
		<div class="flex items-end justify-between">
			<div class="space-y-4">
				<h1 class="text-2xl font-black text-surface-900 tracking-tight">Tasks</h1>
				
				<!-- Progress Bar Section -->
				<div class="flex items-center gap-6">
					<div class="w-[400px] h-2 bg-surface-200 rounded-full overflow-hidden flex shadow-inner">
						<div class="h-full bg-primary-500" style="width: 35%"></div>
						<div class="h-full bg-secondary-500" style="width: 25%"></div>
						<div class="h-full bg-tertiary-500" style="width: 20%"></div>
						<div class="h-full bg-surface-400" style="width: 20%"></div>
					</div>
					<span class="text-sm font-bold text-surface-400">12 Tasks</span>
				</div>

				<!-- Status Indicators -->
				<div class="flex gap-6">
					{#each columns as col, i (col.id)}
						<div class="flex items-center gap-2">
							<div class="w-4 h-4 rounded-full flex items-center justify-center text-[10px] text-white font-bold {['bg-primary-500', 'bg-secondary-500', 'bg-tertiary-500', 'bg-surface-400'][i]}">
								{i+1}
							</div>
							<span class="text-xs font-bold text-surface-400">{col.title}</span>
						</div>
					{/each}
				</div>
			</div>

			<!-- Filters -->
			<div class="flex items-center gap-6">
				<button class="flex items-center gap-2 text-xs font-bold text-surface-400 hover:text-primary-600 transition-colors">
					Due Date: <span class="text-surface-900">Today</span>
					<ChevronDown size={16} />
				</button>
				<button class="flex items-center gap-2 text-xs font-bold text-surface-400 hover:text-primary-600 transition-colors">
					Priority: <span class="text-surface-900">Any</span>
					<ChevronDown size={16} />
				</button>
				<button class="flex items-center gap-2 text-xs font-bold text-surface-400 hover:text-primary-600 transition-colors">
					Tags: <span class="text-surface-900">All</span>
					<ChevronDown size={16} />
				</button>
			</div>
		</div>
	</header>

	<!-- Board Content -->
	<div class="flex-1 flex gap-6 overflow-hidden">
		{#each columns as column (column.id)}
			<div class="w-[320px] flex flex-col gap-4">
				<!-- Column Header -->
				<div class="flex flex-col gap-2">
					<h3 class="text-lg font-black text-surface-900 px-2">{column.title}</h3>
				</div>

				<!-- Cards Area -->
				<div class="flex-1 overflow-y-auto pr-2 space-y-4 custom-scrollbar">
					{#each tasks.filter(t => t.col === column.id) as task (task.id)}
						<div class="bg-white rounded-2xl p-6 shadow-sm border border-surface-100 hover:shadow-md transition-all cursor-pointer group relative">
							<div class="flex justify-between items-start mb-4">
								<h4 class="text-lg font-black text-surface-900 leading-tight group-hover:text-primary-600 transition-colors">{task.title}</h4>
								<button class="text-surface-300 opacity-0 group-hover:opacity-100 transition-opacity">
									<Ellipsis size={18} />
								</button>
							</div>

							<div class="flex items-center gap-2 text-xs text-surface-400 mb-4 font-bold">
								<Link2 size={14} />
								{task.date}
							</div>

							<!-- Priority Badge -->
							<div class="mb-4">
								<span class="px-3 py-1 rounded-lg text-[10px] font-black uppercase tracking-wider {getPriorityClass(task.priority)}">
									{task.priority}
								</span>
							</div>

							<!-- Tags -->
							<div class="flex flex-wrap gap-2 mb-6">
								{#each task.tags as tag (tag)}
									<span class="px-3 py-1 bg-primary-50 text-primary-600 rounded-lg text-[10px] font-bold">
										{tag}
									</span>
								{/each}
							</div>

							<!-- Card Footer -->
							<div class="flex justify-between items-center pt-4 border-t border-surface-50">
								<div class="flex items-center gap-3 text-surface-300">
									<div class="flex items-center gap-1">
										<Paperclip size={14} />
										<span class="text-[10px] font-bold">3</span>
									</div>
									<div class="flex items-center gap-1">
										<MessageSquare size={14} />
										<span class="text-[10px] font-bold">4</span>
									</div>
								</div>

								<div class="flex -space-x-2">
									{#each [1, 2] as j (j)}
										<div class="w-7 h-7 rounded-full border-2 border-white overflow-hidden bg-surface-200">
											<img src={`https://api.dicebear.com/7.x/avataaars/svg?seed=${task.id + j + 100}`} alt="user" />
										</div>
									{/each}
									<button class="w-7 h-7 rounded-full border-2 border-white border-dashed bg-surface-50 flex items-center justify-center text-surface-300 hover:text-primary-600 transition-colors">
										<Plus size={12} />
									</button>
								</div>
							</div>
						</div>
					{/each}

					<!-- Sample Placeholder for empty look -->
					{#if column.id === 'in_progress' && tasks.filter(t => t.col === column.id).length === 1}
						<div class="h-40 bg-surface-100 border-2 border-dashed border-surface-200 rounded-2xl"></div>
					{/if}
				</div>
			</div>
		{/each}


	</div>
</div>

<style>
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
</style>
