<script lang="ts">
	import { Dialog, Portal, Progress } from '@skeletonlabs/skeleton-svelte';
	import { X, Mail, Shield, User as UserIcon, Calendar, CircleCheck, Clock } from '@lucide/svelte';
	import { createUserQuery } from '$lib/state/user.svelte';

	interface Props {
		userId: string | undefined;
		open: boolean;
		onOpenChange: (open: boolean) => void;
	}

	let { userId, open, onOpenChange }: Props = $props();

	const userQuery = createUserQuery(() => userId);
	const user = $derived(userQuery.data);
	const isLoading = $derived(userQuery.isLoading);

	function formatDate(dateStr: string | undefined) {
		if (!dateStr) return 'N/A';
		return new Intl.DateTimeFormat('en-GB', {
			day: '2-digit',
			month: 'long',
			year: 'numeric'
		}).format(new Date(dateStr));
	}
</script>

<Dialog {open} onOpenChange={(e) => onOpenChange(e.open)} closeOnInteractOutside={true}>
	<Portal>
		<Dialog.Backdrop class="fixed inset-0 z-50 bg-surface-950/30 backdrop-blur-[2px] transition-opacity" />
		<Dialog.Positioner class="fixed inset-y-0 right-0 z-50 flex items-stretch justify-end">
			<Dialog.Content class="w-full max-w-[400px] bg-white shadow-2xl flex flex-col overflow-y-auto animate-slide-in">
				<div class="p-6 border-b border-surface-100 flex justify-between items-center bg-surface-50">
					<h2 class="text-xl font-bold text-surface-900 flex items-center gap-2">
						<UserIcon size={20} class="text-primary-500" />
						User Profile
					</h2>
					<button
						onclick={() => onOpenChange(false)}
						class="p-2 hover:bg-surface-200 rounded-full transition-colors text-surface-500"
					>
						<X size={20} />
					</button>
				</div>

				<div class="flex-1 p-8">
					{#if isLoading}
						<div class="h-full flex flex-col items-center justify-center gap-4 py-20">
							<Progress value={null} class="w-fit">
								<Progress.Circle class="size-12">
									<Progress.CircleTrack />
									<Progress.CircleRange class="stroke-primary-500" />
								</Progress.Circle>
							</Progress>
							<p class="text-surface-500 animate-pulse font-medium">Loading user data...</p>
						</div>
					{:else if user}
						<!-- Header Section -->
						<div class="flex flex-col items-center text-center mb-8">
							<div class="relative mb-4">
								{#if user.avatarUrl}
									<img
										src={user.avatarUrl}
										alt={user.firstName}
										class="size-32 rounded-3xl object-cover shadow-xl border-4 border-white"
									/>
								{:else}
									<div class="size-32 rounded-3xl bg-primary-100 flex items-center justify-center shadow-xl border-4 border-white">
										<span class="text-4xl font-bold text-primary-600">
											{user.firstName[0]}{user.lastName[0]}
										</span>
									</div>
								{/if}
								<div class="absolute -bottom-2 -right-2 p-1.5 bg-white rounded-full shadow-md">
									{#if user.isActive}
										<CircleCheck size={24} class="text-emerald-500" />
									{:else}
										<Clock size={24} class="text-surface-400" />
									{/if}
								</div>
							</div>
							<h3 class="text-2xl font-bold text-surface-900">
								{user.firstName} {user.lastName}
							</h3>
							<span class="inline-flex items-center px-3 py-1 mt-2 rounded-full text-xs font-semibold bg-primary-50 text-primary-700 capitalize">
								{user.role}
							</span>
						</div>

						<!-- Info Cards -->
						<div class="space-y-4">
							<div class="bg-surface-50 p-4 rounded-2xl border border-surface-100">
								<div class="flex items-center gap-3 text-surface-500 mb-1">
									<Mail size={16} />
									<span class="text-xs font-medium uppercase tracking-wider">Email Address</span>
								</div>
								<p class="text-surface-900 font-medium break-all">{user.email}</p>
							</div>

							<div class="bg-surface-50 p-4 rounded-2xl border border-surface-100">
								<div class="flex items-center gap-3 text-surface-500 mb-1">
									<Shield size={16} />
									<span class="text-xs font-medium uppercase tracking-wider">Role & Permissions</span>
								</div>
								<p class="text-surface-900 font-medium capitalize">{user.role}</p>
							</div>

							<div class="bg-surface-50 p-4 rounded-2xl border border-surface-100">
								<div class="flex items-center gap-3 text-surface-500 mb-1">
									<Calendar size={16} />
									<span class="text-xs font-medium uppercase tracking-wider">Member Since</span>
								</div>
								<p class="text-surface-900 font-medium">{formatDate(user.createdAt)}</p>
							</div>
						</div>

						<!-- Stats section -->
						<div class="grid grid-cols-2 gap-4 mt-8">
							<div class="bg-indigo-50 p-4 rounded-2xl text-center">
								<span class="block text-2xl font-bold text-indigo-700">--</span>
								<span class="text-xs font-medium text-indigo-600 uppercase">Tasks Active</span>
							</div>
							<div class="bg-emerald-50 p-4 rounded-2xl text-center">
								<span class="block text-2xl font-bold text-emerald-700">--</span>
								<span class="text-xs font-medium text-emerald-600 uppercase">Completed</span>
							</div>
						</div>
					{/if}
				</div>

				<div class="p-6 border-t border-surface-100 bg-surface-50 mt-auto">
					<button
						onclick={() => onOpenChange(false)}
						class="w-full py-3 bg-white border border-surface-200 text-surface-700 font-semibold rounded-xl hover:bg-surface-100 transition-colors shadow-sm"
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
