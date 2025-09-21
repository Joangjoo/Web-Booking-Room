<script setup lang="ts">
import { ref, reactive, computed, onMounted } from "vue";
import type {
  Room,
  BookingForm,
  CapacityFilter,
} from "../types";
import RoomCard from "../components/RoomCard.vue";
import apiClient from "../api";
import { useAuthStore } from '../stores/authStore';
import router from "../router";
import { useNotificationStore } from "../stores/notificationStore";



const searchQuery = ref<string>("");
const selectedCapacityFilter = ref<string>("all");
const isSearchFocused = ref<boolean>(false);
const isModalOpen = ref<boolean>(false);
const selectedRoom = ref<Room | null>(null);
const isBookingLoading = ref<boolean>(false);
const authStore = useAuthStore();
const notificationStore = useNotificationStore();

const bookingForm = reactive<BookingForm>({
  roomId: null,
  date: "",
  startTime: "",
  endTime: "",
  purpose: "",
});

const capacityFilters: CapacityFilter[] = [
  { value: "all", label: "Semua Kapasitas" },
  { value: "small", label: "< 10 Orang" },
  { value: "medium", label: "10-20 Orang" },
  { value: "large", label: "> 20 Orang" },
];

const rooms = ref<Room[]>([]);

const fetchRooms = async () => {
  try {
    const response = await apiClient.get<Room[]>('/rooms'); 
    rooms.value = response.data || []; 
  } catch (error) {
    console.error("Gagal mengambil data ruangan:", error);
    alert("Gagal memuat data ruangan dari server.");
  }
};

onMounted(() => {
  fetchRooms();
});

const filteredRooms = computed<Room[]>(() => {
  let filtered = rooms.value.filter((room) => room.is_available);
  if (searchQuery.value.trim()) {
    const query = searchQuery.value.toLowerCase();
    filtered = filtered.filter(
      (room) =>
        room.name.toLowerCase().includes(query) ||
        room.description.toLowerCase().includes(query) ||
        room.features.some((feature) => feature.toLowerCase().includes(query))
    );
  }

  if (selectedCapacityFilter.value !== "all") {
    filtered = filtered.filter((room) => {
      switch (selectedCapacityFilter.value) {
        case "small":
          return room.capacity < 10;
        case "medium":
          return room.capacity >= 10 && room.capacity <= 20;
        case "large":
          return room.capacity > 20;
        default:
          return true;
      }
    });
  }

  return filtered;
});

const isBookingFormValid = computed<boolean>(() => {
  return !!(
    bookingForm.date &&
    bookingForm.startTime &&
    bookingForm.endTime &&
    bookingForm.purpose.trim() &&
    bookingForm.startTime < bookingForm.endTime
  );
});

const setSearchFocused = (focused: boolean): void => {
  isSearchFocused.value = focused;
};

const setCapacityFilter = (filter: string): void => {
  selectedCapacityFilter.value = filter;
};

const openBookingModal = (room: Room): void => {
  selectedRoom.value = room;
  bookingForm.roomId = room.id;
  isModalOpen.value = true;

  bookingForm.date = "";
  bookingForm.startTime = "";
  bookingForm.endTime = "";
  bookingForm.purpose = "";
};

const closeBookingModal = (): void => {
  isModalOpen.value = false;
  selectedRoom.value = null;
  isBookingLoading.value = false;
};

const handleBooking = async (): Promise<void> => {
  if (!isBookingFormValid.value) {
    alert("Mohon lengkapi semua field booking!");
    return;
  }
  if (bookingForm.startTime >= bookingForm.endTime) {
    alert("Waktu selesai harus setelah waktu mulai!");
    return;
  }

  isBookingLoading.value = true;

  try {
    const startTimeISO = new Date(`${bookingForm.date}T${bookingForm.startTime}`).toISOString();
    const endTimeISO = new Date(`${bookingForm.date}T${bookingForm.endTime}`).toISOString();

    const payload = {
      room_id: bookingForm.roomId,
      start_time: startTimeISO,
      end_time: endTimeISO,
      purpose: bookingForm.purpose,
    };

    const response = await apiClient.post('/protected/bookings', payload);

    notificationStore.showToast({
      type: 'success',
      title: 'Berhasil!',
      message: `Ruangan '${response.data.room_name}' berhasil di-booking.`
    });
    closeBookingModal();

  } catch (error: any) {
    console.error("Booking error:", error);
    const errorMessage = error.response?.data?.error || "Terjadi kesalahan saat booking.";
    alert(`❌ Gagal Booking: ${errorMessage}`);
  } finally {
    isBookingLoading.value = false;
  }
};

const handleLogout = (): void => {
  authStore.logout(); 
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
    <header
      class="relative z-20 bg-slate-800/90 backdrop-blur-sm border-b border-slate-700/50 shadow-lg"
    >
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex items-center justify-between h-16">
          <div class="flex items-center">
            <h1
              class="text-2xl font-bold bg-gradient-to-r from-blue-300 via-slate-200 to-blue-400 bg-clip-text text-transparent"
            >
              📅 Booking Room 
            </h1>
          </div>
          <div class="flex items-center space-x-4">
            <span class="text-slate-300 text-sm" v-if="authStore.user">
              Halo,
              <span class="font-medium text-blue-400">{{
                authStore.user.name
              }}</span>
              👋
            </span>
            <button
              @click="handleLogout"
              class="px-4 py-2 bg-gradient-to-r from-red-600 to-red-700 text-white text-sm font-medium rounded-lg transition-all duration-300 hover:scale-105 hover:shadow-lg hover:shadow-red-600/30 active:scale-95 focus:outline-none focus:ring-4 focus:ring-red-500/40"
            >
              Logout
            </button>
          </div>
        </div>
      </div>
    </header>

    <!-- Main Content -->
    <main class="relative z-10 max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div
        class="bg-slate-800/90 backdrop-blur-sm rounded-xl shadow-xl shadow-blue-950/20 p-6 border border-slate-700/50 mb-8"
      >
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div class="relative">
            <label class="block text-sm font-medium text-slate-300 mb-3">
              🔍 Cari Ruangan
            </label>
            <div class="relative">
              <input
                v-model="searchQuery"
                type="text"
                placeholder="Masukkan nama ruangan..."
                @focus="setSearchFocused(true)"
                @blur="setSearchFocused(false)"
                class="w-full px-4 py-3 pl-10 bg-slate-700/60 text-slate-100 placeholder-slate-400 border-2 border-slate-600 rounded-lg transition-all duration-300 outline-none focus:border-blue-500 focus:bg-slate-700/80 focus:ring-4 focus:ring-blue-500/20 hover:border-slate-500"
              />
              <div class="absolute left-3 top-1/2 transform -translate-y-1/2">
                <svg
                  class="w-5 h-5 text-slate-400"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
                  ></path>
                </svg>
              </div>
              <div
                :class="[
                  'absolute inset-0 rounded-lg bg-gradient-to-r from-blue-500/5 via-slate-500/5 to-blue-600/5 pointer-events-none transition-opacity duration-300',
                  isSearchFocused ? 'opacity-100' : 'opacity-0',
                ]"
              ></div>
            </div>
          </div>
          <div class="relative">
            <label class="block text-sm font-medium text-slate-300 mb-3">
              👥 Filter Kapasitas
            </label>
            <div class="flex flex-wrap gap-2">
              <button
                v-for="filter in capacityFilters"
                :key="filter.value"
                @click="setCapacityFilter(filter.value)"
                :class="[
                  'px-4 py-2 text-sm font-medium rounded-lg transition-all duration-300 border-2',
                  selectedCapacityFilter === filter.value
                    ? 'bg-gradient-to-r from-blue-600 to-blue-700 text-white border-blue-500 shadow-lg shadow-blue-600/30'
                    : 'bg-slate-700/60 text-slate-300 border-slate-600 hover:border-blue-500 hover:bg-slate-700/80 hover:text-blue-300',
                ]"
              >
                {{ filter.label }}
              </button>
            </div>
          </div>
        </div>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <RoomCard
          v-for="room in filteredRooms"
          :key="room.id"
          :room="room"
          @view-details="openBookingModal"
        />
      </div>

      <!-- Empty State -->
      <div v-if="filteredRooms.length === 0" class="text-center py-12">
        <div class="text-6xl mb-4">🏢</div>
        <h3 class="text-xl font-medium text-slate-300 mb-2">
          Tidak ada ruangan ditemukan
        </h3>
        <p class="text-slate-400">Coba ubah filter atau kata kunci pencarian</p>
      </div>
    </main>

    <!-- Booking Modal -->
    <div
      v-if="isModalOpen"
      class="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center p-4 z-50"
      @click="closeBookingModal"
    >
      <div
        @click.stop
        class="bg-slate-800/95 backdrop-blur-sm rounded-xl shadow-2xl shadow-blue-950/40 border border-slate-700/50 w-full max-w-2xl max-h-[90vh] overflow-y-auto transform transition-all duration-300"
      >
        <div
          class="flex items-center justify-between p-6 border-b border-slate-700/50"
        >
          <div>
            <h2
              class="text-2xl font-bold bg-gradient-to-r from-blue-300 to-slate-200 bg-clip-text text-transparent"
            >
              📋 Booking {{ selectedRoom?.name }}
            </h2>
            <p class="text-slate-400 text-sm mt-1">
              Kapasitas: {{ selectedRoom?.capacity }} Orang
            </p>
          </div>
          <button
            @click="closeBookingModal"
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

        <!-- Modal Content -->
        <div class="p-6 space-y-6">
          <!-- Room Details -->
          <div>
            <h3 class="text-lg font-semibold text-slate-200 mb-3">
              Detail Ruangan
            </h3>
            <div class="bg-slate-700/50 rounded-lg p-4">
              <p class="text-slate-300 mb-3">{{ selectedRoom?.description }}</p>
              <div class="flex flex-wrap gap-2">
                <span
                  v-for="feature in selectedRoom?.features"
                  :key="feature"
                  class="px-3 py-1 bg-blue-600/20 text-blue-300 text-sm rounded-full border border-blue-600/30"
                >
                  {{ feature }}
                </span>
              </div>
            </div>
          </div>

          <!-- Date Selection -->
          <div>
            <h3 class="text-lg font-semibold text-slate-200 mb-3">
              📅 Pilih Tanggal
            </h3>
            <input
              v-model="bookingForm.date"
              type="date"
              :min="new Date().toISOString().split('T')[0]"
              class="w-full px-4 py-3 bg-slate-700/60 text-slate-100 border-2 border-slate-600 rounded-lg transition-all duration-300 outline-none focus:border-blue-500 focus:bg-slate-700/80 focus:ring-4 focus:ring-blue-500/20"
              required
            />
          </div>

          <!-- Time Selection -->
          <div class="grid grid-cols-2 gap-4">
            <div>
              <h3 class="text-lg font-semibold text-slate-200 mb-3">
                🕐 Waktu Mulai
              </h3>
              <input
                v-model="bookingForm.startTime"
                type="time"
                class="w-full px-4 py-3 bg-slate-700/60 text-slate-100 border-2 border-slate-600 rounded-lg transition-all duration-300 outline-none focus:border-blue-500 focus:bg-slate-700/80 focus:ring-4 focus:ring-blue-500/20"
                required
              />
            </div>
            <div>
              <h3 class="text-lg font-semibold text-slate-200 mb-3">
                🕐 Waktu Selesai
              </h3>
              <input
                v-model="bookingForm.endTime"
                type="time"
                class="w-full px-4 py-3 bg-slate-700/60 text-slate-100 border-2 border-slate-600 rounded-lg transition-all duration-300 outline-none focus:border-blue-500 focus:bg-slate-700/80 focus:ring-4 focus:ring-blue-500/20"
                required
              />
            </div>
          </div>

          <!-- Purpose -->
          <div>
            <h3 class="text-lg font-semibold text-slate-200 mb-3">
              📝 Tujuan Booking
            </h3>
            <textarea
              v-model="bookingForm.purpose"
              placeholder="Jelaskan tujuan penggunaan ruangan..."
              class="w-full px-4 py-3 bg-slate-700/60 text-slate-100 placeholder-slate-400 border-2 border-slate-600 rounded-lg transition-all duration-300 outline-none focus:border-blue-500 focus:bg-slate-700/80 focus:ring-4 focus:ring-blue-500/20 h-24 resize-none"
              required
            ></textarea>
          </div>
        </div>

        <!-- Modal Footer -->
        <div
          class="flex items-center justify-end space-x-3 p-6 border-t border-slate-700/50"
        >
          <button
            @click="closeBookingModal"
            class="px-6 py-3 text-slate-300 hover:text-slate-100 transition-colors duration-300"
          >
            Batal
          </button>
          <button
            @click="handleBooking"
            :disabled="isBookingLoading || !isBookingFormValid"
            class="px-8 py-3 bg-gradient-to-r from-blue-600 via-blue-700 to-slate-700 text-white font-semibold rounded-lg shadow-lg transform transition-all duration-300 hover:scale-105 hover:from-blue-500 hover:via-blue-600 hover:to-slate-600 hover:shadow-2xl hover:shadow-blue-600/30 active:scale-95 focus:outline-none focus:ring-4 focus:ring-blue-500/40 disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:scale-100"
          >
            <span v-if="!isBookingLoading">🎯 Booking Sekarang</span>
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
              Memproses...
            </span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
