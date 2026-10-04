/**
 * Shared floating-overlay helpers.
 *
 * Overlays (tooltip / dropdown / popover / ...) render inline in the DOM by
 * default, so any ancestor with `overflow: hidden | auto | scroll` (cards,
 * explorer panes, dialogs, scroll areas) clips them. Portaling the content to
 * `document.body` with `position: fixed` escapes all of those ancestors.
 */

/** Svelte action — move the node to `target` (default `document.body`). */
export function portal(node: HTMLElement, target: HTMLElement | string = document.body) {
	const parent =
		typeof target === 'string'
			? (document.querySelector(target) as HTMLElement | null) ?? document.body
			: target;
	parent.appendChild(node);
	return {
		destroy() {
			node.remove();
		}
	};
}

export type FloatingSide = 'top' | 'bottom' | 'left' | 'right';
export type FloatingAlign = 'start' | 'center' | 'end';

export interface FloatingOptions {
	side?: FloatingSide;
	align?: FloatingAlign;
	/** Gap in px between anchor and overlay. */
	offset?: number;
}

/**
 * Position a `position: fixed` overlay relative to an anchor element.
 * Centers/aligns per `side`+`align`, flips when it would overflow the
 * viewport, then clamps with an 8px margin.
 */
export function positionFloating(
	el: HTMLElement,
	anchor: HTMLElement,
	opts: FloatingOptions = {}
): void {
	const { side = 'bottom', align = 'start', offset = 6 } = opts;
	const r = anchor.getBoundingClientRect();
	const tw = el.offsetWidth;
	const th = el.offsetHeight;
	const vw = window.innerWidth;
	const vh = window.innerHeight;
	let top = 0;
	let left = 0;

	if (side === 'top' || side === 'bottom') {
		if (align === 'start') left = r.left;
		else if (align === 'end') left = r.right - tw;
		else left = r.left + r.width / 2 - tw / 2;
		top = side === 'bottom' ? r.bottom + offset : r.top - offset - th;
		if (side === 'bottom' && top + th > vh - 8) top = r.top - offset - th;
		if (side === 'top' && top < 8) top = r.bottom + offset;
	} else {
		if (align === 'start') top = r.top;
		else if (align === 'end') top = r.bottom - th;
		else top = r.top + r.height / 2 - th / 2;
		left = side === 'right' ? r.right + offset : r.left - offset - tw;
		if (side === 'right' && left + tw > vw - 8) left = r.left - offset - tw;
		if (side === 'left' && left < 8) left = r.right + offset;
	}

	left = Math.min(Math.max(8, left), Math.max(8, vw - tw - 8));
	top = Math.min(Math.max(8, top), Math.max(8, vh - th - 8));
	el.style.top = `${Math.round(top)}px`;
	el.style.left = `${Math.round(left)}px`;
}
