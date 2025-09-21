<script setup lang="ts">
import { computed } from "vue";
import { useNotificationStore, type Toast } from "../stores/notificationStore";
import {
  CheckCircleIcon,
  XCircleIcon,
  ExclamationTriangleIcon,
  InformationCircleIcon,
  XMarkIcon,
} from "@heroicons/vue/24/outline"; 

const notificationStore = useNotificationStore();
const toasts = computed(() => notificationStore.toasts);

const getIcon = (type: Toast["type"]) => {
  switch (type) {
    case "success":
      return CheckCircleIcon;
    case "error":
      return XCircleIcon;
    case "warning":
      return ExclamationTriangleIcon;
    case "info":
      return InformationCircleIcon;
  }
};

const removeToast = (id: number) => {
  notificationStore.removeToast(id);
};
</script>

<template>
  <div class="fixed top-4 right-4 z-[10000] flex flex-col space-y-3">
    <TransitionGroup name="toast">
      <div
        v-for="toast in toasts"
        :key="toast.id"
        class="relative w-80 bg-slate-800/95 backdrop-blur-sm rounded-xl shadow-2xl border border-slate-700/50 overflow-hidden"
        :class="{
          'shadow-green-500/20 border-green-500/30': toast.type === 'success',
          'shadow-red-500/20 border-red-500/30': toast.type === 'error',
          'shadow-blue-500/20 border-blue-500/30': toast.type === 'info',
          'shadow-yellow-500/20 border-yellow-500/30': toast.type === 'warning'
        }"
      >
        <div class="absolute top-0 left-0 w-full h-1 bg-slate-700/50">
          <div
            class="h-full transition-all duration-100 linear"
            :class="{
              'bg-green-500': toast.type === 'success',
              'bg-red-500': toast.type === 'error',
              'bg-blue-500': toast.type === 'info',
              'bg-yellow-500': toast.type === 'warning'
            }"
            :style="`animation: progress ${toast.duration}ms linear;`"
          ></div>
        </div>
        <div class="p-4 flex items-start">
          <div
            class="flex-shrink-0 w-6 h-6 rounded-full flex items-center justify-center mr-3 mt-0.5"
            :class="{
              'bg-green-500/20 text-green-400': toast.type === 'success',
              'bg-red-500/20 text-red-400': toast.type === 'error',
              'bg-blue-500/20 text-blue-400': toast.type === 'info',
              'bg-yellow-500/20 text-yellow-400': toast.type === 'warning'
            }"
          >
            <component :is="getIcon(toast.type)" class="w-4 h-4" />
          </div>

          <div class="flex-1 min-w-0">
            <h3 class="text-sm font-semibold text-slate-100 mb-1">
              {{ toast.title }}
            </h3>
            <p class="text-sm text-slate-400 leading-tight">
              {{ toast.message }}
            </p>
          </div>

          <button
            @click="removeToast(toast.id)"
            class="flex-shrink-0 ml-3 text-slate-400 hover:text-slate-200 transition-colors duration-200 p-1 rounded hover:bg-slate-700/50"
          >
            <XMarkIcon class="w-4 h-4" />
          </button>
        </div>
      </div>
    </TransitionGroup>
  </div>
</template>

<style>
@keyframes progress {
  from {
    width: 100%;
  }
  to {
    width: 0%;
  }
}

.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s ease;
}
.toast-enter-from,
.toast-leave-to {
  opacity: 0;
  transform: translateX(100px);
}
.toast-move {
  transition: transform 0.3s ease;
}
</style>
