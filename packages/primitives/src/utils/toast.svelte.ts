export interface ToastItem {
  id: string;
  type: 'default' | 'success' | 'warning' | 'error' | 'info';
  title?: string;
  description?: string;
  duration?: number;
}

class ToastManager {
  toasts = $state<ToastItem[]>([]);

  show(item: Omit<ToastItem, 'id'>) {
    const id = 'toast-' + Math.random().toString(36).substring(2, 9);
    const toastItem: ToastItem = { ...item, id };
    this.toasts.push(toastItem);

    const duration = item.duration ?? 4000;
    if (duration > 0) {
      setTimeout(() => {
        this.dismiss(id);
      }, duration);
    }

    return id;
  }

  /** Alias matching the spec usage: `toast.add({ title, description })`. */
  add(item: Omit<ToastItem, 'id'>) {
    return this.show(item);
  }

  success(title: string, description?: string) {    return this.show({ type: 'success', title, description });
  }

  error(title: string, description?: string) {
    return this.show({ type: 'error', title, description });
  }

  warning(title: string, description?: string) {
    return this.show({ type: 'warning', title, description });
  }

  info(title: string, description?: string) {
    return this.show({ type: 'info', title, description });
  }

  promise<T>(
    task: Promise<T>,
    messages: { loading: string; success: string | ((value: T) => string); error: string | ((err: unknown) => string) }
  ): Promise<T> {
    const id = this.show({ type: 'info', title: messages.loading, duration: 0 });
    task.then(
      (value) => {
        this.dismiss(id);
        const title = typeof messages.success === 'function' ? messages.success(value) : messages.success;
        this.success(title);
      },
      (err: unknown) => {
        this.dismiss(id);
        const title = typeof messages.error === 'function' ? messages.error(err) : messages.error;
        this.error(title);
      }
    );
    return task;
  }

  dismiss(id: string) {
    this.toasts = this.toasts.filter((t: ToastItem) => t.id !== id);
  }
}

export const toast = new ToastManager();
