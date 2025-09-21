<script setup lang="ts">
import type { Room } from "../types";
defineProps<{ isOpen: boolean; room: Room | null }>();
const emit = defineEmits<{ (e: "close"): void; (e: "confirm"): void }>();
</script>

<template>
  <div
    v-if="isOpen"
    class="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center p-4 z-50"
    @click="emit('close')"
  >
    <div
      @click.stop
      class="bg-slate-800/95 rounded-xl w-full max-w-md p-6 text-center"
    >
      <h3 class="text-lg font-medium text-slate-200 mb-2">
        Konfirmasi Penghapusan
      </h3>
      <p class="text-sm text-slate-400 mb-6">
        Apakah Anda yakin ingin menghapus ruangan "<span
          class="font-semibold text-slate-200"
          >{{ room?.name }}</span
        >"?
      </p>
      <div class="flex items-center justify-center space-x-3">
        <button @click="emit('close')" class="px-6 py-3 text-slate-300">
          Batal
        </button>
        <button
          @click="emit('confirm')"
          class="px-6 py-3 bg-gradient-to-r from-red-600 to-red-700 text-white font-semibold rounded-lg"
        >
          🗑️ Ya, Hapus
        </button>
      </div>
    </div>
  </div>
</template>
