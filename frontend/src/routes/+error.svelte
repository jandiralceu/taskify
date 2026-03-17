<script lang="ts">
  import { page } from '$app/state'
  import { resolve } from '$app/paths'
</script>

<div
  class="flex min-h-screen flex-col items-center justify-center space-y-4 px-4 text-center"
>
  <div class="space-y-2">
    <h1 class="text-9xl font-black text-gray-200">{page.status}</h1>
    <p class="text-2xl font-semibold tracking-tight text-gray-800">
      Oops! Algo deu errado.
    </p>
    <p class="max-w-md text-lg leading-relaxed text-gray-500">
      {page.error?.message ||
        'Ocorreu um erro inesperado. Por favor, tente novamente mais tarde.'}
    </p>
  </div>

  <div class="pt-6">
    <a
      href={resolve('/')}
      class="inline-flex h-12 items-center justify-center rounded-lg bg-indigo-600 px-8 font-semibold text-white shadow-lg shadow-indigo-200 transition-all hover:scale-105 hover:bg-indigo-700 focus:ring-4 focus:ring-indigo-300 focus:outline-none active:scale-95"
    >
      Voltar para o Início
    </a>
  </div>
</div>
