<script setup lang="ts">
import { ref, reactive, onMounted } from "vue";
import axios from "axios";
import type { Room, RoomForm } from "../types";
import { useRouter } from "vue-router";

const rooms = ref<Room[]>([]);
const isFormModalOpen = ref<boolean>(false);
const isDeleteModalOpen = ref<boolean>(false);
const isEditMode = ref<boolean>(false);
const isSubmitting = ref<boolean>(false);
const isDeleting = ref<boolean>(false);
const router = useRouter();

const roomForm = reactive<RoomForm>({
    name: "",
    description: "",
    capacity: 1,
    features: [],
    is_available: true,
});

const featuresInput = ref<string>("");
const roomToDelete = ref<Room | null>(null);
const editingRoomId = ref<number | null>(null);

const fetchRooms = async () => {
    try {
        const response = await axios.get("http://localhost:8080/api/rooms");
        rooms.value = response.data || [];
    } catch (error) {
        console.error("Gagal mengambil data ruangan:", error);
        alert("Gagal memuat data ruangan dari server.");
    }
};

onMounted(() => {
    fetchRooms();
});

const openAddModal = (): void => {
    console.log("openAddModal called - isFormModalOpen will be true");
    isEditMode.value = false;
    resetForm();
    isFormModalOpen.value = true;
    console.log("Modal state:", isFormModalOpen.value);
};
const openEditModal = (room: Room): void => {
    isEditMode.value = true;
    editingRoomId.value = room.id;
    roomForm.name = room.name;
    roomForm.description = room.description;
    roomForm.capacity = room.capacity;
    roomForm.features = [...room.features];
    roomForm.is_available = room.is_available;
    featuresInput.value = room.features ? room.features.join(", ") : "";

    isFormModalOpen.value = true;
};

const closeFormModal = (): void => {
    isFormModalOpen.value = false;
    resetForm();
};

const resetForm = (): void => {
    roomForm.name = "";
    roomForm.description = "";
    roomForm.capacity = 1;
    roomForm.features = [];
    roomForm.is_available = true;
    featuresInput.value = "";
    editingRoomId.value = null;
};

const openDeleteModal = (room: Room): void => {
    roomToDelete.value = room;
    isDeleteModalOpen.value = true;
};

const closeDeleteModal = (): void => {
    isDeleteModalOpen.value = false;
    roomToDelete.value = null;
};

const toggleAvailability = (): void => {
    roomForm.is_available = !roomForm.is_available;
};

const parseFeatures = (input: string): string[] => {
    return input
    .split(",")
    .map((feature) => feature.trim())
    .filter((feature) => feature.length > 0);
};

const handleSubmit = async (): Promise<void> => {
    roomForm.features = parseFeatures(featuresInput.value);
        if (roomForm.features.length === 0) {
        alert("Mohon masukkan minimal satu fitur!");
    return;
    }
    isSubmitting.value = true;

    try {
    if (isEditMode.value && editingRoomId.value) {
        await axios.put(
        `http://localhost:8080/api/protected/rooms/${editingRoomId.value}`,
        roomForm
        );
        alert("✅ Ruangan berhasil diperbarui!");
    } else {
        await axios.post("http://localhost:8080/api/protected/rooms", roomForm);
        alert("✅ Ruangan baru berhasil ditambahkan!");
    }
    closeFormModal();
    await fetchRooms();
    } catch (error) {
    console.error("Error saving room:", error);
    alert("❌ Terjadi kesalahan saat menyimpan ruangan!");
    } finally {
    isSubmitting.value = false;
    }
};

const confirmDelete = async (): Promise<void> => {
    if (!roomToDelete.value) return;
    isDeleting.value = true;

    try {
    await axios.delete(
        `http://localhost:8080/api/protected/rooms/${roomToDelete.value.id}`
    );
    alert(`✅ Ruangan "${roomToDelete.value.name}" berhasil dihapus!`);

    closeDeleteModal();
    await fetchRooms();
    } catch (error) {
    console.error("Error deleting room:", error);
    alert("❌ Terjadi kesalahan saat menghapus ruangan!");
    } finally {
    isDeleting.value = false;
    }
};

const handleLogout = () => {
  localStorage.removeItem("token"); 
  alert("✅ Anda telah logout!");
  router.push("/login"); 
};
</script>

<template>
    <div
    class="min-h-screen bg-gradient-to-br from-slate-900 via-blue-900 to-black"
    >
    <div
        class="absolute inset-0 bg-gradient-to-tr from-blue-950/40 via-slate-800/30 to-gray-900/50"
    ></div>

    <div class="relative z-10 max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div class="flex items-center justify-between mb-8">
        <div>
            <h1
            class="text-4xl font-bold bg-gradient-to-r from-blue-300 via-slate-200 to-blue-400 bg-clip-text text-transparent mb-2">
            🏢 Manajemen Ruangan</h1>
            <p class="text-slate-400 text-lg">Kelola semua ruangan dengan mudah dan efisien</p>
        </div>

        <button
            @click="openAddModal"
            class="px-6 py-4 bg-gradient-to-r from-blue-600 via-blue-700 to-slate-700 text-white font-semibold rounded-xl shadow-lg transform transition-all duration-300 hover:scale-105 hover:from-blue-500 hover:via-blue-600 hover:to-slate-600 hover:shadow-2xl hover:shadow-blue-600/30 active:scale-95 focus:outline-none focus:ring-4 focus:ring-blue-500/40 flex items-center space-x-3"
        >
        <svg
            class="w-5 h-5"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
        >
            <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M12 6v6m0 0v6m0-6h6m-6 0H6"
            ></path>
        </svg>
        <span>Tambah Ruangan Baru</span>
        </button>

        <button
      @click="handleLogout"
      class="px-6 py-4 bg-gradient-to-r from-red-600 via-red-700 to-slate-700 text-white font-semibold rounded-xl shadow-lg transform transition-all duration-300 hover:scale-105 hover:from-red-500 hover:via-red-600 hover:to-slate-600 hover:shadow-2xl hover:shadow-red-600/30 active:scale-95 focus:outline-none focus:ring-4 focus:ring-red-500/40 flex items-center space-x-3"
    >
      <svg
        class="w-5 h-5"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M17 16l4-4m0 0l-4-4m4 4H7"
        ></path>
      </svg>
      <span>Logout</span>
    </button>
    </div>

      <div
        class="bg-slate-800/90 backdrop-blur-sm rounded-xl shadow-2xl shadow-blue-950/20 border border-slate-700/50 overflow-hidden"
      >
        <div class="px-6 py-4 border-b border-slate-700/50 bg-slate-700/50">
          <div class="flex items-center justify-between">
            <h2 class="text-xl font-semibold text-slate-200">
              📋 Daftar Ruangan
            </h2>
            <span class="text-sm text-slate-400">
              Total: {{ rooms.length }} ruangan
            </span>
          </div>
        </div>

        <div class="overflow-x-auto">
          <table class="min-w-full divide-y divide-slate-700/50">
            <thead class="bg-slate-700/30">
              <tr>
                <th
                  class="px-6 py-4 text-left text-xs font-medium text-slate-300 uppercase tracking-wider"
                >
                  ID
                </th>
                <th
                  class="px-6 py-4 text-left text-xs font-medium text-slate-300 uppercase tracking-wider"
                >
                  Nama Ruangan
                </th>
                <th
                  class="px-6 py-4 text-left text-xs font-medium text-slate-300 uppercase tracking-wider"
                >
                  Kapasitas
                </th>
                <th
                  class="px-6 py-4 text-left text-xs font-medium text-slate-300 uppercase tracking-wider"
                >
                  Status
                </th>
                <th
                  class="px-6 py-4 text-left text-xs font-medium text-slate-300 uppercase tracking-wider"
                >
                  Fitur
                </th>
                <th
                  class="px-6 py-4 text-left text-xs font-medium text-slate-300 uppercase tracking-wider"
                >
                  Aksi
                </th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-700/50">
              <tr
                v-for="room in rooms"
                :key="room.id"
                class="transition-colors duration-200 hover:bg-slate-700/30 group"
              >
                <td
                  class="px-6 py-4 whitespace-nowrap text-sm font-mono text-slate-300"
                >
                  {{ room.id }}
                </td>

                <td class="px-6 py-4 whitespace-nowrap">
                  <div class="text-sm font-medium text-slate-100">
                    {{ room.name }}
                  </div>
                  <div class="text-sm text-slate-400 truncate max-w-xs">
                    {{ room.description }}
                  </div>
                </td>

                <td class="px-6 py-4 whitespace-nowrap">
                  <div class="flex items-center text-sm text-slate-300">
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
                        d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"
                      ></path>
                    </svg>
                    {{ room.capacity }} orang
                  </div>
                </td>

                <td class="px-6 py-4 whitespace-nowrap">
                  <span
                    :class="[
                      'inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium',
                      room.is_available
                        ? 'bg-green-100/20 text-green-400 border border-green-400/30'
                        : 'bg-red-100/20 text-red-400 border border-red-400/30',
                    ]"
                  >
                    <div
                      :class="[
                        'w-2 h-2 rounded-full mr-1.5',
                        room.is_available ? 'bg-green-400' : 'bg-red-400',
                      ]"
                    ></div>
                    {{ room.is_available ? "Tersedia" : "Tidak Tersedia" }}
                  </span>
                </td>

                <span
                  v-for="feature in (room.features || []).slice(0, 3)"
                  :key="feature"
                  class="inline-flex items-center px-2 py-1 rounded-full text-xs bg-blue-600/20 text-blue-300 border border-blue-600/30"
                >
                  {{ feature }}
                </span>

                <span
                  v-if="(room.features || []).length > 3"
                  class="inline-flex items-center px-2 py-1 rounded-full text-xs bg-slate-600/50 text-slate-400 border border-slate-600/50"
                >
                  +{{ (room.features || []).length - 3 }} lagi
                </span>

                <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                  <div class="flex items-center space-x-2">
                    <button
                      @click="openEditModal(room)"
                      class="inline-flex items-center px-3 py-2 text-sm text-blue-400 hover:text-blue-300 hover:bg-blue-600/20 rounded-lg transition-all duration-200 border border-blue-600/30 hover:border-blue-500/50"
                      title="Edit Ruangan"
                    >
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
                          d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
                        ></path>
                      </svg>
                      Edit
                    </button>

                    <button
                      @click="openDeleteModal(room)"
                      class="inline-flex items-center px-3 py-2 text-sm text-red-400 hover:text-red-300 hover:bg-red-600/20 rounded-lg transition-all duration-200 border border-red-600/30 hover:border-red-500/50"
                      title="Hapus Ruangan"
                    >
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
                          d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                        ></path>
                      </svg>
                      Hapus
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>

          <div v-if="rooms.length === 0" class="text-center py-12">
            <div class="text-6xl mb-4">🏢</div>
            <h3 class="text-xl font-medium text-slate-300 mb-2">
              Belum ada ruangan
            </h3>
            <p class="text-slate-400 mb-6">
              Tambahkan ruangan pertama dengan klik tombol di atas
            </p>
            <button
              @click="openAddModal"
              class="px-6 py-3 bg-gradient-to-r from-blue-600 to-blue-700 text-white font-medium rounded-lg transition-all duration-300 hover:from-blue-500 hover:to-blue-600"
            >
              Tambah Ruangan
            </button>
          </div>
        </div>
      </div>
    </div>

    <div
      v-if="isFormModalOpen"
      class="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center p-4 z-50"
      @click="closeFormModal"
    >
      <div
        @click.stop
        class="bg-slate-800/95 backdrop-blur-sm rounded-xl shadow-2xl shadow-blue-950/40 border border-slate-700/50 w-full max-w-2xl max-h-[90vh] overflow-y-auto transform transition-all duration-300"
      >
        <div
          class="flex items-center justify-between p-6 border-b border-slate-700/50"
        >
          <h2
            class="text-2xl font-bold bg-gradient-to-r from-blue-300 to-slate-200 bg-clip-text text-transparent"
          >
            {{ isEditMode ? "✏️ Edit Ruangan" : "➕ Tambah Ruangan Baru" }}
          </h2>
          <button
            @click="closeFormModal"
            class="text-slate-400 hover:text-slate-200 transition-colors duration-300 p-2 hover:bg-slate-700/50 rounded-lg"
          >
            <svg
              class="w-6 h-6"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M6 18L18 6M6 6l12 12"
              ></path>
            </svg>
          </button>
        </div>

        <form @submit.prevent="handleSubmit" class="p-6 space-y-6">
          <div>
            <label class="block text-sm font-medium text-slate-300 mb-2">
              Nama Ruangan
            </label>
            <input
              v-model="roomForm.name"
              type="text"
              placeholder="Masukkan nama ruangan..."
              class="w-full px-4 py-3 bg-slate-700/60 text-slate-100 placeholder-slate-400 border-2 border-slate-600 rounded-lg transition-all duration-300 outline-none focus:border-blue-500 focus:bg-slate-700/80 focus:ring-4 focus:ring-blue-500/20"
              required
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-slate-300 mb-2">
              Deskripsi
            </label>
            <textarea
              v-model="roomForm.description"
              placeholder="Masukkan deskripsi ruangan..."
              rows="3"
              class="w-full px-4 py-3 bg-slate-700/60 text-slate-100 placeholder-slate-400 border-2 border-slate-600 rounded-lg transition-all duration-300 outline-none focus:border-blue-500 focus:bg-slate-700/80 focus:ring-4 focus:ring-blue-500/20 resize-none"
              required
            ></textarea>
          </div>

          <div>
            <label class="block text-sm font-medium text-slate-300 mb-2">
              Kapasitas
            </label>
            <input
              v-model.number="roomForm.capacity"
              type="number"
              min="1"
              max="500"
              placeholder="Jumlah orang maksimal..."
              class="w-full px-4 py-3 bg-slate-700/60 text-slate-100 placeholder-slate-400 border-2 border-slate-600 rounded-lg transition-all duration-300 outline-none focus:border-blue-500 focus:bg-slate-700/80 focus:ring-4 focus:ring-blue-500/20"
              required
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-slate-300 mb-2">
              Fitur-fitur
            </label>
            <input
              v-model="featuresInput"
              type="text"
              placeholder="Pisahkan dengan koma (contoh: Proyektor, AC, WiFi)"
              class="w-full px-4 py-3 bg-slate-700/60 text-slate-100 placeholder-slate-400 border-2 border-slate-600 rounded-lg transition-all duration-300 outline-none focus:border-blue-500 focus:bg-slate-700/80 focus:ring-4 focus:ring-blue-500/20"
              required
            />
            <p class="mt-1 text-xs text-slate-400">
              Contoh: Proyektor 4K, AC Central, WiFi Cepat, Whiteboard
            </p>
          </div>
          <div
            class="flex items-center justify-between p-4 bg-slate-700/30 rounded-lg"
          >
            <div>
              <h4 class="text-sm font-medium text-slate-200">
                Status Ketersediaan
              </h4>
              <p class="text-xs text-slate-400">
                Apakah ruangan dapat dibooking?
              </p>
            </div>
            <button
              @click="toggleAvailability"
              type="button"
              :class="[
                'relative inline-flex h-6 w-11 items-center rounded-full transition-colors duration-300',
                roomForm.is_available ? 'bg-blue-600' : 'bg-slate-600',
              ]"
            >
              <span
                :class="[
                  'inline-block h-4 w-4 transform rounded-full bg-white transition-transform duration-300',
                  roomForm.is_available ? 'translate-x-6' : 'translate-x-1',
                ]"
              />
            </button>
          </div>
          <div
            class="flex items-center justify-end space-x-3 pt-4 border-t border-slate-700/50"
          >
            <button
              type="button"
              @click="closeFormModal"
              class="px-6 py-3 text-slate-300 hover:text-slate-100 transition-colors duration-300"
            >
              Batal
            </button>
            <button
              type="submit"
              :disabled="isSubmitting"
              class="px-8 py-3 bg-gradient-to-r from-blue-600 via-blue-700 to-slate-700 text-white font-semibold rounded-lg shadow-lg transform transition-all duration-300 hover:scale-105 hover:from-blue-500 hover:via-blue-600 hover:to-slate-600 hover:shadow-2xl hover:shadow-blue-600/30 active:scale-95 focus:outline-none focus:ring-4 focus:ring-blue-500/40 disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:scale-100"
            >
              <span v-if="!isSubmitting">
                {{ isEditMode ? "💾 Simpan Perubahan" : "➕ Tambah Ruangan" }}
              </span>
              <span v-else class="flex items-center">
                <svg
                  class="animate-spin -ml-1 mr-3 h-4 w-4 text-white"
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                >
                  <circle
                    class="opacity-25"
                    cx="12"
                    cy="12"
                    r="10"
                    stroke="currentColor"
                    stroke-width="4"
                  ></circle>
                  <path
                    class="opacity-75"
                    fill="currentColor"
                    d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                  ></path>
                </svg>
                Menyimpan...
              </span>
            </button>
          </div>
        </form>
      </div>
    </div>

    <div
      v-if="isDeleteModalOpen"
      class="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center p-4 z-50"
      @click="closeDeleteModal"
    >
      <div
        @click.stop
        class="bg-slate-800/95 backdrop-blur-sm rounded-xl shadow-2xl shadow-red-950/40 border border-slate-700/50 w-full max-w-md transform transition-all duration-300"
      >
        <div class="p-6 text-center">
          <div
            class="mx-auto flex items-center justify-center h-12 w-12 rounded-full bg-red-600/20 mb-4"
          >
            <svg
              class="h-6 w-6 text-red-500"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z"
              ></path>
            </svg>
          </div>

          <h3 class="text-lg font-medium text-slate-200 mb-2">
            Konfirmasi Penghapusan
          </h3>

          <p class="text-sm text-slate-400 mb-6">
            Apakah Anda yakin ingin menghapus ruangan
            <span class="font-semibold text-slate-200"
              >"{{ roomToDelete?.name }}"</span
            >? <br /><br />
            Aksi ini tidak dapat dibatalkan.
          </p>

          <div class="flex items-center justify-center space-x-3">
            <button
              @click="closeDeleteModal"
              class="px-6 py-3 text-slate-300 hover:text-slate-100 transition-colors duration-300"
            >
              Batal
            </button>
            <button
              @click="confirmDelete"
              :disabled="isDeleting"
              class="px-6 py-3 bg-gradient-to-r from-red-600 to-red-700 text-white font-semibold rounded-lg shadow-lg transform transition-all duration-300 hover:scale-105 hover:from-red-500 hover:to-red-600 hover:shadow-2xl hover:shadow-red-600/30 active:scale-95 focus:outline-none focus:ring-4 focus:ring-red-500/40 disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:scale-100"
            >
              <span v-if="!isDeleting">🗑️ Ya, Hapus</span>
              <span v-else class="flex items-center">
                <svg
                  class="animate-spin -ml-1 mr-3 h-4 w-4 text-white"
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                >
                  <circle
                    class="opacity-25"
                    cx="12"
                    cy="12"
                    r="10"
                    stroke="currentColor"
                    stroke-width="4"
                  ></circle>
                  <path
                    class="opacity-75"
                    fill="currentColor"
                    d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                  ></path>
                </svg>
                Menghapus...
              </span>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
