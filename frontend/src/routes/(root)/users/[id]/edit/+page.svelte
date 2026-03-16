<script lang="ts">
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { ArrowLeft, LoaderCircle, User, ShieldCheck } from '@lucide/svelte';
	import { createUserQuery, updateUserMutation } from '$lib/state/user.svelte';
	import { toaster } from '$lib/state/toast.svelte';
	import Input from '$lib/components/Input.svelte';

	const userId = $derived(page.params.id ?? '');

	const userQuery = createUserQuery(() => userId);
	const updateUser = updateUserMutation();

	let firstName = $state('');
	let lastName = $state('');
	let isActive = $state(true);

	$effect(() => {
		if (userQuery.data) {
			firstName = userQuery.data.firstName;
			lastName = userQuery.data.lastName;
			isActive = userQuery.data.isActive;
		}
	});

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault();

		try {
			await updateUser.mutateAsync({
				id: userId,
				data: { firstName, lastName, isActive }
			});

			toaster.success({
				title: 'User Updated',
				description: `${firstName} ${lastName}'s profile has been updated successfully.`
			});

			goto(resolve('/users'));
		} catch {
			toaster.error({
				title: 'Update Failed',
				description: 'Could not update user. Please try again.'
			});
		}
	}

	function formatDate(dateStr: string) {
		return new Intl.DateTimeFormat('en-GB', {
			day: '2-digit',
			month: 'long',
			year: 'numeric'
		}).format(new Date(dateStr));
	}
</script>

<svelte:head>
	<title>Edit User - Taskify</title>
</svelte:head>

<div class="h-full flex flex-col pt-8">
	<header class="px-8 pb-8">
		<button
			onclick={() => goto(resolve('/users'))}
			class="flex items-center gap-2 text-sm text-surface-500 hover:text-surface-800 transition-colors mb-6"
		>
			<ArrowLeft size={16} />
			Back to Users
		</button>

		<h2 class="text-4xl text-surface-900 tracking-tight leading-tight">
			<span class="font-light">Edit</span>
			<span class="font-normal"> User</span>
		</h2>
	</header>

	<div class="px-8 flex-1 overflow-y-auto pb-8">
		{#if userQuery.isPending}
			<div class="flex flex-col items-center justify-center py-24 text-surface-400">
				<LoaderCircle size={28} class="animate-spin mb-3" />
				<span class="text-sm font-medium">Loading user...</span>
			</div>
		{:else if userQuery.isError}
			<div class="flex flex-col items-center justify-center py-24 text-rose-500">
				<p class="text-sm font-medium">Failed to load user.</p>
			</div>
		{:else if userQuery.data}
			{@const user = userQuery.data}
			<div class="max-w-2xl space-y-6">

				<!-- Avatar / Identity card (readonly) -->
				<div class="bg-white rounded-2xl border border-surface-200 p-6 flex items-center gap-4">
					{#if user.avatarUrl}
						<img
							src={user.avatarUrl}
							alt="{user.firstName} {user.lastName}"
							class="size-14 rounded-xl object-cover border border-surface-200 shrink-0"
						/>
					{:else}
						<div class="size-14 rounded-xl bg-primary-50 border border-primary-100 flex items-center justify-center shrink-0">
							<span class="text-base font-bold text-primary-600">
								{user.firstName[0]}{user.lastName[0]}
							</span>
						</div>
					{/if}
					<div>
						<p class="font-semibold text-surface-900">{user.firstName} {user.lastName}</p>
						<p class="text-sm text-surface-500">{user.email}</p>
					</div>
					{#if user.role === 'admin'}
						<span class="ml-auto inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-[11px] font-semibold bg-violet-50 text-violet-700">
							<ShieldCheck size={11} />
							Admin
						</span>
					{:else}
						<span class="ml-auto inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-[11px] font-semibold bg-sky-50 text-sky-700">
							<User size={11} />
							Employee
						</span>
					{/if}
				</div>

				<!-- Edit form -->
				<form onsubmit={handleSubmit} class="bg-white rounded-2xl border border-surface-200 p-6 space-y-5">
					<h3 class="text-sm font-semibold text-surface-900 mb-1">Editable Information</h3>

					<div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
						<Input
							id="firstName"
							label="First Name"
							placeholder="First name"
							bind:value={firstName}
							required
						/>
						<Input
							id="lastName"
							label="Last Name"
							placeholder="Last name"
							bind:value={lastName}
							required
						/>
					</div>

					<!-- isActive toggle -->
					<div class="flex items-center justify-between p-4 rounded-xl border border-surface-200 bg-surface-50">
						<div>
							<p class="text-sm font-medium text-surface-800">Active Account</p>
							<p class="text-xs text-surface-500 mt-0.5">Inactive users cannot log in</p>
						</div>
						<button
							type="button"
							aria-label="Toggle active status"
							onclick={() => (isActive = !isActive)}
							class="relative inline-flex h-6 w-11 shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 focus:outline-none {isActive ? 'bg-primary-500' : 'bg-surface-300'}"
							role="switch"
							aria-checked={isActive}
						>
							<span
								class="pointer-events-none inline-block size-5 rounded-full bg-white shadow-sm transform transition duration-200 {isActive ? 'translate-x-5' : 'translate-x-0'}"
							></span>
						</button>
					</div>

					<hr class="border-surface-100" />

					<h3 class="text-sm font-semibold text-surface-900 mb-1">Read-only Information</h3>

					<div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
						<!-- Email readonly -->
						<div class="space-y-1">
							<label for="ro-email" class="block text-sm font-medium text-surface-700">Email</label>
							<input
								id="ro-email"
								type="text"
								value={user.email}
								disabled
								class="block h-12 w-full rounded-xl border border-surface-200 bg-surface-100 px-4 text-surface-500 text-sm cursor-not-allowed"
							/>
						</div>

						<!-- Role readonly -->
						<div class="space-y-1">
							<label for="ro-role" class="block text-sm font-medium text-surface-700">Role</label>
							<input
								id="ro-role"
								type="text"
								value={user.role.charAt(0).toUpperCase() + user.role.slice(1)}
								disabled
								class="block h-12 w-full rounded-xl border border-surface-200 bg-surface-100 px-4 text-surface-500 text-sm cursor-not-allowed"
							/>
						</div>

						<!-- Joined readonly -->
						<div class="space-y-1">
							<label for="ro-joined" class="block text-sm font-medium text-surface-700">Joined</label>
							<input
								id="ro-joined"
								type="text"
								value={formatDate(user.createdAt)}
								disabled
								class="block h-12 w-full rounded-xl border border-surface-200 bg-surface-100 px-4 text-surface-500 text-sm cursor-not-allowed"
							/>
						</div>
					</div>

					<!-- Footer -->
					<div class="flex items-center justify-end gap-3 pt-2">
						<button
							type="button"
							onclick={() => goto(resolve('/users'))}
							class="px-6 py-2.5 rounded-xl font-medium text-surface-500 hover:text-surface-900 hover:bg-surface-50 transition-all text-sm"
						>
							Cancel
						</button>
						<button
							type="submit"
							disabled={updateUser.isPending}
							class="bg-primary-500 hover:bg-primary-700 text-white px-8 py-2.5 rounded-xl font-bold shadow-lg shadow-primary-500/20 transition-all active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2 text-sm"
						>
							{#if updateUser.isPending}
								<LoaderCircle size={16} class="animate-spin" />
								Saving...
							{:else}
								Save Changes
							{/if}
						</button>
					</div>
				</form>
			</div>
		{/if}
	</div>
</div>
