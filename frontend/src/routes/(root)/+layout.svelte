<script lang="ts">
  import { useQueryClient } from '@tanstack/svelte-query'
  import { ListTodo, Users, LogOut } from '@lucide/svelte'
  import { page } from '$app/state'
  import { storage, AUTH_KEYS } from '$lib/utils/storage'
  import { goto } from '$app/navigation'
  import { resolveAvatarUrl } from '$lib/utils/avatar'
  import { resolve } from '$app/paths'
  import logoWhite from '$lib/assets/logo_white.webp'
  import { authService } from '$lib/api/auth.service'
  import {
    createProfileQuery,
    createPermissionsQuery,
    PROFILE_QUERY_KEY,
    PERMISSIONS_QUERY_KEY,
    authState,
  } from '$lib/state/user.svelte'

  let { children } = $props()
  const queryClient = useQueryClient()
  const profileQuery = createProfileQuery()
  const permissionsQuery = createPermissionsQuery()

  async function handleLogout() {
    const refreshToken = storage.get(AUTH_KEYS.REFRESH_TOKEN)
    if (refreshToken) {
      try {
        await authService.signout({ refreshToken })
      } catch (error) {
        console.error('Failed to sign out on server:', error)
      }
    }
    storage.remove(AUTH_KEYS.ACCESS_TOKEN)
    storage.remove(AUTH_KEYS.REFRESH_TOKEN)

    // Limpa o estado reativo e o cache
    authState.token = null
    queryClient.setQueryData(PROFILE_QUERY_KEY, null)
    queryClient.removeQueries({ queryKey: PROFILE_QUERY_KEY })
    queryClient.setQueryData(PERMISSIONS_QUERY_KEY, null)
    queryClient.removeQueries({ queryKey: PERMISSIONS_QUERY_KEY })

    goto(resolve('/signin'))
  }

  const baseNavItems = [
    { icon: ListTodo, label: 'Tasks', href: '/', adminOnly: false },
    { icon: Users, label: 'Users', href: '/users', adminOnly: true },
  ] as const

  let navItems = $derived.by(() => {
    const isAdmin = permissionsQuery.data?.permissions.admin_area ?? false
    return baseNavItems.filter(item => !item.adminOnly || isAdmin)
  })

  let activeRoute = $derived(page.url.pathname)

  let userInitials = $derived.by(() => {
    if (!profileQuery.data) return ''
    const f = profileQuery.data.firstName?.[0] || ''
    const l = profileQuery.data.lastName?.[0] || ''
    return (f + l).toUpperCase()
  })
</script>

<div
  class="flex h-screen w-full overflow-hidden bg-[#F7F3F9] font-sans text-surface-900"
>
  <!-- Left Sidebar (Nubank Purple) -->
  <aside
    class="z-30 flex w-[75px] shrink-0 flex-col items-center bg-primary-500 py-6 shadow-2xl"
  >
    <a
      href={resolve('/')}
      class="mb-10 block transition-transform hover:scale-110"
    >
      <img src={logoWhite} alt="Taskify" class="h-auto w-8" />
    </a>

    <nav class="mt-12 flex w-full flex-1 flex-col items-center gap-6">
      {#each navItems as item (item.href)}
        <a
          href={resolve(item.href)}
          class="rounded-xl p-2.5 transition-all duration-300 {activeRoute ===
          resolve(item.href)
            ? 'bg-white text-primary-600 shadow-lg'
            : 'text-white/70 hover:text-white/90'}"
          title={item.label}
        >
          <item.icon size={22} strokeWidth={2.5} />
        </a>
      {/each}

      <div class="my-2 h-px w-8 bg-white/10"></div>

      <!-- Profile Avatar -->
      <a
        href={resolve('/profile')}
        class="flex h-10 w-10 items-center justify-center overflow-hidden rounded-xl border-2 text-xs font-bold transition-all duration-300 {activeRoute ===
        resolve('/profile')
          ? 'scale-110 border-white bg-white text-primary-500 shadow-lg'
          : 'border-transparent bg-white/10 text-white hover:border-white/20 hover:bg-white/20'}"
        title={profileQuery.data
          ? `${profileQuery.data.firstName} ${profileQuery.data.lastName}`
          : 'Meu Perfil'}
      >
        {#if profileQuery.data?.avatarUrl}
          <img
            src={resolveAvatarUrl(profileQuery.data.avatarUrl) ?? ''}
            alt="Profile"
            class="h-full w-full object-cover"
          />
        {:else if userInitials}
          <span>{userInitials}</span>
        {:else}
          <img
            src="https://api.dicebear.com/7.x/avataaars/svg?seed=user"
            alt="Profile"
            class="h-full w-full object-cover opacity-50"
          />
        {/if}
      </a>
    </nav>

    <button
      onclick={handleLogout}
      class="mt-auto p-4 text-white/70 transition-colors hover:text-white/90"
      title="Logout"
    >
      <LogOut size={22} />
    </button>
  </aside>

  <!-- Main Body Wrapper -->
  <div class="relative flex min-w-0 flex-1 flex-col">
    <!-- Page Content Area -->
    <main class="custom-scrollbar flex-1 overflow-x-auto overflow-y-auto">
      <div class="h-full min-w-max">
        {@render children()}
      </div>
    </main>
  </div>
</div>

<style>
  :global(body) {
    background-color: #f7f3f9;
    margin: 0;
  }

  .custom-scrollbar::-webkit-scrollbar {
    height: 8px;
    width: 8px;
  }
  .custom-scrollbar::-webkit-scrollbar-track {
    background: transparent;
  }
  .custom-scrollbar::-webkit-scrollbar-thumb {
    background: #0f5a5a1a;
    border-radius: 20px;
  }
</style>
