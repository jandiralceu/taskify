<script lang="ts">
	import logo from '$lib/assets/logo.webp';
	import { resolve } from '$app/paths';
	import { Mail, Lock } from '@lucide/svelte';
	import Input from '$lib/components/Input.svelte';
	import { z } from 'zod';

	let email = $state('');
	let password = $state('');
	let submitted = $state(false);

	const signinSchema = z.object({
		email: z.email('Please enter a valid email address'),
		password: z.string().min(1, 'Password is required')
	});

	let validationResult = $derived(signinSchema.safeParse({ email, password }));
	let fieldErrors = $derived(submitted && !validationResult.success ? z.treeifyError(validationResult.error).properties : undefined);

	function handleSubmit(e: SubmitEvent) {
		e.preventDefault();
		submitted = true;

		if (!validationResult.success) {
			return;
		}

		console.log('Login attempt:', { email, password });
	}
</script>

<form class="mt-8 space-y-6 w-full max-w-md" onsubmit={handleSubmit} novalidate>
	<h1>
		<img src={logo} alt="Taskify Logo" class="h-12 w-auto" />
	</h1>

	<div class="space-y-1 pb-4">
		<h2 class="text-primary-500 text-3xl font-bold">Sign in to your account</h2>
		<p class="text-surface-600 text-sm dark:text-surface-400">Welcome back! Please sign in to your account.</p>
	</div>

	<div class="space-y-4">
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
	</div>

	<div>
		<button
			type="submit"
			class="group relative flex w-full justify-center rounded-xl border border-transparent bg-primary-500 px-4 py-3 text-sm font-semibold text-white transition-all hover:bg-primary-600 focus:ring-2 focus:ring-primary-500/50 focus:outline-none active:scale-[0.98]"
		>
			Sign in
		</button>
	</div>

	<div class="flex items-center justify-center gap-1 text-sm text-surface-600 dark:text-surface-400">
		<span>Don't have an account?</span>
		<a
			href={resolve('/signup')}
			class="font-medium text-primary-600 hover:text-primary-500 dark:text-primary-400"
		>
			Sign up
		</a>
	</div>
</form>
