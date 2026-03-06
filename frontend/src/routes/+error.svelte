<script lang="ts">
	import { page } from '$app/state';
</script>

<div class="flex min-h-screen flex-col items-center justify-center space-y-4 px-4 text-center">
	<div class="space-y-2">
		<h1 class="text-9xl font-black text-gray-200">{page.status}</h1>
		<p class="text-2xl font-semibold text-gray-800 tracking-tight">Oops! Algo deu errado.</p>
		<p class="max-w-md text-gray-500 text-lg leading-relaxed">
			{page.error?.message || 'Ocorreu um erro inesperado. Por favor, tente novamente mais tarde.'}
		</p>
	</div>

	<div class="pt-6">
		<a
			href="/"
			class="inline-flex h-12 items-center justify-center rounded-lg bg-indigo-600 px-8 font-semibold text-white transition-all hover:bg-indigo-700 hover:scale-105 active:scale-95 shadow-lg shadow-indigo-200 focus:outline-none focus:ring-4 focus:ring-indigo-300"
		>
			Voltar para o Início
		</a>
	</div>
</div>
