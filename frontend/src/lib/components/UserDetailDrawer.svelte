<script lang="ts">
  import { Dialog, Portal, Progress } from '@skeletonlabs/skeleton-svelte'
  import {
    X,
    Mail,
    Shield,
    User as UserIcon,
    Calendar,
    CircleCheck,
    Clock,
  } from '@lucide/svelte'
  import { createUserQuery } from '$lib/state/user.svelte'
  import { resolveAvatarUrl } from '$lib/utils/avatar'

  interface Props {
    userId: string | undefined
    open: boolean
    onOpenChange: (open: boolean) => void
  }

  let { userId, open, onOpenChange }: Props = $props()

  const userQuery = createUserQuery(() => userId)
  const user = $derived(userQuery.data)
  const isLoading = $derived(userQuery.isLoading)

  function formatDate(dateStr: string | undefined) {
    if (!dateStr) return 'N/A'
    return new Intl.DateTimeFormat('en-GB', {
      day: '2-digit',
      month: 'long',
      year: 'numeric',
    }).format(new Date(dateStr))
  }
</script>

<Dialog
  {open}
  onOpenChange={e => onOpenChange(e.open)}
  closeOnInteractOutside={true}
>
  <Portal>
    <Dialog.Backdrop
      class="fixed inset-0 z-50 bg-surface-950/30 backdrop-blur-[2px] transition-opacity"
    />
    <Dialog.Positioner
      class="fixed inset-y-0 right-0 z-50 flex items-stretch justify-end"
    >
      <Dialog.Content
        class="animate-slide-in flex w-full max-w-[400px] flex-col overflow-y-auto bg-white shadow-2xl"
      >
        <div
          class="flex items-center justify-between border-b border-surface-100 bg-surface-50 p-6"
        >
          <h2
            class="flex items-center gap-2 text-xl font-bold text-surface-900"
          >
            <UserIcon size={20} class="text-primary-500" />
            User Profile
          </h2>
          <button
            onclick={() => onOpenChange(false)}
            class="rounded-full p-2 text-surface-500 transition-colors hover:bg-surface-200"
          >
            <X size={20} />
          </button>
        </div>

        <div class="flex-1 p-8">
          {#if isLoading}
            <div
              class="flex h-full flex-col items-center justify-center gap-4 py-20"
            >
              <Progress value={null} class="w-fit">
                <Progress.Circle class="size-12">
                  <Progress.CircleTrack />
                  <Progress.CircleRange class="stroke-primary-500" />
                </Progress.Circle>
              </Progress>
              <p class="animate-pulse font-medium text-surface-500">
                Loading user data...
              </p>
            </div>
          {:else if user}
            <!-- Header Section -->
            <div class="mb-8 flex flex-col items-center text-center">
              <div class="relative mb-4">
                {#if user.avatarUrl}
                  <img
                    src={resolveAvatarUrl(user.avatarUrl) ?? ''}
                    alt={user.firstName}
                    class="size-32 rounded-3xl border-4 border-white object-cover shadow-xl"
                  />
                {:else}
                  <div
                    class="flex size-32 items-center justify-center rounded-3xl border-4 border-white bg-primary-100 shadow-xl"
                  >
                    <span class="text-4xl font-bold text-primary-600">
                      {user.firstName[0]}{user.lastName[0]}
                    </span>
                  </div>
                {/if}
                <div
                  class="absolute -right-2 -bottom-2 rounded-full bg-white p-1.5 shadow-md"
                >
                  {#if user.isActive}
                    <CircleCheck size={24} class="text-emerald-500" />
                  {:else}
                    <Clock size={24} class="text-surface-400" />
                  {/if}
                </div>
              </div>
              <h3 class="text-2xl font-bold text-surface-900">
                {user.firstName}
                {user.lastName}
              </h3>
              <span
                class="mt-2 inline-flex items-center rounded-full bg-primary-50 px-3 py-1 text-xs font-semibold text-primary-700 capitalize"
              >
                {user.role}
              </span>
            </div>

            <!-- Info Cards -->
            <div class="space-y-4">
              <div
                class="rounded-2xl border border-surface-100 bg-surface-50 p-4"
              >
                <div class="mb-1 flex items-center gap-3 text-surface-500">
                  <Mail size={16} />
                  <span class="text-xs font-medium tracking-wider uppercase"
                    >Email Address</span
                  >
                </div>
                <p class="font-medium break-all text-surface-900">
                  {user.email}
                </p>
              </div>

              <div
                class="rounded-2xl border border-surface-100 bg-surface-50 p-4"
              >
                <div class="mb-1 flex items-center gap-3 text-surface-500">
                  <Shield size={16} />
                  <span class="text-xs font-medium tracking-wider uppercase"
                    >Role & Permissions</span
                  >
                </div>
                <p class="font-medium text-surface-900 capitalize">
                  {user.role}
                </p>
              </div>

              <div
                class="rounded-2xl border border-surface-100 bg-surface-50 p-4"
              >
                <div class="mb-1 flex items-center gap-3 text-surface-500">
                  <Calendar size={16} />
                  <span class="text-xs font-medium tracking-wider uppercase"
                    >Member Since</span
                  >
                </div>
                <p class="font-medium text-surface-900">
                  {formatDate(user.createdAt)}
                </p>
              </div>
            </div>

            <!-- Stats section -->
            <div class="mt-8 grid grid-cols-2 gap-4">
              <div class="rounded-2xl bg-indigo-50 p-4 text-center">
                <span class="block text-2xl font-bold text-indigo-700">--</span>
                <span class="text-xs font-medium text-indigo-600 uppercase"
                  >Tasks Active</span
                >
              </div>
              <div class="rounded-2xl bg-emerald-50 p-4 text-center">
                <span class="block text-2xl font-bold text-emerald-700">--</span
                >
                <span class="text-xs font-medium text-emerald-600 uppercase"
                  >Completed</span
                >
              </div>
            </div>
          {/if}
        </div>

        <div class="mt-auto border-t border-surface-100 bg-surface-50 p-6">
          <button
            onclick={() => onOpenChange(false)}
            class="w-full rounded-xl border border-surface-200 bg-white py-3 font-semibold text-surface-700 shadow-sm transition-colors hover:bg-surface-100"
          >
            Close Profile
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
