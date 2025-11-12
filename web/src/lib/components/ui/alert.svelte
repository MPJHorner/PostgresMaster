<script lang="ts">
	import { cn } from '$lib/utils';
	import type { Snippet } from 'svelte';
	import type { HTMLAttributes } from 'svelte/elements';

	interface AlertProps extends HTMLAttributes<HTMLDivElement> {
		variant?: 'default' | 'destructive';
		children?: Snippet;
	}

	let { class: className, variant = 'default', children, ...restProps }: AlertProps = $props();

	const variants = {
		default: 'bg-white text-slate-950 dark:bg-slate-950 dark:text-slate-50',
		destructive:
			'border-red-500/50 text-red-500 dark:border-red-500 [&>svg]:text-red-500 dark:border-red-900/50 dark:text-red-900 dark:dark:border-red-900 dark:[&>svg]:text-red-900'
	};
</script>

<div
	class={cn(
		'relative w-full rounded-lg border border-slate-200 p-4 [&>svg~*]:pl-7 [&>svg+div]:translate-y-[-3px] [&>svg]:absolute [&>svg]:left-4 [&>svg]:top-4 [&>svg]:text-slate-950 dark:border-slate-800 dark:[&>svg]:text-slate-50',
		variants[variant],
		className
	)}
	role="alert"
	{...restProps}
>
	{#if children}
		{@render children()}
	{/if}
</div>
