<script lang="ts">
	import { createMutation } from '@tanstack/svelte-query';
	import { z } from 'zod';
	import { Mail, Lock } from '@lucide/svelte';

	import { resolve } from '$app/paths';
	import { goto } from '$app/navigation';
	import Input from '$lib/components/Input.svelte';
	import logo from '$lib/assets/logo.webp';
	import { authService } from '$lib/api/auth.service';
	import { storage, AUTH_KEYS } from '$lib/utils/storage';
	import { PROFILE_QUERY_KEY } from '$lib/state/user.svelte';
	import { useQueryClient } from '@tanstack/svelte-query';
	import type { SignInRequest, SignInResponse } from '$lib/api/types';

	let email = $state('');
	let password = $state('');
	let submitted = $state(false);
	const queryClient = useQueryClient();

	const signinSchema = z.object({
		email: z.email('Please enter a valid email address'),
		password: z.string().min(1, 'Password is required')
	});

	let validationResult = $derived(signinSchema.safeParse({ email, password }));
	let fieldErrors = $derived(submitted && !validationResult.success ? z.treeifyError(validationResult.error).properties : undefined);

	const signinMutation = createMutation(() => ({
		mutationFn: (data: SignInRequest) => authService.signin(data),
		onSuccess: async (data: SignInResponse) => {
			storage.set(AUTH_KEYS.ACCESS_TOKEN, data.accessToken);
			storage.set(AUTH_KEYS.REFRESH_TOKEN, data.refreshToken);
			
			import('$lib/state/user.svelte').then(m => m.authState.token = data.accessToken);
			await queryClient.invalidateQueries({ queryKey: PROFILE_QUERY_KEY });
			goto(resolve('/'));
		},
		onError: (error: Error) => {
			console.error('Login failed:', error);
		}
	}));

	function handleSubmit(e: SubmitEvent) {
		e.preventDefault();
		submitted = true;

		if (!validationResult.success) {
			return;
		}

		signinMutation.mutate({ email, password });
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
			disabled={signinMutation.isPending}
			class="group relative flex w-full justify-center rounded-xl border border-transparent bg-primary-500 px-4 py-3 text-sm font-semibold text-white transition-all hover:bg-primary-600 focus:ring-2 focus:ring-primary-500/50 focus:outline-none active:scale-[0.98] disabled:opacity-50 disabled:cursor-not-allowed"
		>
			{signinMutation.isPending ? 'Logging in...' : 'Sign in'}
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
