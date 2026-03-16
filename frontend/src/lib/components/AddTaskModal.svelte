<script lang="ts">
	import { X, Flag, Clock, Calendar, LoaderCircle } from '@lucide/svelte';
	import { Dialog, DatePicker, Portal, type DateValue } from '@skeletonlabs/skeleton-svelte';
	import { createTaskMutation } from '$lib/state/tasks.svelte';
	import { toaster } from '$lib/state/toast.svelte';
	import Input from './Input.svelte';

	interface Props {
		isOpen: boolean;
		onClose: () => void;
	}

	let { isOpen, onClose }: Props = $props();

	let title = $state('');
	let description = $state('');
	let priority = $state<'low' | 'medium' | 'high' | 'critical'>('medium');
	let dueDateValue = $state<DateValue[]>([]);
	let estimatedHours = $state<number | undefined>(undefined);

	const createTask = createTaskMutation();

	function handleOpenChange(e: { open: boolean }) {
		if (!e.open) onClose();
	}

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault();
		if (!title) return;

		const selectedDate = dueDateValue.at(0);

		try {
			await createTask.mutateAsync({
				title,
				description,
				priority,
				dueDate: selectedDate ? `${selectedDate.toString()}T00:00:00Z` : undefined,
				estimatedHours
			});

			toaster.create({
				title: 'Task Created',
				description: `"${title}" has been created successfully.`
			});

			onClose();

			// Reset fields
			title = '';
			description = '';
			priority = 'medium';
			dueDateValue = [];
			estimatedHours = undefined;
		} catch {
			toaster.create({
				title: 'Error',
				description: 'Failed to create task. Please try again.',
			});
		}
	}
</script>

<Dialog open={isOpen} onOpenChange={handleOpenChange} closeOnInteractOutside={false}>
	<Portal>
		<Dialog.Backdrop class="fixed inset-0 z-50 bg-surface-950/30 backdrop-blur-[2px] transition-opacity" />
		<Dialog.Positioner class="fixed inset-y-0 right-0 z-50 flex items-stretch justify-end">
			<Dialog.Content class="w-full min-w-md bg-white shadow-2xl flex flex-col animate-slide-in">
				<!-- Header -->
				<div class="sticky top-0 bg-white/95 backdrop-blur-sm z-10 px-8 pt-8 pb-4 border-b border-surface-100">
					<div class="flex items-start justify-between gap-4">
						<div class="flex-1 min-w-0">
							<Dialog.Title class="text-2xl font-medium text-surface-900 tracking-tight">
								Create New Task
							</Dialog.Title>
							<Dialog.Description class="text-sm text-surface-500 mt-1">
								Fill in the details for your new task.
							</Dialog.Description>
						</div>
						<Dialog.CloseTrigger
							class="p-2 -mr-2 -mt-1 rounded-xl text-surface-400 hover:text-surface-900 hover:bg-surface-100 transition-all shrink-0"
						>
							<X size={20} />
						</Dialog.CloseTrigger>
					</div>
				</div>

				<!-- Form -->
				<form id="create-task-form" onsubmit={handleSubmit} class="flex-1 overflow-y-auto px-8 py-6 space-y-6">
					<!-- Title -->
					<Input
						id="title"
						label="Task Title"
						placeholder="What needs to be done?"
						bind:value={title}
						required
					/>

					<!-- Description -->
					<div class="space-y-1">
						<label for="description" class="block text-sm font-medium text-surface-700">
							Description
						</label>
						<textarea
							id="description"
							bind:value={description}
							placeholder="Add some details about this task..."
							class="block w-full rounded-xl border border-surface-300 bg-surface-50 px-4 py-3 text-surface-900 placeholder-surface-600 transition-all focus:border-primary-500 focus:ring-2 focus:ring-primary-500/10 focus:outline-none sm:text-sm min-h-[120px] resize-none"
						></textarea>
					</div>

					<!-- Priority -->
					<div class="space-y-1">
						<label for="priority" class="block text-sm font-medium text-surface-700">
							Priority
						</label>
						<div class="relative">
							<div class="absolute inset-y-0 left-4 flex items-center text-surface-400 pointer-events-none">
								<Flag size={18} />
							</div>
							<select
								id="priority"
								bind:value={priority}
								class="block h-12 w-full appearance-none rounded-xl border border-surface-300 bg-surface-50 pl-11 pr-10 text-surface-900 transition-all focus:border-primary-500 focus:ring-2 focus:ring-primary-500/10 focus:outline-none sm:text-sm cursor-pointer"
							>
								<option value="low">Low</option>
								<option value="medium">Medium</option>
								<option value="high">High</option>
								<option value="critical">Critical</option>
							</select>
						</div>
					</div>

					<!-- Due Date -->
					<div class="space-y-1">
						<DatePicker value={dueDateValue} onValueChange={(e) => (dueDateValue = e.value)}>
							<DatePicker.Label class="block text-sm font-medium text-surface-700">
								Due Date
							</DatePicker.Label>
							<DatePicker.Control class="relative">
								<DatePicker.Input
									placeholder="Pick a date"
									class="block h-12 w-full rounded-xl border border-surface-300 bg-surface-50 px-4 pr-12 text-surface-900 placeholder-surface-400 transition-all focus:border-primary-500 focus:ring-2 focus:ring-primary-500/10 focus:outline-none sm:text-sm"
								/>
								<DatePicker.Trigger class="absolute inset-y-2.5 right-2.5 flex items-center justify-center text-surface-600 hover:text-surface-800 transition-colors">
									<Calendar size={18} />
								</DatePicker.Trigger>
							</DatePicker.Control>
							<DatePicker.Positioner>
								<DatePicker.Content>
										<DatePicker.View view="day">
											<DatePicker.Context>
												{#snippet children(ctx)}
													<DatePicker.ViewControl>
														<DatePicker.PrevTrigger />
														<DatePicker.ViewTrigger>
															<DatePicker.RangeText />
														</DatePicker.ViewTrigger>
														<DatePicker.NextTrigger />
													</DatePicker.ViewControl>
													<DatePicker.Table>
														<DatePicker.TableHead>
															<DatePicker.TableRow>
																{#each ctx().weekDays as weekDay, id (id)}
																	<DatePicker.TableHeader>{weekDay.short}</DatePicker.TableHeader>
																{/each}
															</DatePicker.TableRow>
														</DatePicker.TableHead>
														<DatePicker.TableBody>
															{#each ctx().weeks as week, id (id)}
																<DatePicker.TableRow>
																	{#each week as day, id (id)}
																		<DatePicker.TableCell value={day}>
																			<DatePicker.TableCellTrigger>{day.day}</DatePicker.TableCellTrigger>
																		</DatePicker.TableCell>
																	{/each}
																</DatePicker.TableRow>
															{/each}
														</DatePicker.TableBody>
													</DatePicker.Table>
												{/snippet}
											</DatePicker.Context>
										</DatePicker.View>
										<DatePicker.View view="month">
											<DatePicker.Context>
												{#snippet children(ctx)}
													<DatePicker.ViewControl>
														<DatePicker.PrevTrigger />
														<DatePicker.ViewTrigger>
															<DatePicker.RangeText />
														</DatePicker.ViewTrigger>
														<DatePicker.NextTrigger />
													</DatePicker.ViewControl>
													<DatePicker.Table>
														<DatePicker.TableBody>
															{#each ctx().getMonthsGrid({ columns: 4, format: 'short' }) as months, id (id)}
																<DatePicker.TableRow>
																	{#each months as month, id (id)}
																		<DatePicker.TableCell value={month.value}>
																			<DatePicker.TableCellTrigger>{month.label}</DatePicker.TableCellTrigger>
																		</DatePicker.TableCell>
																	{/each}
																</DatePicker.TableRow>
															{/each}
														</DatePicker.TableBody>
													</DatePicker.Table>
												{/snippet}
											</DatePicker.Context>
										</DatePicker.View>
										<DatePicker.View view="year">
											<DatePicker.Context>
												{#snippet children(ctx)}
													<DatePicker.ViewControl>
														<DatePicker.PrevTrigger />
														<DatePicker.ViewTrigger>
															<DatePicker.RangeText />
														</DatePicker.ViewTrigger>
														<DatePicker.NextTrigger />
													</DatePicker.ViewControl>
													<DatePicker.Table>
														<DatePicker.TableBody>
															{#each ctx().getYearsGrid({ columns: 4 }) as years, id (id)}
																<DatePicker.TableRow>
																	{#each years as year, id (id)}
																		<DatePicker.TableCell value={year.value}>
																			<DatePicker.TableCellTrigger>{year.label}</DatePicker.TableCellTrigger>
																		</DatePicker.TableCell>
																	{/each}
																</DatePicker.TableRow>
															{/each}
														</DatePicker.TableBody>
													</DatePicker.Table>
												{/snippet}
											</DatePicker.Context>
										</DatePicker.View>
								</DatePicker.Content>
							</DatePicker.Positioner>
						</DatePicker>
					</div>

					<!-- Estimated Hours -->
					<div class="space-y-1">
						<label for="estimated-hours" class="block text-sm font-medium text-surface-700">
							Estimated Hours
						</label>
						<div class="relative">
							<div class="absolute inset-y-0 left-4 flex items-center text-surface-400 pointer-events-none">
								<Clock size={18} />
							</div>
							<input
								id="estimated-hours"
								type="number"
								min="0"
								step="0.5"
								bind:value={estimatedHours}
								placeholder="0.0"
								class="block h-12 w-full rounded-xl border border-surface-300 bg-surface-50 pl-11 pr-4 text-surface-900 transition-all focus:border-primary-500 focus:ring-2 focus:ring-primary-500/10 focus:outline-none sm:text-sm"
								/>
							</div>
						</div>

					<!-- Spacer to push footer actions away from form fields -->
					<div class="pt-2"></div>
				</form>

				<!-- Footer Actions (sticky at bottom) -->
				<div class="px-8 py-5 border-t border-surface-100 bg-white flex items-center justify-end gap-3">
					<Dialog.CloseTrigger
						type="button"
						class="px-6 py-2.5 rounded-xl font-medium text-surface-500 hover:text-surface-900 hover:bg-surface-50 transition-all"
					>
						Cancel
					</Dialog.CloseTrigger>
					<button
						type="submit"
						form="create-task-form"
						disabled={createTask.isPending}
						onclick={(e) => {
							// Programmatically submit the form since the button is outside the <form>
							e.preventDefault();
							const form = document.querySelector('form') as HTMLFormElement;
							form?.requestSubmit();
						}}
						class="bg-primary-500 hover:bg-primary-700 text-white px-8 py-2.5 rounded-xl font-medium shadow-lg shadow-primary-500/20 transition-all active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
					>
						{#if createTask.isPending}
							<LoaderCircle size={18} class="animate-spin" />
							Creating...
						{:else}
							Create Task
						{/if}
					</button>
				</div>
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
