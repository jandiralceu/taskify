<script lang="ts">
  import type { Snippet } from 'svelte'
  import { LoaderCircle } from '@lucide/svelte'

  interface Props {
    variant?: 'primary' | 'danger' | 'ghost'
    type?: 'button' | 'submit' | 'reset'
    disabled?: boolean
    loading?: boolean
    loadingText?: string
    onclick?: (e: MouseEvent) => void
    class?: string
    children: Snippet
  }

  let {
    variant = 'primary',
    type = 'button',
    disabled = false,
    loading = false,
    loadingText = 'Loading...',
    onclick,
    class: className = '',
    children,
  }: Props = $props()

  const base =
    'inline-flex items-center gap-2 rounded-xl text-sm font-medium transition-all active:scale-95 disabled:cursor-not-allowed disabled:opacity-50'

  const variants = {
    primary: 'bg-primary-500 hover:bg-primary-700 text-white px-8 py-2.5',
    danger: 'bg-rose-500 hover:bg-rose-600 text-white px-5 py-2.5',
    ghost:
      'font-medium text-surface-500 hover:text-surface-900 hover:bg-surface-50 px-5 py-2.5',
  }
</script>

<button
  {type}
  disabled={disabled || loading}
  {onclick}
  class="{base} {variants[variant]} {className}"
>
  {#if loading}
    <LoaderCircle size={15} class="animate-spin" />
    {loadingText}
  {:else}
    {@render children()}
  {/if}
</button>
