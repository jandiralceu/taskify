<script lang="ts">
	import type { Snippet } from 'svelte';
	import { Eye, EyeOff } from '@lucide/svelte';

	interface Props {
		value: string;
		type?: 'text' | 'email' | 'password';
		id: string;
		name?: string;
		placeholder?: string;
		label?: string;
		required?: boolean;
		icon?: Snippet;
		error?: string;
	}

	let {
		value = $bindable(),
		type = 'text',
		id,
		name,
		placeholder,
		label,
		required = false,
		icon,
		error
	}: Props = $props();

	let showPassword = $state(false);
	let inputType = $derived(type === 'password' && showPassword ? 'text' : type);
</script>

<div class="w-full space-y-1">
	{#if label}
		<label for={id} class="block text-sm font-medium text-surface-700 dark:text-surface-300">
			{label}
		</label>
	{/if}
	<div class="relative">
		{#if icon}
			<div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3 text-surface-500">
				{@render icon()}
			</div>
		{/if}

		<input
			{id}
			{name}
			{required}
			{placeholder}
			type={inputType}
			bind:value={value}
			class="block h-12 w-full rounded-xl border bg-surface-50 px-3.5 text-surface-900 placeholder-surface-600 transition-all focus:ring-2 focus:outline-none sm:text-sm dark:bg-surface-900 dark:text-surface-50 dark:placeholder-surface-400 {icon ? 'pl-10' : 'pl-4'} {type === 'password' ? 'pr-10' : 'pr-4'} {error ? 'border-red-500 focus:border-red-500 focus:ring-red-500/20 dark:border-red-500' : 'border-surface-300 focus:border-primary-500 focus:ring-primary-500/20 dark:border-surface-700'}"
		/>

		{#if type === 'password'}
			<button
				type="button"
				class="absolute inset-y-0 right-0 flex items-center pr-3 text-surface-400 transition-colors hover:text-primary-500 dark:hover:text-primary-400"
				onclick={() => (showPassword = !showPassword)}
				tabindex="-1"
			>
				{#if showPassword}
					<Eye size={20} />
				{:else}
					<EyeOff size={20} />
				{/if}
			</button>
		{/if}
	</div>
	{#if error}
		<p class="mt-1 text-xs text-red-500 font-medium">{error}</p>
	{/if}
</div>
