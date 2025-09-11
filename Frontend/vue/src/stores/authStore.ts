import { defineStore } from "pinia";
import axios from "axios";

const API_URL = "http://localhost:8080/api";

interface UserState {
  token: string | null;
}

interface RegisterPayload {
  name: string;
  email: string;
  password: string;
}

export const useAuthStore = defineStore("auth", {
  state: (): UserState => ({
    token: localStorage.getItem("token") || null,
  }),

  getters: {
    isAuthenticated: (state) => !!state.token,
  },

  actions: {
    setToken(tokenString: string) {
      this.token = tokenString;
      localStorage.setItem("token", tokenString);
      axios.defaults.headers.common["Authorization"] = `Bearer ${tokenString}`;
    },

    logout() {
      this.token = null;
      localStorage.removeItem("token");
      delete axios.defaults.headers.common["Authorization"];
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
        console.error(
          "Login API error:",
          error.response?.data?.error || error.message
        );
        const errorMessage =
          error.response?.data?.error || "Terjadi kesalahan saat login.";
        throw new Error(errorMessage);
      }
    },

    async register(payload: RegisterPayload){
      try{
        await axios.post(`${API_URL}/register`,{
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
    }
  },
});
