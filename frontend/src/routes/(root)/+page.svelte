<script lang="ts">
	import { Plus, Ellipsis, MessageSquare, SquareCheck } from '@lucide/svelte';
	import { createProfileQuery } from '$lib/state/user.svelte';

	const profileQuery = createProfileQuery();

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
			title: '3.4. Simplicity', 
			description: 'One of the key elements to consider when designing is simplicity...',
			progress: 0,
			comments: 4,
			subtasksTotal: 0,
			subtasksCompleted: 0
		},
		{ 
			id: 2, 
			col: 'pending', 
			title: '3.5. Consistency', 
			progress: 0,
			comments: 2
		},
		{ 
			id: 3, 
			col: 'in_progress', 
			title: '3.1. Design Research', 
			progress: 80,
			subtasksTotal: 6,
			subtasksCompleted: 5,
			tag: 'after review',
			comments: 2
		},
		{ 
			id: 4, 
			col: 'in_progress', 
			title: '3.2. Content 🍕', 
			description: 'There are benefits to starting with no content, or even fake content...',
			progress: 25,
			subtasksTotal: 4,
			subtasksCompleted: 1,
			comments: 1
		},
		{ 
			id: 5, 
			col: 'completed', 
			title: '2.2. Design Thinking & Ethics', 
			progress: 100,
			subtasksTotal: 3,
			subtasksCompleted: 3,
			tag: 'in review',
			comments: 1
		},
		{ 
			id: 6, 
			col: 'completed', 
			title: '1. Getting Started 🚀', 
			description: "If you've ever wanted to pursue a career in design, learn the ins-and-outs...",
			tag: 'reviewed',
			comments: 1,
			subtasksTotal: 1,
			subtasksCompleted: 1
		}
	];

	let formattedDate = $derived.by(() => {
		return new Intl.DateTimeFormat('en-GB', {
			day: 'numeric',
			month: 'long',
			year: 'numeric'
		}).format(new Date());
	});
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
		<button class="bg-[#820AD1] hover:bg-[#6a08aa] text-white px-6 py-2.5 rounded-xl font-medium shadow-lg shadow-[#820AD1]/20 transition-all active:scale-95 flex items-center gap-2">
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

					<!-- Cards Area -->
					<div class="flex-1 overflow-y-auto pr-2 space-y-4 custom-scrollbar">
						{#each tasks.filter(t => t.col === column.id) as task (task.id)}
							<div class="bg-white rounded-xl p-5 shadow-[0_2px_15px_-3px_rgba(0,0,0,0.07),0_10px_20px_-2px_rgba(0,0,0,0.04)] hover:shadow-lg transition-all cursor-pointer group border border-surface-50/50">
								<div class="flex justify-between items-start mb-3">
									<h4 class="text-[15px] font-bold text-surface-800 leading-snug">{task.title}</h4>
								</div>

								{#if task.description}
									<p class="text-xs text-surface-400 leading-relaxed mb-4 line-clamp-2">
										{task.description}
									</p>
								{/if}

								<!-- Progress Bar (Image Style) -->
								{#if task.progress !== undefined && (task.subtasksTotal ?? 0) > 0}
									<div class="mb-4">
										<div class="flex items-center justify-between mb-1.5">
											<div class="flex-1 h-1.5 bg-indigo-50 rounded-full overflow-hidden mr-3">
												<div class="h-full bg-indigo-600 rounded-full" style="width: {task.progress}%"></div>
											</div>
											<span class="text-[10px] font-bold text-surface-400">{task.subtasksCompleted ?? 0}/{task.subtasksTotal ?? 0}</span>
										</div>
									</div>
								{/if}

								<div class="flex items-center justify-between mt-auto">
									<div class="flex gap-2">
										{#if task.tag}
											<span class="px-2.5 py-1 bg-indigo-50 text-indigo-700 rounded text-[10px] font-bold">
												{task.tag}
											</span>
										{/if}
									</div>

									<div class="flex items-center gap-3 text-surface-300">
										{#if (task.subtasksTotal ?? 0) > 0}
											<div class="flex items-center gap-1">
												<SquareCheck size={13} strokeWidth={2.5} />
												<span class="text-[10px] font-bold">{task.subtasksTotal ?? 0}</span>
											</div>
										{/if}
										{#if task.comments}
											<div class="flex items-center gap-1">
												<MessageSquare size={13} strokeWidth={2.5} />
												<span class="text-[10px] font-bold">{task.comments}</span>
											</div>
										{/if}
									</div>
								</div>
							</div>
						{/each}

						<!-- Plus Button at bottom of column -->
						<button class="w-full py-4 rounded-xl border-2 border-dashed border-surface-100/50 text-surface-200 hover:text-indigo-500 hover:border-indigo-100 hover:bg-indigo-50/30 transition-all flex items-center justify-center">
							<Plus size={20} />
						</button>
					</div>
				</div>
			{/each}
		</div>
	</div>
</div>

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
