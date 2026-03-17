<script lang="ts">
  import { X, Flag, Clock, Calendar, LoaderCircle } from '@lucide/svelte'
  import {
    Dialog,
    DatePicker,
    Portal,
    type DateValue,
  } from '@skeletonlabs/skeleton-svelte'
  import { createTaskMutation } from '$lib/state/tasks.svelte'
  import { toaster } from '$lib/state/toast.svelte'
  import Input from './Input.svelte'

  interface Props {
    isOpen: boolean
    onClose: () => void
  }

  let { isOpen, onClose }: Props = $props()

  let title = $state('')
  let description = $state('')
  let priority = $state<'low' | 'medium' | 'high' | 'critical'>('medium')
  let dueDateValue = $state<DateValue[]>([])
  let estimatedHours = $state<number | undefined>(undefined)

  const createTask = createTaskMutation()

  function handleOpenChange(e: { open: boolean }) {
    if (!e.open) onClose()
  }

  async function handleSubmit(e: SubmitEvent) {
    e.preventDefault()
    if (!title) return

    const selectedDate = dueDateValue.at(0)

    try {
      await createTask.mutateAsync({
        title,
        description,
        priority,
        dueDate: selectedDate
          ? `${selectedDate.toString()}T00:00:00Z`
          : undefined,
        estimatedHours,
      })

      toaster.create({
        title: 'Task Created',
        description: `"${title}" has been created successfully.`,
      })

      onClose()

      // Reset fields
      title = ''
      description = ''
      priority = 'medium'
      dueDateValue = []
      estimatedHours = undefined
    } catch {
      toaster.create({
        title: 'Error',
        description: 'Failed to create task. Please try again.',
      })
    }
  }
</script>

<Dialog
  open={isOpen}
  onOpenChange={handleOpenChange}
  closeOnInteractOutside={false}
>
  <Portal>
    <Dialog.Backdrop
      class="fixed inset-0 z-50 bg-surface-950/30 backdrop-blur-[2px] transition-opacity"
    />
    <Dialog.Positioner
      class="fixed inset-y-0 right-0 z-50 flex items-stretch justify-end"
    >
      <Dialog.Content
        class="animate-slide-in flex w-full min-w-md flex-col bg-white shadow-2xl"
      >
        <!-- Header -->
        <div
          class="sticky top-0 z-10 border-b border-surface-100 bg-white/95 px-8 pt-8 pb-4 backdrop-blur-sm"
        >
          <div class="flex items-start justify-between gap-4">
            <div class="min-w-0 flex-1">
              <Dialog.Title
                class="text-2xl font-medium tracking-tight text-surface-900"
              >
                Create New Task
              </Dialog.Title>
              <Dialog.Description class="mt-1 text-sm text-surface-800">
                Fill in the details for your new task.
              </Dialog.Description>
            </div>
            <Dialog.CloseTrigger
              class="-mt-1 -mr-2 shrink-0 rounded-xl p-2 text-surface-400 transition-all hover:bg-surface-100 hover:text-surface-900"
            >
              <X size={20} />
            </Dialog.CloseTrigger>
          </div>
        </div>

        <!-- Form -->
        <form
          id="create-task-form"
          onsubmit={handleSubmit}
          class="flex-1 space-y-6 overflow-y-auto px-8 py-6"
        >
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
            <label
              for="description"
              class="block text-sm font-medium text-surface-800"
            >
              Description
            </label>
            <textarea
              id="description"
              bind:value={description}
              placeholder="Add some details about this task..."
              class="block min-h-[120px] w-full resize-none rounded-xl border border-surface-300 bg-surface-50 px-4 py-3 text-surface-900 placeholder-surface-600 transition-all focus:border-primary-500 focus:ring-2 focus:ring-primary-500/10 focus:outline-none sm:text-sm"
            ></textarea>
          </div>

          <!-- Priority -->
          <div class="space-y-1">
            <label
              for="priority"
              class="block text-sm font-medium text-surface-800"
            >
              Priority
            </label>
            <div class="relative">
              <div
                class="pointer-events-none absolute inset-y-0 left-4 flex items-center text-surface-400"
              >
                <Flag size={18} />
              </div>
              <select
                id="priority"
                bind:value={priority}
                class="block h-12 w-full cursor-pointer appearance-none rounded-xl border border-surface-300 bg-surface-50 pr-10 pl-11 text-surface-900 transition-all focus:border-primary-500 focus:ring-2 focus:ring-primary-500/10 focus:outline-none sm:text-sm"
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
            <DatePicker
              value={dueDateValue}
              onValueChange={e => (dueDateValue = e.value)}
            >
              <DatePicker.Label
                class="block text-sm font-medium text-surface-700"
              >
                Due Date
              </DatePicker.Label>
              <DatePicker.Control class="relative">
                <DatePicker.Input
                  placeholder="Pick a date"
                  class="block h-12 w-full rounded-xl border border-surface-300 bg-surface-50 px-4 pr-12 text-surface-900 placeholder-surface-400 transition-all focus:border-primary-500 focus:ring-2 focus:ring-primary-500/10 focus:outline-none sm:text-sm"
                />
                <DatePicker.Trigger
                  class="absolute inset-y-2.5 right-2.5 flex items-center justify-center text-surface-600 transition-colors hover:text-surface-800"
                >
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
                                <DatePicker.TableHeader
                                  >{weekDay.short}</DatePicker.TableHeader
                                >
                              {/each}
                            </DatePicker.TableRow>
                          </DatePicker.TableHead>
                          <DatePicker.TableBody>
                            {#each ctx().weeks as week, id (id)}
                              <DatePicker.TableRow>
                                {#each week as day, id (id)}
                                  <DatePicker.TableCell value={day}>
                                    <DatePicker.TableCellTrigger
                                      >{day.day}</DatePicker.TableCellTrigger
                                    >
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
                            {#each ctx().getMonthsGrid( { columns: 4, format: 'short' } ) as months, id (id)}
                              <DatePicker.TableRow>
                                {#each months as month, id (id)}
                                  <DatePicker.TableCell value={month.value}>
                                    <DatePicker.TableCellTrigger
                                      >{month.label}</DatePicker.TableCellTrigger
                                    >
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
                            {#each ctx().getYearsGrid( { columns: 4 } ) as years, id (id)}
                              <DatePicker.TableRow>
                                {#each years as year, id (id)}
                                  <DatePicker.TableCell value={year.value}>
                                    <DatePicker.TableCellTrigger
                                      >{year.label}</DatePicker.TableCellTrigger
                                    >
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
            <label
              for="estimated-hours"
              class="block text-sm font-medium text-surface-700"
            >
              Estimated Hours
            </label>
            <div class="relative">
              <div
                class="pointer-events-none absolute inset-y-0 left-4 flex items-center text-surface-400"
              >
                <Clock size={18} />
              </div>
              <input
                id="estimated-hours"
                type="number"
                min="0"
                step="0.5"
                bind:value={estimatedHours}
                placeholder="0.0"
                class="block h-12 w-full rounded-xl border border-surface-300 bg-surface-50 pr-4 pl-11 text-surface-900 transition-all focus:border-primary-500 focus:ring-2 focus:ring-primary-500/10 focus:outline-none sm:text-sm"
              />
            </div>
          </div>

          <!-- Spacer to push footer actions away from form fields -->
          <div class="pt-2"></div>
        </form>

        <!-- Footer Actions (sticky at bottom) -->
        <div
          class="flex items-center justify-end gap-3 border-t border-surface-100 bg-white px-8 py-5"
        >
          <Dialog.CloseTrigger
            type="button"
            class="rounded-xl px-6 py-2.5 font-medium text-surface-500 transition-all hover:bg-surface-50 hover:text-surface-900"
          >
            Cancel
          </Dialog.CloseTrigger>
          <button
            type="submit"
            form="create-task-form"
            disabled={createTask.isPending}
            onclick={e => {
              // Programmatically submit the form since the button is outside the <form>
              e.preventDefault()
              const form = document.querySelector('form') as HTMLFormElement
              form?.requestSubmit()
            }}
            class="flex items-center gap-2 rounded-xl bg-primary-500 px-8 py-2.5 font-medium text-white shadow-lg shadow-primary-500/20 transition-all hover:bg-primary-700 active:scale-95 disabled:cursor-not-allowed disabled:opacity-50"
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
