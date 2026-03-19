<script lang="ts">
  import { QueryClient, QueryClientProvider } from '@tanstack/svelte-query'
  import { Toast } from '@skeletonlabs/skeleton-svelte'
  import { toaster } from '$lib/state/toast.svelte'
  import '@fontsource-variable/roboto'
  import { 
    CircleCheck, 
    CircleAlert, 
    Info, 
    TriangleAlert, 
    X
  } from '@lucide/svelte'
  import '../app.css'

  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        enabled: true,
        retry: 1,
        staleTime: 60 * 1000,
      },
    },
  })

  let { children } = $props()
</script>


<!-- Slide-in from top animation context -->
<style>
  :global(.toast-top-end) {
    top: 2rem !important;
    right: 2rem !important;
    bottom: auto !important;
  }
</style>

<QueryClientProvider client={queryClient}>
  {@render children()}
</QueryClientProvider>

<Toast.Group {toaster}>
  {#snippet children(toast)}
    {@const type = toast.type || 'info'}
    {@const themes = {
      success: {
        border: 'border-emerald-500',
        bg: 'bg-emerald-50/90 dark:bg-emerald-950/20',
        text: 'text-emerald-800 dark:text-emerald-200',
        icon: CircleCheck,
      },
      error: {
        border: 'border-rose-500',
        bg: 'bg-rose-50/90 dark:bg-rose-950/20',
        text: 'text-rose-800 dark:text-rose-200',
        icon: CircleAlert,
      },
      warning: {
        border: 'border-amber-500',
        bg: 'bg-amber-50/90 dark:bg-amber-950/20',
        text: 'text-amber-800 dark:text-amber-200',
        icon: TriangleAlert,
      },
      info: {
        border: 'border-sky-500',
        bg: 'bg-sky-50/90 dark:bg-sky-950/20',
        text: 'text-sky-800 dark:text-sky-200',
        icon: Info,
      },
    }}
    {@const theme = themes[type as keyof typeof themes] || themes.info}
    <Toast
      {toast}
      class="flex min-w-[320px] max-w-md gap-4 overflow-hidden rounded-2xl border-l-[6px] {theme.border} {theme.bg} p-5 shadow-2xl backdrop-blur-md transition-all duration-300 animate-in slide-in-from-top-4"
    >
      <div class="flex size-10 shrink-0 items-center justify-center rounded-xl bg-white/50 shadow-sm dark:bg-white/5">
        <theme.icon size={22} class={theme.text} />
      </div>
      <Toast.Message class="flex-1">
        <Toast.Title class="text-base font-bold text-surface-900 dark:text-surface-50"
          >{toast.title}</Toast.Title
        >
        <Toast.Description
          class="mt-0.5 text-sm leading-relaxed text-surface-600 dark:text-surface-400"
          >{toast.description}</Toast.Description
        >
      </Toast.Message>
      <Toast.CloseTrigger
        class="flex size-8 items-center justify-center rounded-lg text-surface-400 transition-colors hover:bg-white/50 hover:text-surface-900 dark:hover:bg-white/5 dark:hover:text-white"
      >
        <X size={16} />
      </Toast.CloseTrigger>
    </Toast>
  {/snippet}
</Toast.Group>
