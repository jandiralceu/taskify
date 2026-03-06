<script lang="ts">
	import { resolve } from '$app/paths';
	let email = $state('');
	let password = $state('');

	function handleSubmit(e: SubmitEvent) {
		e.preventDefault();
		console.log('Login attempt:', { email, password });
		// Implementar lógica de login aqui
	}
</script>

<form class="mt-8 space-y-6" onsubmit={handleSubmit}>
	<div class="-space-y-px rounded-md shadow-sm">
		<div>
			<label for="email-address" class="sr-only">E-mail</label>
			<input
				id="email-address"
				name="email"
				type="email"
				required
				bind:value={email}
				class="relative block w-full appearance-none rounded-none rounded-t-md border border-gray-300 px-3 py-2 text-gray-900 placeholder-gray-500 focus:z-10 focus:border-indigo-500 focus:ring-indigo-500 focus:outline-none sm:text-sm dark:border-gray-700 dark:bg-gray-800 dark:text-white"
				placeholder="Endereço de e-mail"
			/>
		</div>
		<div>
			<label for="password" class="sr-only">Senha</label>
			<input
				id="password"
				name="password"
				type="password"
				required
				bind:value={password}
				class="relative block w-full appearance-none rounded-none rounded-b-md border border-gray-300 px-3 py-2 text-gray-900 placeholder-gray-500 focus:z-10 focus:border-indigo-500 focus:ring-indigo-500 focus:outline-none sm:text-sm dark:border-gray-700 dark:bg-gray-800 dark:text-white"
				placeholder="Senha"
			/>
		</div>
	</div>

	<div class="flex items-center justify-between">
		<div class="text-sm">
			<a
				href={resolve('/signup')}
				class="font-medium text-indigo-600 hover:text-indigo-500 dark:text-indigo-400"
			>
				Não tem uma conta? Cadastre-se
			</a>
		</div>
	</div>

	<div>
		<button
			type="submit"
			class="group relative flex w-full justify-center rounded-md border border-transparent bg-indigo-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-indigo-700 focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 focus:outline-none"
		>
			Entrar
		</button>
	</div>
</form>
