<script lang="ts">
	import logo from '$lib/assets/logo.webp';
	import { resolve } from '$app/paths';
	import { Mail, Lock } from '@lucide/svelte';
	import Input from '$lib/components/Input.svelte';

	let email = $state('');
	let password = $state('');

	function handleSubmit(e: SubmitEvent) {
		e.preventDefault();
		console.log('Login attempt:', { email, password });
	}
</script>

<form class="mt-8 space-y-6 w-full max-w-md" onsubmit={handleSubmit}>
	<h1>
		<img src={logo} alt="Taskify Logo" class="h-12 w-auto" />
	</h1>

	<h2 class="text-primary-500 text-3xl font-bold">Sign in to your account</h2>
	<p class="text-sm text-gray-600">Welcome back! Please sign in to your account.</p>

	<div class="space-y-4">
		<Input
			id="email-address"
			name="email"
			type="email"
			placeholder="Email"
			bind:value={email}
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
			Entrar
		</button>
	</div>

	<div class="flex items-center justify-center">
		<div class="text-sm">
			<a
				href={resolve('/signup')}
				class="font-medium text-primary-600 hover:text-primary-500 dark:text-primary-400"
			>
				Não tem uma conta? Cadastre-se
			</a>
		</div>
	</div>
</form>
