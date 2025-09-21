<script setup lang="ts">
import type { Room } from "../types";

defineProps<{
  room: Room;
}>();

const emit = defineEmits<{
  (e: "view-details", room: Room): void;
}>();
</script>

<template>
  <div
    class="bg-slate-800/90 backdrop-blur-sm rounded-xl shadow-xl shadow-blue-950/20 border border-slate-700/50 overflow-hidden transform transition-all duration-300 hover:scale-105 hover:shadow-2xl hover:shadow-blue-900/30 cursor-pointer"
  >
    <div
      class="h-48 bg-gradient-to-br from-slate-700 via-blue-800 to-slate-800 relative"
    >
      <div class="absolute inset-0 flex items-center justify-center">
        <span class="text-6xl opacity-30">🏢</span>
      </div>
      <div
        v-if="!room.is_available"
        class="absolute inset-0 bg-black/60 flex items-center justify-center"
      >
        <span
          class="px-4 py-2 border border-yellow-500 text-yellow-500 rounded-lg font-medium"
          >Tidak Tersedia</span
        >
      </div>
    </div>

    <div class="p-6">
      <h3 class="text-xl font-bold text-slate-100 mb-2 truncate">
        {{ room.name }}
      </h3>

      <div class="flex items-center text-slate-400 text-sm mb-3">
        <span class="flex items-center">
          <svg
            class="w-4 h-4 mr-1"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"
            ></path>
          </svg>
          Kapasitas: {{ room.capacity }} Orang
        </span>
      </div>

      <p
        class="text-slate-300 text-sm leading-relaxed mb-4 h-20 overflow-hidden"
      >
        {{ room.description }}
      </p>

      <button
        @click="emit('view-details', room)"
        :disabled="!room.is_available"
        class="w-full py-3 bg-gradient-to-r from-blue-600 via-blue-700 to-slate-700 text-white font-semibold rounded-lg shadow-lg transform transition-all duration-300 hover:scale-105 hover:from-blue-500 hover:via-blue-600 hover:to-slate-600 hover:shadow-2xl hover:shadow-blue-600/30 active:scale-95 focus:outline-none focus:ring-4 focus:ring-blue-500/40 disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:scale-100"
      >
        📋 Lihat Detail & Jadwal
      </button>
    </div>
  </div>
</template>
