<script lang="ts">
	import { QueryClient, QueryClientProvider } from '@tanstack/svelte-query';
	import '../app.css';
	import '@fontsource-variable/roboto';
	import favicon from '$lib/assets/favicon.ico';

	const queryClient = new QueryClient({
		defaultOptions: {
			queries: {
				enabled: true,
				retry: 1,
				staleTime: 60 * 1000
			}
		}
	});

	let { children } = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<QueryClientProvider client={queryClient}>
	{@render children()}
</QueryClientProvider>
