<script lang="ts">
  import { ExternalLink, Ban } from '@lucide/svelte'
  import type { TaskResponse } from '$lib/api/types'
  import { priorityConfig, statusConfig } from '$lib/utils/task'
  import { resolveAvatarUrl } from '$lib/utils/avatar'
  import Button from '$lib/components/Button.svelte'

  interface Props {
    tasks: TaskResponse[]
    onViewDetails: (task: TaskResponse) => void
  }

  let { tasks, onViewDetails }: Props = $props()
</script>

<div class="overflow-hidden rounded-2xl border border-surface-200 bg-white shadow-sm">
  <table class="w-full text-left border-collapse">
    <thead>
      <tr class="border-b border-surface-100 bg-surface-50/50">
        <th class="px-6 py-4 text-[11px] font-bold uppercase tracking-wider text-surface-500">Task Title</th>
        <th class="px-6 py-4 text-[11px] font-bold uppercase tracking-wider text-surface-500">Status</th>
        <th class="px-6 py-4 text-[11px] font-bold uppercase tracking-wider text-surface-500">Priority</th>
        <th class="px-6 py-4 text-[11px] font-bold uppercase tracking-wider text-surface-500">Assignee</th>
        <th class="px-6 py-4 text-[11px] font-bold uppercase tracking-wider text-surface-500 text-right">Actions</th>
      </tr>
    </thead>
    <tbody class="divide-y divide-surface-100">
      {#each tasks as task (task.id)}
        <tr class="group transition-all hover:bg-surface-50/80">
          <td class="px-6 py-4">
            <span class="text-sm font-medium text-surface-900">{task.title}</span>
          </td>
          <td class="px-6 py-4">
            <span class="inline-flex items-center rounded-full px-2.5 py-0.5 text-[11px] font-semibold {statusConfig[task.status].class}">
              {statusConfig[task.status].label}
            </span>
          </td>
          <td class="px-6 py-4">
            <span class="inline-flex items-center rounded-full px-2.5 py-0.5 text-[11px] font-semibold {priorityConfig[task.priority].class}">
              {priorityConfig[task.priority].label}
            </span>
          </td>
          <td class="px-6 py-4">
             <div class="flex items-center gap-2">
                {#if task.assignee.avatarUrl}
                  <img 
                    src={resolveAvatarUrl(task.assignee.avatarUrl) ?? ''} 
                    alt="" 
                    class="size-6 rounded-full object-cover border border-surface-200" 
                  />
                {:else}
                  <div class="flex size-6 items-center justify-center rounded-full bg-indigo-100 text-[10px] font-bold text-indigo-700">
                    {task.assignee.firstName[0]}{task.assignee.lastName[0]}
                  </div>
                {/if}
                <span class="text-sm text-surface-600 font-medium">{task.assignee.firstName} {task.assignee.lastName}</span>
             </div>
          </td>
          <td class="px-6 py-4 text-right">
            <div class="flex items-center justify-end gap-1 opacity-0 transition-opacity group-hover:opacity-100">
               <Button 
                variant="ghost" 
                onclick={() => onViewDetails(task)}
                class="size-8 !p-0 text-surface-400 hover:text-surface-900"
                title="View details"
               >
                 <ExternalLink size={14} />
               </Button>
            </div>
          </td>
        </tr>
      {:else}
        <tr>
          <td colspan="5" class="px-6 py-20 text-center">
            <div class="flex flex-col items-center gap-3">
               <div class="flex size-12 items-center justify-center rounded-2xl bg-surface-50 text-surface-200">
                 <Ban size={24} />
               </div>
               <div class="space-y-1">
                 <p class="text-sm font-medium text-surface-900">No archived tasks</p>
                 <p class="text-xs text-surface-500">Tasks you archive will appear here.</p>
               </div>
            </div>
          </td>
        </tr>
      {/each}
    </tbody>
  </table>
</div>
