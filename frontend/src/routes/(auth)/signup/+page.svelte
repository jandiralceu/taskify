<script lang="ts">
	let name = $state('');
	let email = $state('');
	let password = $state('');
	let confirmPassword = $state('');

	function handleSubmit(e: SubmitEvent) {
		e.preventDefault();
		if (password !== confirmPassword) {
			alert('As senhas não coincidem!');
			return;
		}
		console.log('Signup attempt:', { name, email, password });
		// Implementar lógica de cadastro aqui
	}
</script>

<form class="mt-8 space-y-6" onsubmit={handleSubmit}>
	<div class="-space-y-px rounded-md shadow-sm">
		<div>
			<label for="name" class="sr-only">Nome Completo</label>
			<input
				id="name"
				name="name"
				type="text"
				required
				bind:value={name}
				class="relative block w-full appearance-none rounded-none rounded-t-md border border-gray-300 px-3 py-2 text-gray-900 placeholder-gray-500 focus:z-10 focus:border-indigo-500 focus:ring-indigo-500 focus:outline-none sm:text-sm dark:border-gray-700 dark:bg-gray-800 dark:text-white"
				placeholder="Nome completo"
			/>
		</div>
		<div>
			<label for="email-address" class="sr-only">E-mail</label>
			<input
				id="email-address"
				name="email"
				type="email"
				required
				bind:value={email}
				class="relative block w-full appearance-none rounded-none border border-gray-300 px-3 py-2 text-gray-900 placeholder-gray-500 focus:z-10 focus:border-indigo-500 focus:ring-indigo-500 focus:outline-none sm:text-sm dark:border-gray-700 dark:bg-gray-800 dark:text-white"
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
				class="relative block w-full appearance-none rounded-none border border-gray-300 px-3 py-2 text-gray-900 placeholder-gray-500 focus:z-10 focus:border-indigo-500 focus:ring-indigo-500 focus:outline-none sm:text-sm dark:border-gray-700 dark:bg-gray-800 dark:text-white"
				placeholder="Senha"
			/>
		</div>
		<div>
			<label for="confirm-password" class="sr-only">Confirmar Senha</label>
			<input
				id="confirm-password"
				name="confirm-password"
				type="password"
				required
				bind:value={confirmPassword}
				class="relative block w-full appearance-none rounded-none rounded-b-md border border-gray-300 px-3 py-2 text-gray-900 placeholder-gray-500 focus:z-10 focus:border-indigo-500 focus:ring-indigo-500 focus:outline-none sm:text-sm dark:border-gray-700 dark:bg-gray-800 dark:text-white"
				placeholder="Confirmar senha"
			/>
		</div>
	</div>

	<div class="flex items-center justify-between">
		<div class="text-sm">
			<a
				href="/signin"
				class="font-medium text-indigo-600 hover:text-indigo-500 dark:text-indigo-400"
			>
				Já tem uma conta? Entre aqui
			</a>
		</div>
	</div>

	<div>
		<button
			type="submit"
			class="group relative flex w-full justify-center rounded-md border border-transparent bg-indigo-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-indigo-700 focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 focus:outline-none"
		>
			Criar Conta
		</button>
	</div>
</form>
