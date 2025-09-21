import { ref } from "vue";
import { defineStore } from "pinia";
import type { Room, RoomForm } from "../types"; 
import apiClient from "../api"; 

export const useRoomStore = defineStore("room", () => {
  const rooms = ref<Room[]>([]);
  const isLoading = ref<boolean>(false);
  const isSubmitting = ref<boolean>(false);
  const fetchRooms = async () => {
    isLoading.value = true;
    try {
      const response = await apiClient.get<Room[]>("/rooms");
      rooms.value = response.data || [];
    } catch (error) {
      console.error("Gagal mengambil data ruangan:", error);
      alert("Gagal memuat data ruangan dari server.");
    } finally {
      isLoading.value = false;
    }
  };

  const createRoom = async (roomData: RoomForm) => {
    isSubmitting.value = true;
    try {
      await apiClient.post("/protected/rooms", roomData);
      await fetchRooms(); 
    } catch (error) {
      console.error("Gagal membuat ruangan:", error);
      throw new Error("Gagal menyimpan ruangan baru."); 
    } finally {
      isSubmitting.value = false;
    }
  };

  const updateRoom = async (id: number, roomData: RoomForm) => {
    isSubmitting.value = true;
    try {
      await apiClient.put(`/protected/rooms/${id}`, roomData);
      await fetchRooms(); 
    } catch (error) {
      console.error("Gagal memperbarui ruangan:", error);
      throw new Error("Gagal memperbarui ruangan.");
    } finally {
      isSubmitting.value = false;
    }
  };

  const deleteRoom = async (id: number) => {
    try {
      await apiClient.delete(`/protected/rooms/${id}`);
      await fetchRooms(); 
    } catch (error) {
      console.error("Gagal menghapus ruangan:", error);
      throw new Error("Gagal menghapus ruangan.");
    }
  };

  return {
    rooms,
    isLoading,
    isSubmitting,
    fetchRooms,
    createRoom,
    updateRoom,
    deleteRoom,
  };
});
