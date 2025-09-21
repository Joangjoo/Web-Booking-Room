<script setup lang="ts">
import { ref, reactive, watch } from "vue";
import type { Room, RoomForm } from "../types";

const props = defineProps<{
  isOpen: boolean;
  isEditMode: boolean;
  initialData?: Room | null;
}>();

const emit = defineEmits<{
  (e: "close"): void;
  (e: "save", data: RoomForm): void;
}>();

const roomForm = reactive<RoomForm>({
  name: "",
  description: "",
  capacity: 1,
  features: [],
  is_available: true,
});
const featuresInput = ref("");

watch(
  () => props.initialData,
  (newData) => {
    if (props.isEditMode && newData) {
      roomForm.name = newData.name;
      roomForm.description = newData.description;
      roomForm.capacity = newData.capacity;
      roomForm.features = [...newData.features];
      roomForm.is_available = newData.is_available;
      featuresInput.value = newData.features ? newData.features.join(", ") : "";
    } else {
      roomForm.name = "";
      roomForm.description = "";
      roomForm.capacity = 1;
      roomForm.features = [];
      roomForm.is_available = true;
      featuresInput.value = "";
    }
  }
);

const toggleAvailability = (): void => {
  roomForm.is_available = !roomForm.is_available;
};

const parseFeatures = (input: string): string[] => {
  return input
    .split(",")
    .map((f) => f.trim())
    .filter((f) => f.length > 0);
};

const handleSubmit = () => {
  roomForm.features = parseFeatures(featuresInput.value);
  if (roomForm.features.length === 0) {
    alert("Mohon masukkan minimal satu fitur!");
    return;
  }
  emit("save", { ...roomForm });
};
</script>

<template>
  <div
    v-if="isOpen"
    class="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center p-4 z-50"
    @click="emit('close')"
  >
    <div
      @click.stop
      class="bg-slate-800/95 backdrop-blur-sm rounded-xl w-full max-w-2xl max-h-[90vh] overflow-y-auto"
    >
      <div class="flex items-center justify-between p-6 border-b border-slate-700/50">
        <h2 class="text-2xl font-bold bg-gradient-to-r from-blue-300 to-slate-200 bg-clip-text text-transparent">
          {{ isEditMode ? "✏️ Edit Ruangan" : "➕ Tambah Ruangan Baru" }}
        </h2>
        <button @click="emit('close')" class="text-slate-400 hover:text-slate-200">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>
        </button>
      </div>
      <form @submit.prevent="handleSubmit" class="p-6 space-y-6">
        <div>
          <label class="block text-sm font-medium text-slate-300 mb-2">Nama Ruangan</label>
          <input v-model="roomForm.name" type="text" placeholder="Masukkan nama ruangan..." class="w-full px-4 py-3 bg-slate-700/60 text-slate-100 placeholder-slate-400 border-2 border-slate-600 rounded-lg transition-all duration-300 outline-none focus:border-blue-500 focus:bg-slate-700/80 focus:ring-4 focus:ring-blue-500/20" required />
        </div>
        <div>
          <label class="block text-sm font-medium text-slate-300 mb-2">Deskripsi</label>
          <textarea v-model="roomForm.description" placeholder="Masukkan deskripsi ruangan..." rows="3" class="w-full px-4 py-3 bg-slate-700/60 text-slate-100 placeholder-slate-400 border-2 border-slate-600 rounded-lg transition-all duration-300 outline-none focus:border-blue-500 focus:bg-slate-700/80 focus:ring-4 focus:ring-blue-500/20 resize-none" required></textarea>
        </div>
        <div>
          <label class="block text-sm font-medium text-slate-300 mb-2">Kapasitas</label>
          <input v-model.number="roomForm.capacity" type="number" min="1" max="500" placeholder="Jumlah orang maksimal..." class="w-full px-4 py-3 bg-slate-700/60 text-slate-100 placeholder-slate-400 border-2 border-slate-600 rounded-lg transition-all duration-300 outline-none focus:border-blue-500 focus:bg-slate-700/80 focus:ring-4 focus:ring-blue-500/20" required />
        </div>
        <div>
          <label class="block text-sm font-medium text-slate-300 mb-2">Fitur-fitur</label>
          <input v-model="featuresInput" type="text" placeholder="Pisahkan dengan koma (contoh: Proyektor, AC, WiFi)" class="w-full px-4 py-3 bg-slate-700/60 text-slate-100 placeholder-slate-400 border-2 border-slate-600 rounded-lg transition-all duration-300 outline-none focus:border-blue-500 focus:bg-slate-700/80 focus:ring-4 focus:ring-blue-500/20" required />
          <p class="mt-1 text-xs text-slate-400">Contoh: Proyektor 4K, AC Central, WiFi Cepat</p>
        </div>
        <div class="flex items-center justify-between p-4 bg-slate-700/30 rounded-lg">
          <div>
            <h4 class="text-sm font-medium text-slate-200">Status Ketersediaan</h4>
            <p class="text-xs text-slate-400">Apakah ruangan dapat dibooking?</p>
          </div>
          <button @click="toggleAvailability" type="button" :class="['relative inline-flex h-6 w-11 items-center rounded-full transition-colors duration-300', roomForm.is_available ? 'bg-blue-600' : 'bg-slate-600']">
            <span :class="['inline-block h-4 w-4 transform rounded-full bg-white transition-transform duration-300', roomForm.is_available ? 'translate-x-6' : 'translate-x-1']" />
          </button>
        </div>

        <div class="flex items-center justify-end space-x-3 pt-4 border-t border-slate-700/50">
          <button type="button" @click="emit('close')" class="px-6 py-3 text-slate-300 hover:text-slate-100">Batal</button>
          <button type="submit" class="px-8 py-3 bg-gradient-to-r from-blue-600 to-blue-700 text-white font-semibold rounded-lg">
            {{ isEditMode ? '💾 Simpan Perubahan' : '➕ Tambah Ruangan' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>