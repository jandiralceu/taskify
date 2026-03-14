<script lang="ts">
	import {
		ListTodo,
		Users,
		Search as SearchIcon,
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
</script>

<div class="h-screen w-full bg-[#F7F3F9] flex overflow-hidden font-sans text-surface-900">
	<!-- Left Sidebar (Nubank Purple) -->
	<aside class="w-[75px] bg-[#2E1065] flex flex-col items-center py-6 z-30 shrink-0 shadow-2xl">
		<a href={resolve('/')} class="mb-10 block hover:scale-110 transition-transform">
			<img src={logoWhite} alt="Taskify" class="w-8 h-auto" />
		</a>

		<nav class="flex-1 flex flex-col gap-6 w-full items-center mt-12">
			{#each navItems as item (item.href)}
				<a 
					href={resolve(item.href)}
					class="p-2.5 rounded-xl transition-all duration-300 {activeRoute === resolve(item.href) ? 'text-[#2E1065] bg-white shadow-lg' : 'text-white/40 hover:text-white/70'}"
					title={item.label}
				>
					<item.icon size={22} strokeWidth={2.5} />
				</a>
			{/each}

			<div class="w-8 h-px bg-white/10 my-2"></div>

			<!-- Profile Avatar -->
			<a 
				href={resolve('/profile')}
				class="w-10 h-10 rounded-xl overflow-hidden border-2 transition-all duration-300 {activeRoute === resolve('/profile') ? 'border-white shadow-lg scale-110' : 'border-transparent opacity-50 hover:opacity-100 hover:border-white/20'}"
				title={profileQuery.data ? `${profileQuery.data.firstName} ${profileQuery.data.lastName}` : 'Meu Perfil'}
			>
				<img 
					src="https://api.dicebear.com/7.x/avataaars/svg?seed={profileQuery.data?.firstName || 'user'}" 
					alt="Profile" 
					class="w-full h-full object-cover"
				/>
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
		
		<!-- Top Navigation Header -->
		<header class="h-20 flex items-center justify-between px-8 shrink-0 bg-transparent">
			<!-- Search Bar -->
			<div class="w-full max-w-[400px] relative">
				<div class="absolute inset-y-0 left-4 flex items-center text-surface-400">
					<SearchIcon size={18} />
				</div>
				<input 
					type="text" 
					placeholder="Search" 
					class="w-full h-11 bg-white border-none rounded-full pl-12 pr-6 text-sm shadow-sm focus:ring-2 focus:ring-[#820AD1]/10 placeholder:text-surface-300"
				/>
			</div>

			<!-- Right Side Placeholder -->
			<div class="w-[180px]"></div>
		</header>

		<!-- Page Content Area -->
		<main class="flex-1 overflow-x-auto overflow-y-hidden custom-scrollbar">
			<div class="h-full min-w-max p-8 pt-0">
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
