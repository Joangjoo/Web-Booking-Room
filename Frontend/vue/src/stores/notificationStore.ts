import { defineStore } from "pinia";
import { ref } from "vue";

export interface Toast {
  id: number;
  type: "success" | "error" | "warning" | "info";
  title: string;
  message: string;
  duration?: number;
}

export const useNotificationStore = defineStore("notification", () => {
  const toasts = ref<Toast[]>([]);
  const timers = new Map<number, number>();

  const showToast = (toast: Omit<Toast, "id">) => {
    const id = Date.now();
    const duration = toast.duration || 5000;

    toasts.value.push({ id, ...toast, duration });

    const timer = window.setTimeout(() => {
      removeToast(id);
    }, duration);
    timers.set(id, timer);
  };

  const removeToast = (id: number) => {
    const index = toasts.value.findIndex((toast) => toast.id === id);
    if (index !== -1) {
      toasts.value.splice(index, 1);
      if (timers.has(id)) {
        clearTimeout(timers.get(id));
        timers.delete(id);
      }
    }
  };

  return { toasts, showToast, removeToast };
});
