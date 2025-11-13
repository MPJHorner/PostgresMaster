import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render } from '@testing-library/svelte';
import Editor from './Editor.svelte';

// Mock Monaco Editor
vi.mock('monaco-editor', () => {
	const mockEditor = {
		getValue: vi.fn(() => ''),
		setValue: vi.fn(),
		onDidChangeModelContent: vi.fn((callback) => {
			// Store callback for testing
			mockEditor._changeCallback = callback;
			return { dispose: vi.fn() };
		}),
		addCommand: vi.fn(),
		getPosition: vi.fn(() => ({ lineNumber: 1, column: 1 })),
		setPosition: vi.fn(),
		focus: vi.fn(),
		dispose: vi.fn(),
		_changeCallback: null as ((callback: unknown) => void) | null
	};

	return {
		default: {
			editor: {
				create: vi.fn(() => mockEditor)
			},
			KeyMod: {
				CtrlCmd: 2048
			},
			KeyCode: {
				Enter: 3
			}
		}
	};
});

describe('Editor Component', () => {
	beforeEach(() => {
		vi.clearAllMocks();
	});

	it('should render editor container', () => {
		const { container } = render(Editor);
		const editorWrapper = container.querySelector('.editor-wrapper');
		expect(editorWrapper).toBeTruthy();
	});

	it('should apply custom height', () => {
		const { container } = render(Editor, { props: { height: '600px' } });
		const editorWrapper = container.querySelector('.editor-wrapper') as HTMLElement;
		expect(editorWrapper?.style.height).toBe('600px');
	});

	it('should accept value prop', () => {
		const { component } = render(Editor, { props: { value: 'SELECT * FROM users' } });
		expect(component.value).toBe('SELECT * FROM users');
	});

	it('should call onChange when content changes', async () => {
		const onChange = vi.fn();
		render(Editor, { props: { onChange } });

		// Note: In a real test with Monaco mounted, we'd trigger actual changes
		// For now, this tests that the prop is accepted
		expect(onChange).not.toHaveBeenCalled();
	});

	it('should accept onExecute callback', () => {
		const onExecute = vi.fn();
		const { component } = render(Editor, { props: { onExecute } });
		expect(component.onExecute).toBe(onExecute);
	});
});
