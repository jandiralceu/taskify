<script lang="ts">
	import authBg from '$lib/assets/auth_bg.jpg';
	let { children } = $props();
</script>

<div class="grid grid-cols-[600px_1fr] min-h-screen">
	<main class="flex items-center justify-center">
		{@render children()}
	</main>

	<div 
		class="bg-cover bg-center bg-no-repeat"
		style="background-image: url({authBg})"
	>
	</div>
</div>
