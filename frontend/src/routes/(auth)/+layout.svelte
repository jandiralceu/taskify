<script lang="ts">
	let { children } = $props();
</script>

<div
	class="flex min-h-screen items-center justify-center bg-gray-50 px-4 py-12 sm:px-6 lg:px-8 dark:bg-gray-900"
>
	<div class="w-full max-w-md space-y-8">
		<!-- Aqui você pode colocar elementos comuns de autenticação, como uma Logo -->
		<div class="text-center">
			<h2 class="mt-6 text-3xl font-extrabold tracking-tight text-gray-900 dark:text-white">
				Taskify
			</h2>
		</div>

		<!-- O snippet {@render children()} é OBRIGATÓRIO. 
             É aqui que o conteúdo das páginas (login, cadastro, etc) será injetado. -->
		{@render children()}
	</div>
</div>
