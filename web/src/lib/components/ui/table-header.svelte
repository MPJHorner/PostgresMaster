<script lang="ts">
	import { cn } from '$lib/utils';
	import type { Snippet } from 'svelte';
	import type { HTMLAttributes } from 'svelte/elements';

	interface TableHeaderProps extends HTMLAttributes<HTMLTableSectionElement> {
		children?: Snippet;
	}

	let { class: className, children, ...restProps }: TableHeaderProps = $props();
</script>

<thead class={cn('[&_tr]:border-b', className)} {...restProps}>
	{#if children}
		{@render children()}
	{/if}
</thead>
