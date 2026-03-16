<script lang="ts">
	import {
		ListTodo,
		Users,
		LogOut
	} from '@lucide/svelte';
	import { page } from '$app/state';
	import { storage, AUTH_KEYS } from '$lib/utils/storage';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import logoWhite from '$lib/assets/logo_white.webp';
	import { authService } from '$lib/api/auth.service';
	import { createProfileQuery, PROFILE_QUERY_KEY, authState } from '$lib/state/user.svelte';
	import { useQueryClient } from '@tanstack/svelte-query';

	let { children } = $props();
	const queryClient = useQueryClient();
	const profileQuery = createProfileQuery();

	async function handleLogout() {
		const refreshToken = storage.get(AUTH_KEYS.REFRESH_TOKEN);
		if (refreshToken) {
			try {
				await authService.signout({ refreshToken });
			} catch (error) {
				console.error('Failed to sign out on server:', error);
			}
		}
		storage.remove(AUTH_KEYS.ACCESS_TOKEN);
		storage.remove(AUTH_KEYS.REFRESH_TOKEN);
		
		// Limpa o estado reativo e o cache
		authState.token = null;
		queryClient.setQueryData(PROFILE_QUERY_KEY, null);
		queryClient.removeQueries({ queryKey: PROFILE_QUERY_KEY });
		
		goto(resolve('/signin'));
	}

	const navItems = [
		{ icon: ListTodo, label: 'Tasks', href: '/' },
		{ icon: Users, label: 'Users', href: '/users' },
	] as const;

	let activeRoute = $derived(page.url.pathname);

	let userInitials = $derived.by(() => {
		if (!profileQuery.data) return '';
		const f = profileQuery.data.firstName?.[0] || '';
		const l = profileQuery.data.lastName?.[0] || '';
		return (f + l).toUpperCase();
	});
</script>

<div class="h-screen w-full bg-[#F7F3F9] flex overflow-hidden font-sans text-surface-900">
	<!-- Left Sidebar (Nubank Purple) -->
	<aside class="w-[75px] bg-primary-500 flex flex-col items-center py-6 z-30 shrink-0 shadow-2xl">
		<a href={resolve('/')} class="mb-10 block hover:scale-110 transition-transform">
			<img src={logoWhite} alt="Taskify" class="w-8 h-auto" />
		</a>

		<nav class="flex-1 flex flex-col gap-6 w-full items-center mt-12">
			{#each navItems as item (item.href)}
				<a 
					href={resolve(item.href)}
					class="p-2.5 rounded-xl transition-all duration-300 {activeRoute === resolve(item.href) ? 'text-primary-500 bg-white shadow-lg' : 'text-white/40 hover:text-white/70'}"
					title={item.label}
				>
					<item.icon size={22} strokeWidth={2.5} />
				</a>
			{/each}

			<div class="w-8 h-px bg-white/10 my-2"></div>

			<!-- Profile Avatar -->
			<a 
				href={resolve('/profile')}
				class="w-10 h-10 rounded-xl overflow-hidden border-2 transition-all duration-300 flex items-center justify-center text-xs font-bold {activeRoute === resolve('/profile') ? 'border-white shadow-lg scale-110 bg-white text-primary-500' : 'border-transparent bg-white/10 text-white hover:bg-white/20 hover:border-white/20'}"
				title={profileQuery.data ? `${profileQuery.data.firstName} ${profileQuery.data.lastName}` : 'Meu Perfil'}
			>
				{#if profileQuery.data?.avatarUrl}
					<img 
						src={profileQuery.data.avatarUrl} 
						alt="Profile" 
						class="w-full h-full object-cover"
					/>
				{:else if userInitials}
					<span>{userInitials}</span>
				{:else}
					<img 
						src="https://api.dicebear.com/7.x/avataaars/svg?seed=user" 
						alt="Profile" 
						class="w-full h-full object-cover opacity-50"
					/>
				{/if}
			</a>
		</nav>

		<button 
			onclick={handleLogout}
			class="mt-auto p-4 text-white/40 hover:text-rose-400 transition-colors"
			title="Logout"
		>
			<LogOut size={22} />
		</button>
	</aside>

	<!-- Main Body Wrapper -->
	<div class="flex-1 flex flex-col min-w-0 relative">
		
		<!-- Page Content Area -->
		<main class="flex-1 overflow-x-auto overflow-y-auto custom-scrollbar">
			<div class="h-full min-w-max">
				{@render children()}
			</div>
		</main>
	</div>
</div>

<style>
	:global(body) {
		background-color: #F7F3F9; /* Um tom lavanda bem sutil para alinhar com o roxo da marca */
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
