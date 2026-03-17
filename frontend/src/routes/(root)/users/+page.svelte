<script lang="ts">
  import {
    Search,
    LoaderCircle,
    ChevronLeft,
    ChevronRight,
    X,
    ShieldCheck,
    User,
    FilterX,
    Pencil,
    Trash2,
    TriangleAlert,
  } from '@lucide/svelte'
  import { goto } from '$app/navigation'
  import { resolveAvatarUrl } from '$lib/utils/avatar'
  import { resolve } from '$app/paths'
  import { Dialog, Portal } from '@skeletonlabs/skeleton-svelte'
  import { getUsersQuery, deleteUserMutation } from '$lib/state/user.svelte'
  import Button from '$lib/components/Button.svelte'
  import { toaster } from '$lib/state/toast.svelte'
  import type { UserRole, UserResponse } from '$lib/api/types'

  const PAGE_SIZE = 10

  let searchInput = $state('')
  let debouncedSearch = $state('')
  let filterRole = $state<UserRole | ''>('')
  let sortField = $state('createdAt')
  let sortOrder = $state<'asc' | 'desc'>('desc')
  let currentPage = $state(1)

  let debounceTimer: ReturnType<typeof setTimeout>
  function onSearchInput(value: string) {
    searchInput = value
    clearTimeout(debounceTimer)
    debounceTimer = setTimeout(() => {
      debouncedSearch = value
      currentPage = 1
    }, 350)
  }

  const usersQuery = getUsersQuery(() => ({
    page: currentPage,
    limit: PAGE_SIZE,
    firstName: debouncedSearch || undefined,
    role: (filterRole as UserRole) || undefined,
    sort: sortField,
    order: sortOrder,
  }))

  const hasActiveFilters = $derived(debouncedSearch !== '' || filterRole !== '')

  function clearFilters() {
    searchInput = ''
    debouncedSearch = ''
    filterRole = ''
    currentPage = 1
  }

  function setSort(field: string) {
    if (sortField === field) {
      sortOrder = sortOrder === 'asc' ? 'desc' : 'asc'
    } else {
      sortField = field
      sortOrder = 'asc'
    }
    currentPage = 1
  }

  function goToPage(page: number) {
    currentPage = page
  }

  function formatDate(dateStr: string) {
    return new Intl.DateTimeFormat('en-GB', {
      day: '2-digit',
      month: 'short',
      year: 'numeric',
    }).format(new Date(dateStr))
  }

  function getInitials(firstName: string, lastName: string) {
    return `${firstName[0] ?? ''}${lastName[0] ?? ''}`.toUpperCase()
  }

  const totalPages = $derived(usersQuery.data?.totalPages ?? 1)
  const totalUsers = $derived(usersQuery.data?.total ?? 0)

  function getPageNumbers(total: number, current: number): (number | '...')[] {
    if (total <= 7) return Array.from({ length: total }, (_, i) => i + 1)
    const pages: (number | '...')[] = [1]
    if (current > 3) pages.push('...')
    for (
      let i = Math.max(2, current - 1);
      i <= Math.min(total - 1, current + 1);
      i++
    ) {
      pages.push(i)
    }
    if (current < total - 2) pages.push('...')
    pages.push(total)
    return pages
  }

  // Delete confirmation
  const deleteUser = deleteUserMutation()
  let userToDelete = $state<UserResponse | null>(null)
  let isConfirmOpen = $state(false)

  function openDeleteConfirm(user: UserResponse) {
    userToDelete = user
    isConfirmOpen = true
  }

  async function confirmDelete() {
    if (!userToDelete) return
    const name = `${userToDelete.firstName} ${userToDelete.lastName}`
    try {
      await deleteUser.mutateAsync(userToDelete.id)
      isConfirmOpen = false
      userToDelete = null
      toaster.success({
        title: 'User Deleted',
        description: `${name} has been removed from the system.`,
      })
    } catch {
      toaster.error({
        title: 'Delete Failed',
        description: `Could not delete ${name}. Please try again.`,
      })
    }
  }
</script>

<svelte:head>
  <title>Users - Taskify</title>
</svelte:head>

<div class="flex h-full flex-col pt-8">
  <!-- Page Header -->
  <header class="px-8 pb-8">
    <div class="space-y-1">
      <h2 class="text-4xl leading-tight tracking-tight text-surface-900">
        <span class="font-light">Team</span>
        <span class="font-normal"> Members</span>
      </h2>
      <p class="text-xl font-light text-surface-800">
        {#if usersQuery.isSuccess}
          {totalUsers} {totalUsers === 1 ? 'user' : 'users'} in the system
        {:else}
          Manage and view all users
        {/if}
      </p>
    </div>
  </header>

  <!-- Filters Bar -->
  <div class="mb-6 flex flex-wrap items-center gap-3 px-8">
    <!-- Search -->
    <div class="relative max-w-xs min-w-[220px] flex-1">
      <div
        class="pointer-events-none absolute inset-y-0 left-3 flex items-center text-surface-400"
      >
        <Search size={16} />
      </div>
      <input
        type="text"
        value={searchInput}
        oninput={e => onSearchInput((e.target as HTMLInputElement).value)}
        placeholder="Search by first name..."
        class="h-10 w-full rounded-xl border border-surface-300 bg-white pr-4 pl-9 text-sm text-surface-900 placeholder-surface-400 transition-all focus:border-primary-500 focus:ring-2 focus:ring-primary-500/10 focus:outline-none"
      />
    </div>

    <!-- Role Filter -->
    <select
      bind:value={filterRole}
      onchange={() => (currentPage = 1)}
      class="h-10 cursor-pointer appearance-none rounded-xl border border-surface-300 bg-white px-3 pr-8 text-sm text-surface-700 transition-all focus:border-primary-500 focus:ring-2 focus:ring-primary-500/10 focus:outline-none"
    >
      <option value="">All roles</option>
      <option value="admin">Admin</option>
      <option value="employee">Employee</option>
    </select>

    <!-- Sort -->
    <select
      bind:value={sortField}
      onchange={() => (currentPage = 1)}
      class="h-10 cursor-pointer appearance-none rounded-xl border border-surface-300 bg-white px-3 pr-8 text-sm text-surface-700 transition-all focus:border-primary-500 focus:ring-2 focus:ring-primary-500/10 focus:outline-none"
    >
      <option value="createdAt">Sort: Joined date</option>
      <option value="firstName">Sort: First name</option>
      <option value="lastName">Sort: Last name</option>
      <option value="email">Sort: Email</option>
    </select>

    <button
      onclick={() => setSort(sortField)}
      class="flex h-10 items-center gap-1.5 rounded-xl border border-surface-300 bg-white px-3 text-sm text-surface-700 transition-all hover:bg-surface-50"
      title="Toggle sort order"
    >
      {sortOrder === 'asc' ? '↑ Asc' : '↓ Desc'}
    </button>

    <!-- Clear Filters -->
    {#if hasActiveFilters}
      <button
        onclick={clearFilters}
        class="flex h-10 items-center gap-1.5 rounded-xl border border-surface-300 bg-white px-3 text-sm text-surface-500 transition-all hover:bg-surface-50 hover:text-surface-900"
      >
        <FilterX size={15} />
        Clear
      </button>
    {/if}
  </div>

  <!-- Table -->
  <div class="flex-1 overflow-y-auto px-8 pb-8">
    {#if usersQuery.isPending}
      <div
        class="flex flex-col items-center justify-center py-24 text-surface-400"
      >
        <LoaderCircle size={28} class="mb-3 animate-spin" />
        <span class="text-sm font-medium">Loading users...</span>
      </div>
    {:else if usersQuery.isError}
      <div
        class="flex flex-col items-center justify-center py-24 text-rose-500"
      >
        <X size={28} class="mb-3" />
        <span class="text-sm font-medium">Failed to load users</span>
      </div>
    {:else if usersQuery.data}
      {@const users = usersQuery.data.data}

      {#if users.length === 0}
        <div
          class="flex flex-col items-center justify-center py-24 text-surface-400"
        >
          <User size={32} class="mb-3 opacity-40" />
          <p class="text-sm font-medium">No users found</p>
          {#if hasActiveFilters}
            <button
              onclick={clearFilters}
              class="mt-2 text-xs text-primary-500 hover:underline"
            >
              Clear filters
            </button>
          {/if}
        </div>
      {:else}
        <!-- Users Table -->
        <div
          class="overflow-hidden rounded-2xl border border-surface-200 bg-white"
        >
          <table class="w-full">
            <thead>
              <tr class="border-b border-surface-100">
                <th
                  class="cursor-pointer px-6 py-3.5 text-left text-xs font-semibold tracking-wide text-surface-500 uppercase transition-colors hover:text-surface-700"
                  onclick={() => setSort('firstName')}
                >
                  <span class="flex items-center gap-1">
                    User
                    {#if sortField === 'firstName'}
                      <span class="text-primary-500"
                        >{sortOrder === 'asc' ? '↑' : '↓'}</span
                      >
                    {/if}
                  </span>
                </th>
                <th
                  class="cursor-pointer px-6 py-3.5 text-left text-xs font-semibold tracking-wide text-surface-500 uppercase transition-colors hover:text-surface-700"
                  onclick={() => setSort('email')}
                >
                  <span class="flex items-center gap-1">
                    Email
                    {#if sortField === 'email'}
                      <span class="text-primary-500"
                        >{sortOrder === 'asc' ? '↑' : '↓'}</span
                      >
                    {/if}
                  </span>
                </th>
                <th
                  class="px-6 py-3.5 text-left text-xs font-semibold tracking-wide text-surface-500 uppercase"
                >
                  Role
                </th>
                <th
                  class="px-6 py-3.5 text-left text-xs font-semibold tracking-wide text-surface-500 uppercase"
                >
                  Status
                </th>
                <th
                  class="cursor-pointer px-6 py-3.5 text-left text-xs font-semibold tracking-wide text-surface-500 uppercase transition-colors hover:text-surface-700"
                  onclick={() => setSort('createdAt')}
                >
                  <span class="flex items-center gap-1">
                    Joined
                    {#if sortField === 'createdAt'}
                      <span class="text-primary-500"
                        >{sortOrder === 'asc' ? '↑' : '↓'}</span
                      >
                    {/if}
                  </span>
                </th>
                <th
                  class="px-6 py-3.5 text-right text-xs font-semibold tracking-wide text-surface-500 uppercase"
                >
                  Actions
                </th>
              </tr>
            </thead>
            <tbody class="divide-y divide-surface-50">
              {#each users as user (user.id)}
                <tr class="group transition-colors hover:bg-surface-50/60">
                  <!-- Avatar + Name -->
                  <td class="px-6 py-4">
                    <div class="flex items-center gap-3">
                      {#if user.avatarUrl}
                        <img
                          src={resolveAvatarUrl(user.avatarUrl) ?? ''}
                          alt="{user.firstName} {user.lastName}"
                          class="size-9 shrink-0 rounded-xl border border-surface-200 object-cover"
                        />
                      {:else}
                        <div
                          class="flex size-9 shrink-0 items-center justify-center rounded-xl border border-primary-100 bg-primary-50"
                        >
                          <span class="text-[11px] font-bold text-primary-600">
                            {getInitials(user.firstName, user.lastName)}
                          </span>
                        </div>
                      {/if}
                      <div class="min-w-0">
                        <p
                          class="truncate text-sm font-medium text-surface-900"
                        >
                          {user.firstName}
                          {user.lastName}
                        </p>
                      </div>
                    </div>
                  </td>

                  <!-- Email -->
                  <td class="px-6 py-4">
                    <span class="truncate text-sm text-surface-600"
                      >{user.email}</span
                    >
                  </td>

                  <!-- Role -->
                  <td class="px-6 py-4">
                    {#if user.role === 'admin'}
                      <span
                        class="inline-flex items-center gap-1 rounded-full bg-violet-50 px-2.5 py-1 text-[11px] font-semibold text-violet-700"
                      >
                        <ShieldCheck size={11} />
                        Admin
                      </span>
                    {:else}
                      <span
                        class="inline-flex items-center gap-1 rounded-full bg-sky-50 px-2.5 py-1 text-[11px] font-semibold text-sky-700"
                      >
                        <User size={11} />
                        Employee
                      </span>
                    {/if}
                  </td>

                  <!-- Status -->
                  <td class="px-6 py-4">
                    {#if user.isActive}
                      <span
                        class="inline-flex items-center gap-1.5 text-[11px] font-semibold text-emerald-600"
                      >
                        <span class="size-1.5 rounded-full bg-emerald-500"
                        ></span>
                        Active
                      </span>
                    {:else}
                      <span
                        class="inline-flex items-center gap-1.5 text-[11px] font-semibold text-surface-400"
                      >
                        <span class="size-1.5 rounded-full bg-surface-300"
                        ></span>
                        Inactive
                      </span>
                    {/if}
                  </td>

                  <!-- Joined -->
                  <td class="px-6 py-4">
                    <span class="text-sm text-surface-500"
                      >{formatDate(user.createdAt)}</span
                    >
                  </td>

                  <!-- Actions -->
                  <td class="px-6 py-4">
                    <div class="flex items-center justify-end gap-1">
                      <button
                        onclick={() => goto(resolve(`/users/${user.id}/edit`))}
                        class="rounded-lg p-2 text-surface-400 transition-all hover:bg-primary-50 hover:text-primary-500"
                        title="Edit user"
                      >
                        <Pencil size={15} />
                      </button>
                      <button
                        onclick={() => openDeleteConfirm(user)}
                        class="rounded-lg p-2 text-surface-400 transition-all hover:bg-rose-50 hover:text-rose-500"
                        title="Delete user"
                      >
                        <Trash2 size={15} />
                      </button>
                    </div>
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>

        <!-- Pagination -->
        {#if totalPages > 1}
          <div class="mt-6 flex items-center justify-between">
            <p class="text-sm text-surface-500">
              Showing
              <span class="font-medium text-surface-700">
                {(currentPage - 1) * PAGE_SIZE + 1}–{Math.min(
                  currentPage * PAGE_SIZE,
                  totalUsers
                )}
              </span>
              of <span class="font-medium text-surface-700">{totalUsers}</span> users
            </p>

            <div class="flex items-center gap-1">
              <button
                onclick={() => goToPage(currentPage - 1)}
                disabled={currentPage === 1}
                class="flex size-9 items-center justify-center rounded-xl border border-surface-200 text-surface-500 transition-all hover:bg-surface-50 disabled:cursor-not-allowed disabled:opacity-40"
              >
                <ChevronLeft size={16} />
              </button>

              {#each getPageNumbers(totalPages, currentPage) as page, i (i)}
                {#if page === '...'}
                  <span
                    class="flex size-9 items-center justify-center text-sm text-surface-400"
                    >…</span
                  >
                {:else}
                  <button
                    onclick={() => goToPage(page as number)}
                    class="flex size-9 items-center justify-center rounded-xl border text-sm font-medium transition-all
										{currentPage === page
                      ? 'border-primary-500 bg-primary-500 text-white shadow-sm'
                      : 'border-surface-200 text-surface-600 hover:bg-surface-50'}"
                  >
                    {page}
                  </button>
                {/if}
              {/each}

              <button
                onclick={() => goToPage(currentPage + 1)}
                disabled={currentPage === totalPages}
                class="flex size-9 items-center justify-center rounded-xl border border-surface-200 text-surface-500 transition-all hover:bg-surface-50 disabled:cursor-not-allowed disabled:opacity-40"
              >
                <ChevronRight size={16} />
              </button>
            </div>
          </div>
        {/if}
      {/if}
    {/if}
  </div>
</div>

<!-- Delete Confirmation Dialog -->
<Dialog
  role="alertdialog"
  open={isConfirmOpen}
  onOpenChange={e => {
    if (!e.open) {
      isConfirmOpen = false
      userToDelete = null
    }
  }}
>
  <Portal>
    <Dialog.Backdrop
      class="fixed inset-0 z-50 bg-surface-950/40 backdrop-blur-sm"
    />
    <Dialog.Positioner
      class="fixed inset-0 z-50 flex items-center justify-center p-4"
    >
      <Dialog.Content
        class="w-full max-w-md space-y-4 rounded-2xl border border-surface-100 bg-white p-6 shadow-2xl"
      >
        <div class="flex items-start gap-4">
          <div
            class="flex size-10 shrink-0 items-center justify-center rounded-xl bg-rose-50 text-rose-500"
          >
            <TriangleAlert size={20} />
          </div>
          <div>
            <Dialog.Title class="text-base font-bold text-surface-900">
              Delete User
            </Dialog.Title>
            <Dialog.Description class="mt-1 text-sm text-surface-500">
              Are you sure you want to delete
              <span class="font-semibold text-surface-700">
                {userToDelete?.firstName}
                {userToDelete?.lastName}
              </span>? This action cannot be undone.
            </Dialog.Description>
          </div>
        </div>

        <div class="flex items-center justify-end gap-3 pt-2">
          <Dialog.CloseTrigger
            type="button"
            class="rounded-xl px-5 py-2.5 text-sm font-medium text-surface-500 transition-all hover:bg-surface-50 hover:text-surface-900"
          >
            Cancel
          </Dialog.CloseTrigger>
          <Button
            variant="danger"
            onclick={confirmDelete}
            loading={deleteUser.isPending}
            loadingText="Deleting..."
          >
            <Trash2 size={15} />
            Delete
          </Button>
        </div>
      </Dialog.Content>
    </Dialog.Positioner>
  </Portal>
</Dialog>
