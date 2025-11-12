<script lang="ts">
	import { cn } from '$lib/utils';
	import type { HTMLAttributes } from 'svelte/elements';

	interface SeparatorProps extends HTMLAttributes<HTMLDivElement> {
		orientation?: 'horizontal' | 'vertical';
		decorative?: boolean;
	}

	let {
		class: className,
		orientation = 'horizontal',
		decorative = false,
		...restProps
	}: SeparatorProps = $props();
</script>

<div
	class={cn(
		'shrink-0 bg-slate-200 dark:bg-slate-800',
		orientation === 'horizontal' ? 'h-[1px] w-full' : 'h-full w-[1px]',
		className
	)}
	role={decorative ? 'none' : 'separator'}
	aria-orientation={orientation}
	{...restProps}
></div>
