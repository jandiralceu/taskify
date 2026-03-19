<script lang="ts">
  import { QueryClient, QueryClientProvider } from '@tanstack/svelte-query'
  let { children } = $props()
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
        staleTime: Infinity,
      },
    },
  })
</script>

<QueryClientProvider client={queryClient}>
  {#if typeof children === 'function'}
    {@render children()}
  {:else}
    <svelte:component this={children} />
  {/if}
</QueryClientProvider>
