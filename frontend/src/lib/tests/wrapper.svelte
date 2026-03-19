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
    {#if children.name === 'children' || children.length === 0}
      {@render children()}
    {:else}
      {@const Component = children}
      <Component />
    {/if}
  {:else if children}
    {@const Component = children}
    <Component />
  {/if}
</QueryClientProvider>
