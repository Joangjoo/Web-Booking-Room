<script setup lang="ts">
import { ref, onMounted } from "vue";
import type { Room, RoomForm, BookingDetails } from "../types";
import { useRoomStore } from "../stores/roomStore";
import { useAuthStore } from "../stores/authStore";
import RoomFormModal from "../components/RoomFormModal.vue";
import DeleteConfirmationModal from "../components/DeleteConfirmationModal.vue";
import { useNotificationStore } from "../stores/notificationStore";
import { useRouter } from "vue-router";
import apiClient from "../api";

const roomStore = useRoomStore();
const authStore = useAuthStore();
const isFormModalOpen = ref(false);
const isDeleteModalOpen = ref(false);
const isEditMode = ref(false);
const selectedRoom = ref<Room | null>(null);
const notificationStore = useNotificationStore();
const router = useRouter();
const bookingDetails = ref<BookingDetails[]>([]);
const showBookingDetails = ref(false);

const fetchBookings = async () => {
  try {
    const response = await apiClient.get<BookingDetails[]>(
      "/protected/bookings"
    );
    bookingDetails.value = response.data || [];
  } catch (error) {
    console.error("Gagal mengambil data booking:", error);
  }
};

const formatDateTime = (isoString: string) => {
  const date = new Date(isoString);
  return date.toLocaleString("id-ID", {
    dateStyle: "medium",
    timeStyle: "short",
  });
};

onMounted(() => {
  roomStore.fetchRooms();
  fetchBookings();
});

const openAddModal = () => {
  isEditMode.value = false;
  selectedRoom.value = null;
  isFormModalOpen.value = true;
};

const openEditModal = (room: Room) => {
  isEditMode.value = true;
  selectedRoom.value = room;
  isFormModalOpen.value = true;
};

const openDeleteModal = (room: Room) => {
  selectedRoom.value = room;
  isDeleteModalOpen.value = true;
};

const handleSaveRoom = async (formData: RoomForm) => {
  try {
    if (isEditMode.value && selectedRoom.value) {
      await roomStore.updateRoom(selectedRoom.value.id, formData);

      notificationStore.showToast({
        type: "success",
        title: "Berhasil!",
        message: `Perubahan pada '${formData.name}' telah disimpan.`,
      });
    } else {
      await roomStore.createRoom(formData);

      notificationStore.showToast({
        type: "success",
        title: "Berhasil!",
        message: `Ruangan '${formData.name}' berhasil ditambahkan.`,
      });
    }
    isFormModalOpen.value = false;
  } catch (error: any) {
    notificationStore.showToast({
      type: "error",
      title: "Gagal!",
      message: error.message || "Gagal menyimpan ruangan.",
    });
  }
};

const confirmDelete = async () => {
  if (selectedRoom.value) {
    try {
      await roomStore.deleteRoom(selectedRoom.value.id);
      notificationStore.showToast({
        type: "success",
        title: "Berhasil!",
        message: `Ruangan '${selectedRoom.value.name}' telah dihapus.`,
      });
    } catch (error: any) {
      notificationStore.showToast({
        type: "error",
        title: "Gagal!",
        message: error.message || "Gagal menghapus ruangan.",
      });
    }
  }
  isDeleteModalOpen.value = false;
};



const handleLogout = () => {
  authStore.logout();
  router.push("/login");
};
</script>

<template>
  <div class="min-h-screen bg-gradient-to-br from-slate-900 via-blue-900 to-black">
    <div class="relative z-10 max-w-7xl mx-auto p-8">
      
      <!-- Header -->
      <div class="flex items-center justify-between mb-8">
        <div>
          <h1 class="text-4xl font-bold bg-gradient-to-r from-blue-300 to-blue-400 bg-clip-text text-transparent">
            🏢 Manajemen Ruangan
          </h1>
          <p class="text-slate-400">Kelola semua ruangan dan riwayat booking.</p>
        </div>
        <div class="flex items-center space-x-4">
          <button
            @click="openAddModal"
            class="px-6 py-3 bg-gradient-to-r from-blue-600 via-blue-700 to-slate-700 text-white font-semibold rounded-xl shadow-lg transform transition-all duration-300 hover:scale-105 active:scale-95"
          >
            Tambah Ruangan
          </button>
          <button
            @click="handleLogout"
            class="px-6 py-3 bg-gradient-to-r from-red-600 via-red-700 to-slate-700 text-white font-semibold rounded-xl shadow-lg transform transition-all duration-300 hover:scale-105 active:scale-95"
          >
            Logout
          </button>
        </div>
      </div>

      <!-- Tabel Ruangan -->
      <div class="bg-slate-800/90 rounded-xl shadow-2xl border border-slate-700/50 overflow-hidden">
        <div class="p-4 bg-slate-700/50 border-b border-slate-700/50 flex justify-between items-center">
          <h3 class="font-semibold text-slate-200">Daftar Ruangan</h3>
          <button
            @click="showBookingDetails = !showBookingDetails"
            class="px-4 py-2 text-sm bg-gradient-to-r from-indigo-600 to-indigo-800 text-white rounded-lg shadow-md hover:scale-105 transition"
          >
            {{ showBookingDetails ? "Tutup Riwayat Booking" : "Lihat Riwayat Booking" }}
          </button>
        </div>
        <table class="min-w-full">
          <thead class="bg-slate-700/30">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-slate-300 uppercase">ID</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-slate-300 uppercase">Nama Ruangan</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-slate-300 uppercase">Kapasitas</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-slate-300 uppercase">Status</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-slate-300 uppercase">Fitur</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-slate-300 uppercase">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-700/50">
            <tr
              v-for="room in roomStore.rooms"
              :key="room.id"
              class="hover:bg-slate-700/30"
            >
              <td class="px-6 py-4 text-sm font-mono text-slate-300">{{ room.id }}</td>
              <td class="px-6 py-4 text-sm text-slate-100">{{ room.name }}</td>
              <td class="px-6 py-4 text-sm text-slate-300">{{ room.capacity }} orang</td>
              <td class="px-6 py-4">
                <span
                  :class="[
                    'px-2.5 py-0.5 rounded-full text-xs',
                    room.is_available ? 'bg-green-100/20 text-green-400' : 'bg-red-100/20 text-red-400'
                  ]"
                >
                  {{ room.is_available ? "Tersedia" : "Tidak Tersedia" }}
                </span>
              </td>
              <td class="px-6 py-4">
                <div class="flex flex-wrap gap-1 max-w-xs">
                  <span
                    v-for="feature in (room.features || []).slice(0, 2)"
                    :key="feature"
                    class="px-2 py-1 rounded-full text-xs bg-blue-600/20 text-blue-300"
                  >
                    {{ feature }}
                  </span>
                  <span
                    v-if="(room.features || []).length > 2"
                    class="px-2 py-1 rounded-full text-xs bg-slate-600/50 text-slate-400"
                  >
                    +{{ (room.features || []).length - 2 }} lagi
                  </span>
                </div>
              </td>
              <td class="px-6 py-4">
                <button @click="openEditModal(room)" class="text-blue-400 hover:text-blue-300">Edit</button>
                <button @click="openDeleteModal(room)" class="ml-4 text-red-400 hover:text-red-300">Hapus</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Riwayat Booking (Toggle) -->
      <div
        v-if="showBookingDetails"
        class="mt-8 bg-slate-800/90 rounded-xl shadow-2xl border border-slate-700/50 overflow-hidden"
      >
        <div class="p-4 bg-slate-700/50 border-b border-slate-700/50">
          <h3 class="font-semibold text-slate-200">Riwayat Booking</h3>
        </div>
        <table class="min-w-full">
          <thead class="bg-slate-700/30">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-slate-300 uppercase">Booking ID</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-slate-300 uppercase">Nama Ruangan</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-slate-300 uppercase">Nama Pemesan</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-slate-300 uppercase">Waktu Mulai</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-slate-300 uppercase">Waktu Selesai</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-slate-300 uppercase">Status</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-700/50">
            <tr
              v-for="booking in bookingDetails"
              :key="booking.id"
              class="hover:bg-slate-700/30"
            >
              <td class="px-6 py-4 text-sm font-mono text-slate-300">{{ booking.id }}</td>
              <td class="px-6 py-4 text-sm text-slate-100">{{ booking.room_name }}</td>
              <td class="px-6 py-4 text-sm text-slate-300">{{ booking.user_name }}</td>
              <td class="px-6 py-4 text-sm text-slate-300">{{ formatDateTime(booking.start_time) }}</td>
              <td class="px-6 py-4 text-sm text-slate-300">{{ formatDateTime(booking.end_time) }}</td>
              <td class="px-6 py-4 text-sm">
                <span
                  class="px-2.5 py-0.5 rounded-full text-xs font-medium capitalize"
                  :class="booking.status === 'confirmed'
                    ? 'bg-green-100/20 text-green-400'
                    : 'bg-yellow-100/20 text-yellow-400'"
                >
                  {{ booking.status }}
                </span>
              </td>
            </tr>
          </tbody>
        </table>
        <div
          v-if="bookingDetails.length === 0"
          class="text-center py-12 text-slate-400"
        >
          <p>Belum ada riwayat booking.</p>
        </div>
      </div>
    </div>

    <!-- Modal -->
    <RoomFormModal
      :is-open="isFormModalOpen"
      :is-edit-mode="isEditMode"
      :initial-data="selectedRoom"
      @close="isFormModalOpen = false"
      @save="handleSaveRoom"
    />
    <DeleteConfirmationModal
      :is-open="isDeleteModalOpen"
      :room="selectedRoom"
      @close="isDeleteModalOpen = false"
      @confirm="confirmDelete"
    />
  </div>
</template>
