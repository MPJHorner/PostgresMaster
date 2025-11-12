<script lang="ts">
	import { cn } from '$lib/utils';
	import type { Snippet } from 'svelte';
	import type { HTMLAttributes } from 'svelte/elements';

	interface TableRowProps extends HTMLAttributes<HTMLTableRowElement> {
		children?: Snippet;
	}

	let { class: className, children, ...restProps }: TableRowProps = $props();
</script>

<tr
	class={cn(
		'border-b transition-colors hover:bg-slate-100/50 data-[state=selected]:bg-slate-100 dark:hover:bg-slate-800/50 dark:data-[state=selected]:bg-slate-800',
		className
	)}
	{...restProps}
>
	{#if children}
		{@render children()}
	{/if}
</tr>
