<script lang="ts">
	import {
		Search,
		LoaderCircle,
		ChevronLeft,
		ChevronRight,
		X,
		ShieldCheck,
		User,
		FilterX
	} from '@lucide/svelte';
	import { getUsersQuery } from '$lib/state/user.svelte';
	import type { UserRole } from '$lib/api/types';

	const PAGE_SIZE = 10;

	let searchInput = $state('');
	let debouncedSearch = $state('');
	let filterRole = $state<UserRole | ''>('');
	let sortField = $state('createdAt');
	let sortOrder = $state<'asc' | 'desc'>('desc');
	let currentPage = $state(1);

	let debounceTimer: ReturnType<typeof setTimeout>;
	function onSearchInput(value: string) {
		searchInput = value;
		clearTimeout(debounceTimer);
		debounceTimer = setTimeout(() => {
			debouncedSearch = value;
			currentPage = 1;
		}, 350);
	}

	const usersQuery = getUsersQuery(() => ({
		page: currentPage,
		limit: PAGE_SIZE,
		firstName: debouncedSearch || undefined,
		role: (filterRole as UserRole) || undefined,
		sort: sortField,
		order: sortOrder
	}));

	const hasActiveFilters = $derived(debouncedSearch !== '' || filterRole !== '');

	function clearFilters() {
		searchInput = '';
		debouncedSearch = '';
		filterRole = '';
		currentPage = 1;
	}

	function setSort(field: string) {
		if (sortField === field) {
			sortOrder = sortOrder === 'asc' ? 'desc' : 'asc';
		} else {
			sortField = field;
			sortOrder = 'asc';
		}
		currentPage = 1;
	}

	function goToPage(page: number) {
		currentPage = page;
	}

	function formatDate(dateStr: string) {
		return new Intl.DateTimeFormat('en-GB', {
			day: '2-digit',
			month: 'short',
			year: 'numeric'
		}).format(new Date(dateStr));
	}

	function getInitials(firstName: string, lastName: string) {
		return `${firstName[0] ?? ''}${lastName[0] ?? ''}`.toUpperCase();
	}

	const totalPages = $derived(usersQuery.data?.totalPages ?? 1);
	const totalUsers = $derived(usersQuery.data?.total ?? 0);

	function getPageNumbers(total: number, current: number): (number | '...')[] {
		if (total <= 7) return Array.from({ length: total }, (_, i) => i + 1);
		const pages: (number | '...')[] = [1];
		if (current > 3) pages.push('...');
		for (let i = Math.max(2, current - 1); i <= Math.min(total - 1, current + 1); i++) {
			pages.push(i);
		}
		if (current < total - 2) pages.push('...');
		pages.push(total);
		return pages;
	}
</script>

<svelte:head>
	<title>Users - Taskify</title>
</svelte:head>

<div class="h-full flex flex-col pt-8">
	<!-- Page Header -->
	<header class="px-8 pb-8">
		<div class="space-y-1">
			<h2 class="text-4xl text-surface-900 tracking-tight leading-tight">
				<span class="font-light">Team</span>
				<span class="font-normal"> Members</span>
			</h2>
			<p class="text-surface-800 text-xl font-light">
				{#if usersQuery.isSuccess}
					{totalUsers} {totalUsers === 1 ? 'user' : 'users'} in the system
				{:else}
					Manage and view all users
				{/if}
			</p>
		</div>
	</header>

	<!-- Filters Bar -->
	<div class="px-8 mb-6 flex flex-wrap items-center gap-3">
		<!-- Search -->
		<div class="relative flex-1 min-w-[220px] max-w-xs">
			<div class="absolute inset-y-0 left-3 flex items-center pointer-events-none text-surface-400">
				<Search size={16} />
			</div>
			<input
				type="text"
				value={searchInput}
				oninput={(e) => onSearchInput((e.target as HTMLInputElement).value)}
				placeholder="Search by first name..."
				class="w-full h-10 pl-9 pr-4 rounded-xl border border-surface-300 bg-white text-sm text-surface-900 placeholder-surface-400 focus:outline-none focus:border-primary-500 focus:ring-2 focus:ring-primary-500/10 transition-all"
			/>
		</div>

		<!-- Role Filter -->
		<select
			bind:value={filterRole}
			onchange={() => (currentPage = 1)}
			class="h-10 px-3 pr-8 rounded-xl border border-surface-300 bg-white text-sm text-surface-700 focus:outline-none focus:border-primary-500 focus:ring-2 focus:ring-primary-500/10 transition-all appearance-none cursor-pointer"
		>
			<option value="">All roles</option>
			<option value="admin">Admin</option>
			<option value="employee">Employee</option>
		</select>

		<!-- Sort -->
		<select
			bind:value={sortField}
			onchange={() => (currentPage = 1)}
			class="h-10 px-3 pr-8 rounded-xl border border-surface-300 bg-white text-sm text-surface-700 focus:outline-none focus:border-primary-500 focus:ring-2 focus:ring-primary-500/10 transition-all appearance-none cursor-pointer"
		>
			<option value="createdAt">Sort: Joined date</option>
			<option value="firstName">Sort: First name</option>
			<option value="lastName">Sort: Last name</option>
			<option value="email">Sort: Email</option>
		</select>

		<button
			onclick={() => setSort(sortField)}
			class="h-10 px-3 rounded-xl border border-surface-300 bg-white text-sm text-surface-700 hover:bg-surface-50 transition-all flex items-center gap-1.5"
			title="Toggle sort order"
		>
			{sortOrder === 'asc' ? '↑ Asc' : '↓ Desc'}
		</button>

		<!-- Clear Filters -->
		{#if hasActiveFilters}
			<button
				onclick={clearFilters}
				class="h-10 px-3 rounded-xl border border-surface-300 bg-white text-sm text-surface-500 hover:text-surface-900 hover:bg-surface-50 transition-all flex items-center gap-1.5"
			>
				<FilterX size={15} />
				Clear
			</button>
		{/if}
	</div>

	<!-- Table -->
	<div class="flex-1 overflow-y-auto px-8 pb-8">
		{#if usersQuery.isPending}
			<div class="flex flex-col items-center justify-center py-24 text-surface-400">
				<LoaderCircle size={28} class="animate-spin mb-3" />
				<span class="text-sm font-medium">Loading users...</span>
			</div>
		{:else if usersQuery.isError}
			<div class="flex flex-col items-center justify-center py-24 text-rose-500">
				<X size={28} class="mb-3" />
				<span class="text-sm font-medium">Failed to load users</span>
			</div>
		{:else if usersQuery.data}
			{@const users = usersQuery.data.data}

			{#if users.length === 0}
				<div class="flex flex-col items-center justify-center py-24 text-surface-400">
					<User size={32} class="mb-3 opacity-40" />
					<p class="text-sm font-medium">No users found</p>
					{#if hasActiveFilters}
						<button onclick={clearFilters} class="mt-2 text-xs text-primary-500 hover:underline">
							Clear filters
						</button>
					{/if}
				</div>
			{:else}
				<!-- Users Table -->
				<div class="bg-white rounded-2xl border border-surface-200 overflow-hidden">
					<table class="w-full">
						<thead>
							<tr class="border-b border-surface-100">
								<th
									class="text-left px-6 py-3.5 text-xs font-semibold text-surface-500 uppercase tracking-wide cursor-pointer hover:text-surface-700 transition-colors"
									onclick={() => setSort('firstName')}
								>
									<span class="flex items-center gap-1">
										User
										{#if sortField === 'firstName'}
											<span class="text-primary-500">{sortOrder === 'asc' ? '↑' : '↓'}</span>
										{/if}
									</span>
								</th>
								<th
									class="text-left px-6 py-3.5 text-xs font-semibold text-surface-500 uppercase tracking-wide cursor-pointer hover:text-surface-700 transition-colors"
									onclick={() => setSort('email')}
								>
									<span class="flex items-center gap-1">
										Email
										{#if sortField === 'email'}
											<span class="text-primary-500">{sortOrder === 'asc' ? '↑' : '↓'}</span>
										{/if}
									</span>
								</th>
								<th class="text-left px-6 py-3.5 text-xs font-semibold text-surface-500 uppercase tracking-wide">
									Role
								</th>
								<th class="text-left px-6 py-3.5 text-xs font-semibold text-surface-500 uppercase tracking-wide">
									Status
								</th>
								<th
									class="text-left px-6 py-3.5 text-xs font-semibold text-surface-500 uppercase tracking-wide cursor-pointer hover:text-surface-700 transition-colors"
									onclick={() => setSort('createdAt')}
								>
									<span class="flex items-center gap-1">
										Joined
										{#if sortField === 'createdAt'}
											<span class="text-primary-500">{sortOrder === 'asc' ? '↑' : '↓'}</span>
										{/if}
									</span>
								</th>
							</tr>
						</thead>
						<tbody class="divide-y divide-surface-50">
							{#each users as user (user.id)}
								<tr class="hover:bg-surface-50/60 transition-colors group">
									<!-- Avatar + Name -->
									<td class="px-6 py-4">
										<div class="flex items-center gap-3">
											{#if user.avatarUrl}
												<img
													src={user.avatarUrl}
													alt="{user.firstName} {user.lastName}"
													class="size-9 rounded-xl object-cover border border-surface-200 shrink-0"
												/>
											{:else}
												<div class="size-9 rounded-xl bg-primary-50 border border-primary-100 flex items-center justify-center shrink-0">
													<span class="text-[11px] font-bold text-primary-600">
														{getInitials(user.firstName, user.lastName)}
													</span>
												</div>
											{/if}
											<div class="min-w-0">
												<p class="text-sm font-medium text-surface-900 truncate">
													{user.firstName} {user.lastName}
												</p>
											</div>
										</div>
									</td>

									<!-- Email -->
									<td class="px-6 py-4">
										<span class="text-sm text-surface-600 truncate">{user.email}</span>
									</td>

									<!-- Role -->
									<td class="px-6 py-4">
										{#if user.role === 'admin'}
											<span class="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-[11px] font-semibold bg-violet-50 text-violet-700">
												<ShieldCheck size={11} />
												Admin
											</span>
										{:else}
											<span class="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-[11px] font-semibold bg-sky-50 text-sky-700">
												<User size={11} />
												Employee
											</span>
										{/if}
									</td>

									<!-- Status -->
									<td class="px-6 py-4">
										{#if user.isActive}
											<span class="inline-flex items-center gap-1.5 text-[11px] font-semibold text-emerald-600">
												<span class="size-1.5 rounded-full bg-emerald-500"></span>
												Active
											</span>
										{:else}
											<span class="inline-flex items-center gap-1.5 text-[11px] font-semibold text-surface-400">
												<span class="size-1.5 rounded-full bg-surface-300"></span>
												Inactive
											</span>
										{/if}
									</td>

									<!-- Joined -->
									<td class="px-6 py-4">
										<span class="text-sm text-surface-500">{formatDate(user.createdAt)}</span>
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
								{(currentPage - 1) * PAGE_SIZE + 1}–{Math.min(currentPage * PAGE_SIZE, totalUsers)}
							</span>
							of <span class="font-medium text-surface-700">{totalUsers}</span> users
						</p>

						<div class="flex items-center gap-1">
							<button
								onclick={() => goToPage(currentPage - 1)}
								disabled={currentPage === 1}
								class="size-9 flex items-center justify-center rounded-xl border border-surface-200 text-surface-500 hover:bg-surface-50 disabled:opacity-40 disabled:cursor-not-allowed transition-all"
							>
								<ChevronLeft size={16} />
							</button>

							{#each getPageNumbers(totalPages, currentPage) as page, i (i)}
								{#if page === '...'}
									<span class="size-9 flex items-center justify-center text-sm text-surface-400">…</span>
								{:else}
									<button
										onclick={() => goToPage(page as number)}
										class="size-9 flex items-center justify-center rounded-xl border text-sm font-medium transition-all
										{currentPage === page
											? 'bg-primary-500 border-primary-500 text-white shadow-sm'
											: 'border-surface-200 text-surface-600 hover:bg-surface-50'}"
									>
										{page}
									</button>
								{/if}
							{/each}

							<button
								onclick={() => goToPage(currentPage + 1)}
								disabled={currentPage === totalPages}
								class="size-9 flex items-center justify-center rounded-xl border border-surface-200 text-surface-500 hover:bg-surface-50 disabled:opacity-40 disabled:cursor-not-allowed transition-all"
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
