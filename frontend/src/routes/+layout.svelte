<script lang="ts">
	import { QueryClient, QueryClientProvider } from '@tanstack/svelte-query';
	import { Toast } from '@skeletonlabs/skeleton-svelte';
	import { toaster } from '$lib/state/toast.svelte';
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

<Toast.Group {toaster}>
	{#snippet children(toast)}
		<Toast {toast} class="bg-white dark:bg-surface-900 border border-surface-200 dark:border-surface-700 shadow-xl rounded-xl p-4">
			<Toast.Message>
				<Toast.Title class="font-bold text-surface-900 dark:text-surface-50">{toast.title}</Toast.Title>
				<Toast.Description class="text-sm text-surface-600 dark:text-surface-400">{toast.description}</Toast.Description>
			</Toast.Message>
			<Toast.CloseTrigger class="btn-icon btn-icon-sm hover:bg-surface-100 dark:hover:bg-surface-800" />
		</Toast>
	{/snippet}
</Toast.Group>
