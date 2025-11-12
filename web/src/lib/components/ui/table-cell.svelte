<script lang="ts">
	import { cn } from '$lib/utils';
	import type { Snippet } from 'svelte';
	import type { HTMLAttributes } from 'svelte/elements';

	interface TableCellProps extends HTMLAttributes<HTMLTableCellElement> {
		children?: Snippet;
	}

	let { class: className, children, ...restProps }: TableCellProps = $props();
</script>

<td
	class={cn('p-4 align-middle [&:has([role=checkbox])]:pr-0', className)}
	{...restProps}
>
	{#if children}
		{@render children()}
	{/if}
</td>
