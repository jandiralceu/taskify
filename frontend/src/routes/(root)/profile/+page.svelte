<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { LoaderCircle, Camera, ShieldCheck, User, KeyRound, Trash2, TriangleAlert, Eye, EyeOff } from '@lucide/svelte';
	import { Dialog, Portal } from '@skeletonlabs/skeleton-svelte';
	import {
		createProfileQuery,
		updateProfileMutation,
		changePasswordMutation,
		deleteProfileMutation
	} from '$lib/state/user.svelte';
	import { toaster } from '$lib/state/toast.svelte';
	import Input from '$lib/components/Input.svelte';

	const profileQuery = createProfileQuery();
	const updateProfile = updateProfileMutation();
	const changePassword = changePasswordMutation();
	const deleteProfile = deleteProfileMutation();

	// Edit info state
	let firstName = $state('');
	let lastName = $state('');

	// Change password state
	let oldPassword = $state('');
	let newPassword = $state('');
	let confirmPassword = $state('');
	let showOldPassword = $state(false);
	let showNewPassword = $state(false);

	// Delete confirmation
	let isDeleteOpen = $state(false);
	let deleteConfirmText = $state('');
	const DELETE_PHRASE = 'delete my account';

	$effect(() => {
		if (profileQuery.data) {
			firstName = profileQuery.data.firstName;
			lastName = profileQuery.data.lastName;
		}
	});

	function formatDate(dateStr: string) {
		return new Intl.DateTimeFormat('en-GB', {
			day: '2-digit',
			month: 'long',
			year: 'numeric'
		}).format(new Date(dateStr));
	}

	async function handleUpdateInfo(e: SubmitEvent) {
		e.preventDefault();
		try {
			await updateProfile.mutateAsync({ firstName, lastName });
			toaster.success({
				title: 'Profile Updated',
				description: 'Your personal information has been saved.'
			});
		} catch {
			toaster.error({
				title: 'Update Failed',
				description: 'Could not update your profile. Please try again.'
			});
		}
	}

	async function handleChangePassword(e: SubmitEvent) {
		e.preventDefault();
		if (newPassword !== confirmPassword) {
			toaster.error({ title: 'Password Mismatch', description: 'New passwords do not match.' });
			return;
		}
		try {
			await changePassword.mutateAsync({ oldPassword, newPassword });
			oldPassword = '';
			newPassword = '';
			confirmPassword = '';
			toaster.success({
				title: 'Password Changed',
				description: 'Your password has been updated successfully.'
			});
		} catch {
			toaster.error({
				title: 'Change Failed',
				description: 'Could not change your password. Check your current password and try again.'
			});
		}
	}

	async function handleDeleteAccount() {
		try {
			await deleteProfile.mutateAsync();
			goto(resolve('/signin'));
		} catch {
			toaster.error({
				title: 'Delete Failed',
				description: 'Could not delete your account. Please try again.'
			});
		}
	}

	const canConfirmDelete = $derived(deleteConfirmText.toLowerCase() === DELETE_PHRASE);
</script>

<svelte:head>
	<title>My Profile - Taskify</title>
</svelte:head>

<div class="h-full flex flex-col pt-8">
	<header class="px-8 pb-8">
		<h2 class="text-4xl text-surface-900 tracking-tight leading-tight">
			<span class="font-light">My</span>
			<span class="font-normal"> Profile</span>
		</h2>
	</header>

	<div class="px-8 flex-1 overflow-y-auto pb-8">
		{#if profileQuery.isPending}
			<div class="flex flex-col items-center justify-center py-24 text-surface-400">
				<LoaderCircle size={28} class="animate-spin mb-3" />
				<span class="text-sm font-medium">Loading profile...</span>
			</div>
		{:else if profileQuery.isError}
			<div class="flex flex-col items-center justify-center py-24 text-rose-500">
				<p class="text-sm font-medium">Failed to load profile.</p>
			</div>
		{:else if profileQuery.data}
			{@const user = profileQuery.data}
			<div class="max-w-2xl space-y-6">

				<!-- Identity card -->
				<div class="bg-white rounded-2xl border border-surface-200 p-6 flex items-center gap-5">
					<!-- Avatar -->
					<div class="relative shrink-0">
						{#if user.avatarUrl}
							<img
								src={user.avatarUrl}
								alt="{user.firstName} {user.lastName}"
								class="size-20 rounded-2xl object-cover border border-surface-200"
							/>
						{:else}
							<div class="size-20 rounded-2xl bg-primary-50 border border-primary-100 flex items-center justify-center">
								<span class="text-2xl font-bold text-primary-600">
									{user.firstName[0]}{user.lastName[0]}
								</span>
							</div>
						{/if}
						<button
							type="button"
							title="Change avatar (coming soon)"
							class="absolute -bottom-2 -right-2 size-7 rounded-lg bg-white border border-surface-200 shadow-sm flex items-center justify-center text-surface-400 hover:text-primary-500 transition-colors cursor-not-allowed opacity-50"
						>
							<Camera size={13} />
						</button>
					</div>

					<div class="flex-1 min-w-0">
						<p class="text-lg font-semibold text-surface-900 truncate">
							{user.firstName} {user.lastName}
						</p>
						<p class="text-sm text-surface-500 truncate">{user.email}</p>
						<div class="flex items-center gap-2 mt-2">
							{#if user.role === 'admin'}
								<span class="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-[11px] font-semibold bg-violet-50 text-violet-700">
									<ShieldCheck size={11} /> Admin
								</span>
							{:else}
								<span class="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-[11px] font-semibold bg-sky-50 text-sky-700">
									<User size={11} /> Employee
								</span>
							{/if}
							<span class="text-xs text-surface-400">Member since {formatDate(user.createdAt)}</span>
						</div>
					</div>
				</div>

				<!-- Edit personal info -->
				<form
					onsubmit={handleUpdateInfo}
					class="bg-white rounded-2xl border border-surface-200 p-6 space-y-5"
				>
					<div class="flex items-center gap-2 mb-1">
						<User size={16} class="text-surface-400" />
						<h3 class="text-sm font-semibold text-surface-900">Personal Information</h3>
					</div>

					<div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
						<Input id="firstName" label="First Name" placeholder="First name" bind:value={firstName} required />
						<Input id="lastName" label="Last Name" placeholder="Last name" bind:value={lastName} required />
					</div>

					<!-- Read-only fields -->
					<div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
						<div class="space-y-1">
							<label for="ro-email" class="block text-sm font-medium text-surface-700">Email</label>
							<input
								id="ro-email"
								type="text"
								value={user.email}
								disabled
								class="block h-12 w-full rounded-xl border border-surface-200 bg-surface-100 px-4 text-surface-400 text-sm cursor-not-allowed"
							/>
						</div>
						<div class="space-y-1">
							<label for="ro-role" class="block text-sm font-medium text-surface-700">Role</label>
							<input
								id="ro-role"
								type="text"
								value={user.role.charAt(0).toUpperCase() + user.role.slice(1)}
								disabled
								class="block h-12 w-full rounded-xl border border-surface-200 bg-surface-100 px-4 text-surface-400 text-sm cursor-not-allowed"
							/>
						</div>
					</div>

					<div class="flex justify-end pt-1">
						<button
							type="submit"
							disabled={updateProfile.isPending}
							class="bg-primary-500 hover:bg-primary-700 text-white px-8 py-2.5 rounded-xl font-bold shadow-lg shadow-primary-500/20 transition-all active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2 text-sm"
						>
							{#if updateProfile.isPending}
								<LoaderCircle size={15} class="animate-spin" />
								Saving...
							{:else}
								Save Changes
							{/if}
						</button>
					</div>
				</form>

				<!-- Change password -->
				<form
					onsubmit={handleChangePassword}
					class="bg-white rounded-2xl border border-surface-200 p-6 space-y-5"
				>
					<div class="flex items-center gap-2 mb-1">
						<KeyRound size={16} class="text-surface-400" />
						<h3 class="text-sm font-semibold text-surface-900">Change Password</h3>
					</div>

					<div class="space-y-1">
						<label for="oldPassword" class="block text-sm font-medium text-surface-700">Current Password</label>
						<div class="relative">
							<input
								id="oldPassword"
								type={showOldPassword ? 'text' : 'password'}
								placeholder="Current password"
								bind:value={oldPassword}
								required
								class="block h-12 w-full rounded-xl border border-surface-300 bg-surface-50 px-4 pr-10 text-sm text-surface-900 placeholder-surface-500 focus:outline-none focus:border-primary-500 focus:ring-2 focus:ring-primary-500/20 transition-all"
							/>
							<button
								type="button"
								onclick={() => (showOldPassword = !showOldPassword)}
								class="absolute inset-y-0 right-0 flex items-center pr-3 text-surface-400 hover:text-primary-500 transition-colors"
								tabindex="-1"
							>
								{#if showOldPassword}<Eye size={18} />{:else}<EyeOff size={18} />{/if}
							</button>
						</div>
					</div>

					<div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
						<div class="space-y-1">
							<label for="newPassword" class="block text-sm font-medium text-surface-700">New Password</label>
							<div class="relative">
								<input
									id="newPassword"
									type={showNewPassword ? 'text' : 'password'}
									placeholder="New password"
									bind:value={newPassword}
									required
									minlength={8}
									class="block h-12 w-full rounded-xl border border-surface-300 bg-surface-50 px-4 pr-10 text-sm text-surface-900 placeholder-surface-500 focus:outline-none focus:border-primary-500 focus:ring-2 focus:ring-primary-500/20 transition-all"
								/>
								<button
									type="button"
									onclick={() => (showNewPassword = !showNewPassword)}
									class="absolute inset-y-0 right-0 flex items-center pr-3 text-surface-400 hover:text-primary-500 transition-colors"
									tabindex="-1"
								>
									{#if showNewPassword}<Eye size={18} />{:else}<EyeOff size={18} />{/if}
								</button>
							</div>
						</div>

						<div class="space-y-1">
							<label for="confirmPassword" class="block text-sm font-medium text-surface-700">Confirm New Password</label>
							<input
								id="confirmPassword"
								type="password"
								placeholder="Confirm new password"
								bind:value={confirmPassword}
								required
								class="block h-12 w-full rounded-xl border transition-all
								{confirmPassword && confirmPassword !== newPassword
									? 'border-rose-400 focus:border-rose-400 focus:ring-rose-400/20'
									: 'border-surface-300 focus:border-primary-500 focus:ring-primary-500/20'}
								bg-surface-50 px-4 text-sm text-surface-900 placeholder-surface-500 focus:outline-none focus:ring-2"
							/>
							{#if confirmPassword && confirmPassword !== newPassword}
								<p class="text-xs text-rose-500 font-medium mt-1">Passwords do not match</p>
							{/if}
						</div>
					</div>

					<div class="flex justify-end pt-1">
						<button
							type="submit"
							disabled={changePassword.isPending || (!!confirmPassword && confirmPassword !== newPassword)}
							class="bg-primary-500 hover:bg-primary-700 text-white px-8 py-2.5 rounded-xl font-bold shadow-lg shadow-primary-500/20 transition-all active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2 text-sm"
						>
							{#if changePassword.isPending}
								<LoaderCircle size={15} class="animate-spin" />
								Updating...
							{:else}
								Update Password
							{/if}
						</button>
					</div>
				</form>

				<!-- Danger zone -->
				<div class="bg-white rounded-2xl border border-rose-200 p-6 space-y-4">
					<div class="flex items-center gap-2">
						<TriangleAlert size={16} class="text-rose-500" />
						<h3 class="text-sm font-semibold text-rose-600">Danger Zone</h3>
					</div>

					<div class="flex items-center justify-between p-4 rounded-xl border border-rose-100 bg-rose-50">
						<div>
							<p class="text-sm font-medium text-surface-800">Delete Account</p>
							<p class="text-xs text-surface-500 mt-0.5">
								Permanently remove your account and all associated data. This cannot be undone.
							</p>
						</div>
						<button
							type="button"
							onclick={() => { deleteConfirmText = ''; isDeleteOpen = true; }}
							class="ml-4 shrink-0 px-4 py-2 rounded-xl text-sm font-bold text-rose-600 border border-rose-300 hover:bg-rose-500 hover:text-white hover:border-rose-500 transition-all"
						>
							Delete Account
						</button>
					</div>
				</div>

			</div>
		{/if}
	</div>
</div>

<!-- Delete Account Confirmation Dialog -->
<Dialog
	role="alertdialog"
	open={isDeleteOpen}
	onOpenChange={(e) => { if (!e.open) isDeleteOpen = false; }}
>
	<Portal>
		<Dialog.Backdrop class="fixed inset-0 z-50 bg-surface-950/40 backdrop-blur-sm" />
		<Dialog.Positioner class="fixed inset-0 z-50 flex items-center justify-center p-4">
			<Dialog.Content class="w-full max-w-md bg-white rounded-2xl shadow-2xl border border-surface-100 p-6 space-y-4">
				<div class="flex items-start gap-4">
					<div class="shrink-0 size-10 rounded-xl bg-rose-50 flex items-center justify-center text-rose-500">
						<Trash2 size={20} />
					</div>
					<div>
						<Dialog.Title class="text-base font-bold text-surface-900">
							Delete Your Account
						</Dialog.Title>
						<Dialog.Description class="text-sm text-surface-500 mt-1">
							This will permanently delete your account. To confirm, type
							<span class="font-semibold text-surface-700 select-all">"{DELETE_PHRASE}"</span> below.
						</Dialog.Description>
					</div>
				</div>

				<input
					type="text"
					placeholder={DELETE_PHRASE}
					bind:value={deleteConfirmText}
					class="block w-full h-11 rounded-xl border border-surface-300 bg-surface-50 px-4 text-sm text-surface-900 placeholder-surface-400 focus:outline-none focus:border-rose-400 focus:ring-2 focus:ring-rose-400/20 transition-all"
				/>

				<div class="flex items-center justify-end gap-3 pt-1">
					<Dialog.CloseTrigger
						type="button"
						class="px-5 py-2.5 rounded-xl text-sm font-medium text-surface-500 hover:text-surface-900 hover:bg-surface-50 transition-all"
					>
						Cancel
					</Dialog.CloseTrigger>
					<button
						type="button"
						onclick={handleDeleteAccount}
						disabled={!canConfirmDelete || deleteProfile.isPending}
						class="px-5 py-2.5 rounded-xl text-sm font-bold bg-rose-500 hover:bg-rose-600 text-white transition-all active:scale-95 disabled:opacity-40 disabled:cursor-not-allowed flex items-center gap-2"
					>
						{#if deleteProfile.isPending}
							<LoaderCircle size={15} class="animate-spin" />
							Deleting...
						{:else}
							<Trash2 size={15} />
							Delete My Account
						{/if}
					</button>
				</div>
			</Dialog.Content>
		</Dialog.Positioner>
	</Portal>
</Dialog>
