import { defineStore } from "pinia";
import axios from "axios";
import { jwtDecode } from "jwt-decode";
import router from "../router";

const API_URL = "http://localhost:8080/api";

interface UserState {
  token: string | null;
  user: UserData | null;
}

interface UserData {
  id: number;
  name: string;
  role: string;
}

interface RegisterPayload {
  name: string;
  email: string;
  password: string;
}

export const useAuthStore = defineStore("auth", {
  state: (): UserState => ({
    token: localStorage.getItem("token") || null,
    user: null,
  }),

  getters: {
    isAuthenticated: (state) => !!state.token,
    isAdmin: (state) => state.user?.role === "admin",
  },

  actions: {
    setToken(tokenString: string) {
      this.token = tokenString;
      localStorage.setItem("token", tokenString);
      axios.defaults.headers.common["Authorization"] = `Bearer ${tokenString}`;

      try {
        const decoded: any = jwtDecode(tokenString);
        this.user = {
          id: decoded.user_id,
          name: decoded.name,
          role: decoded.role,
        };
      } catch (error) {
        console.error("Gagal decode token:", error);
        this.logout(); // Jika token tidak valid, langsung logout
      }
    },

    logout() {
      this.token = null;
      this.user = null; // Pastikan data user juga dibersihkan
      localStorage.removeItem("token");
      delete axios.defaults.headers.common["Authorization"];
      router.push("/login"); // Arahkan ke halaman login
    },

    /**
     * Aksi Login: Menghubungi backend API untuk autentikasi.
     * @param email Email pengguna.
     * @param password Password pengguna.
     */
    async login(email: string, password: string) {
      try {
        const response = await axios.post(`${API_URL}/login`, {
          email: email,
          password: password,
        });
        if (response.data && response.data.token) {
          this.setToken(response.data.token); 
        } else {
          throw new Error("Respons token tidak ditemukan dari server");
        }
      } catch (error: any) {
        const errorMessage =
          error.response?.data?.error || "Terjadi kesalahan saat login.";
        throw new Error(errorMessage);
      }
    },

    async register(payload: RegisterPayload) {
      try {
        await axios.post(`${API_URL}/register`, {
          name: payload.name,
          email: payload.email,
          password: payload.password,
        });
      } catch (error: any) {
        console.error(
          "Register API error:",
          error.response?.data?.error || error.message
        );
        const errorMessage =
          error.response?.data?.error || "Terjadi kesalahan saat register.";
        throw new Error(errorMessage);
      }
    },
  },
});
