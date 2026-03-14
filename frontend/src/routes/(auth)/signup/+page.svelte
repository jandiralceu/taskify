<script lang="ts">
	import { z } from 'zod';
	import { Mail, Lock, User, Users, Shield, Check } from '@lucide/svelte';

	import { resolve } from '$app/paths';
	import Input from '$lib/components/Input.svelte';
	import logo from '$lib/assets/logo.webp';

	let firstName = $state('');
	let lastName = $state('');
	let email = $state('');
	let password = $state('');
	let role = $state<'admin' | 'employee'>('employee');
	let submitted = $state(false);

	const signupSchema = z.object({
		firstName: z.string().min(2, 'Must be at least 2 characters').max(100, 'Max 100 characters'),
		lastName: z.string().min(2, 'Must be at least 2 characters').max(100, 'Max 100 characters'),
		email: z.string().email('Please enter a valid email address').max(255, 'Max 255 characters'),
		password: z.string().min(8, 'Password must be at least 8 characters')
	});

	let validationResult = $derived(signupSchema.safeParse({ firstName, lastName, email, password }));
	let fieldErrors = $derived(submitted && !validationResult.success ? z.treeifyError(validationResult.error).properties : undefined);

	function handleSubmit(e: SubmitEvent) {
		e.preventDefault();
		submitted = true;

		if (!validationResult.success) {
			return;
		}

		console.log('Signup attempt:', { firstName, lastName, email, password, role });
		// Ready to send to API
	}
</script>

<form class="w-full max-w-md mt-8 space-y-6" onsubmit={handleSubmit} novalidate>
	<h1>
		<img src={logo} alt="Taskify Logo" class="h-12 w-auto" />
	</h1>

	<div class="space-y-1 pb-4">
		<h2 class="text-primary-500 text-3xl font-bold">Create your account</h2>
		<p class="text-surface-600 text-sm dark:text-surface-400">Join Taskify today and start managing your tasks efficiently.</p>
	</div>

	<div class="space-y-4">
		<div class="grid grid-cols-2 gap-4">
			<Input
				id="first-name"
				name="first_name"
				placeholder="First Name"
				bind:value={firstName}
				error={fieldErrors?.firstName?.errors?.[0]}
				required
			>
				{#snippet icon()}
					<User size={20} />
				{/snippet}
			</Input>

			<Input
				id="last-name"
				name="last_name"
				placeholder="Last Name"
				bind:value={lastName}
				error={fieldErrors?.lastName?.errors?.[0]}
				required
			>
				{#snippet icon()}
					<User size={20} />
				{/snippet}
			</Input>
		</div>

		<Input
			id="email-address"
			name="email"
			type="email"
			placeholder="Email"
			bind:value={email}
			error={fieldErrors?.email?.errors?.[0]}
			required
		>
			{#snippet icon()}
				<Mail size={20} />
			{/snippet}
		</Input>

		<Input
			id="password"
			name="password"
			type="password"
			placeholder="Password"
			bind:value={password}
			error={fieldErrors?.password?.errors?.[0]}
			required
		>
			{#snippet icon()}
				<Lock size={20} />
			{/snippet}
		</Input>

		<div class="space-y-3">
			<span class="block text-sm font-medium text-surface-700 dark:text-surface-300">Select your Role</span>
			<div class="flex flex-col gap-4">
				<label class="relative flex cursor-pointer rounded-xl border p-4 transition-all {role === 'employee' ? 'border-primary-500 bg-primary-50 dark:border-primary-500 dark:bg-primary-900/20' : 'border-surface-200 bg-white hover:border-surface-300 dark:border-surface-700 dark:bg-surface-900 dark:hover:border-surface-600'}">
					<input
						type="radio"
						name="role"
						value="employee"
						bind:group={role}
						class="sr-only"
					/>
					<div class="flex w-full items-center gap-4">
						<div class="flex h-12 w-12 shrink-0 items-center justify-center rounded-lg transition-colors {role === 'employee' ? 'bg-primary-100 text-primary-600 dark:bg-primary-900/50 dark:text-primary-400' : 'bg-surface-100 text-surface-500 dark:bg-surface-800'}">
							<Users size={24} />
						</div>
						<div class="flex-1">
							<div class="font-semibold text-surface-900 dark:text-surface-50">Employee</div>
							<div class="text-xs text-surface-500 dark:text-surface-400 mt-0.5">Regular access to tasks and projects.</div>
						</div>
						<div class="flex h-6 w-6 shrink-0 items-center justify-center rounded-full border-2 transition-colors {role === 'employee' ? 'border-primary-500 bg-primary-500 text-white' : 'border-surface-300 dark:border-surface-600'}">
							{#if role === 'employee'}
								<Check size={14} strokeWidth={3} />
							{/if}
						</div>
					</div>
				</label>

				<label class="relative flex cursor-pointer rounded-xl border p-4 transition-all {role === 'admin' ? 'border-primary-500 bg-primary-50 dark:border-primary-500 dark:bg-primary-900/20' : 'border-surface-200 bg-white hover:border-surface-300 dark:border-surface-700 dark:bg-surface-900 dark:hover:border-surface-600'}">
					<input
						type="radio"
						name="role"
						value="admin"
						bind:group={role}
						class="sr-only"
					/>
					<div class="flex w-full items-center gap-4">
						<div class="flex h-12 w-12 shrink-0 items-center justify-center rounded-lg transition-colors {role === 'admin' ? 'bg-primary-100 text-primary-600 dark:bg-primary-900/50 dark:text-primary-400' : 'bg-surface-100 text-surface-500 dark:bg-surface-800'}">
							<Shield size={24} />
						</div>
						<div class="flex-1">
							<div class="font-semibold text-surface-900 dark:text-surface-50">Admin</div>
							<div class="text-xs text-surface-500 dark:text-surface-400 mt-0.5">Full access to manage users and system settings.</div>
						</div>
						<div class="flex h-6 w-6 shrink-0 items-center justify-center rounded-full border-2 transition-colors {role === 'admin' ? 'border-primary-500 bg-primary-500 text-white' : 'border-surface-300 dark:border-surface-600'}">
							{#if role === 'admin'}
								<Check size={14} strokeWidth={3} />
							{/if}
						</div>
					</div>
				</label>
			</div>
		</div>
	</div>

	<div>
		<button
			type="submit"
			class="group relative flex w-full justify-center rounded-xl border border-transparent bg-primary-500 px-4 py-3 text-sm font-semibold text-white transition-all hover:bg-primary-600 focus:ring-2 focus:ring-primary-500/50 focus:outline-none active:scale-[0.98]"
		>
			Create Account
		</button>
	</div>

	<div class="flex items-center justify-center gap-1 text-sm text-surface-600 dark:text-surface-400">
		<span>Already have an account?</span>
		<a
			href={resolve('/signin')}
			class="font-medium text-primary-600 hover:text-primary-500 dark:text-primary-400"
		>
			Sign in
		</a>
	</div>
</form>
