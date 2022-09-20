<script lang="ts">
	import { goto } from '$app/navigation';

	function changeRoute(href: string) {
		goto('/' + href);
	}
</script>

<header>
	<div class="flex flex-wrap w-full justify-between">
		<h1 class="text-4xl font-bold">YACOID</h1>
		<div class="flex gap-2">
			<button class="btn btn-disabled">Submissions</button>
			<button class="btn btn-secondary" on:click={() => changeRoute('definitions')}
				>All Definitions</button
			>
			<button class="btn btn-secondary" on:click={() => changeRoute('about_us')}>About us</button>
			<button class="btn btn-secondary" on:click={() => changeRoute('login')}>Login</button>
		</div>
	</div>
</header>

<style lang="postcss">
	header {
		@apply p-4 w-full gap-2 bg-base-200;
	}
</style>
