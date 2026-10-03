export function composeEventHandlers<E extends Event>(
  originalHandler?: ((event: E) => void) | null,
  ourHandler?: ((event: E) => void) | null,
  { checkForDefaultPrevented = true } = {}
) {
  return function handleEvent(event: E) {
    originalHandler?.(event);

    if (checkForDefaultPrevented === false || !event.defaultPrevented) {
      return ourHandler?.(event);
    }
  };
}
