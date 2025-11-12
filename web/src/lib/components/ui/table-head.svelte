<script lang="ts">
	import { cn } from '$lib/utils';
	import type { Snippet } from 'svelte';
	import type { HTMLAttributes } from 'svelte/elements';

	interface TableHeadProps extends HTMLAttributes<HTMLTableCellElement> {
		children?: Snippet;
	}

	let { class: className, children, ...restProps }: TableHeadProps = $props();
</script>

<th
	class={cn(
		'h-12 px-4 text-left align-middle font-medium text-slate-500 [&:has([role=checkbox])]:pr-0 dark:text-slate-400',
		className
	)}
	{...restProps}
>
	{#if children}
		{@render children()}
	{/if}
</th>
