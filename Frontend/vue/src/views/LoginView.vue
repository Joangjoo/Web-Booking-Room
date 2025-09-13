<script setup lang="ts">
import { ref, reactive } from "vue";
import { useRouter } from "vue-router";
import { useAuthStore } from "../stores/authStore";
import type { LoginFormData } from "../types";



const formData = reactive<LoginFormData>({
  email: "",
  password: "",
});

const isEmailFocused = ref<boolean>(false);
const isPasswordFocused = ref<boolean>(false);
const isLoading = ref<boolean>(false);
const authStore = useAuthStore();
const router = useRouter();
const apiErrorMessage = ref<string>("");
const showPassword = ref<boolean>(false);

const setEmailFocused = (focused: boolean): void => {
  isEmailFocused.value = focused;
};

const setPasswordFocused = (focused: boolean): void => {
  isPasswordFocused.value = focused;
};

const handleLogin = async () => {
  apiErrorMessage.value = "";
  isLoading.value = true;

  try {
    await authStore.login(formData.email, formData.password);

    if (authStore.isAdmin) {
      router.push('/admin'); 
    } else {
      router.push('/dashboard'); 
    }
  } catch (error: any) {
    apiErrorMessage.value =
      error.message || "Terjadi kesalahan tidak diketahui.";
  } finally {
    isLoading.value = false;
  }
};

const handleForgotPassword = () => {
  alert("Fitur Lupa Password belum diimplementasikan.");
};
</script>

<template>
  <div
    class="min-h-screen flex items-center justify-center p-4 bg-gradient-to-br from-slate-900 via-blue-900 to-black"
  >
    <div
      class="absolute inset-0 bg-gradient-to-tr from-blue-950/40 via-slate-800/30 to-gray-900/50"
    ></div>
    <div class="relative z-10 w-full max-w-md">
      <div
        class="bg-slate-800/90 backdrop-blur-sm rounded-xl shadow-2xl shadow-blue-950/40 p-8 border border-slate-700/50 transform transition-all duration-300 hover:scale-[1.02] hover:shadow-3xl hover:shadow-blue-900/30"
      >
        <div class="text-center mb-8">
          <h1
            class="text-3xl font-bold bg-gradient-to-r from-blue-300 via-slate-200 to-blue-400 bg-clip-text text-transparent mb-3"
          >
            Selamat Datang 👋
          </h1>
          <p class="text-slate-400 text-sm">
            Masuk ke akun kamu untuk melanjutkan
          </p>
        </div>

        <div class="space-y-6">
          <div class="relative group">
            <label class="block text-sm font-medium text-slate-300 mb-2">
              Email
            </label>
            <div class="relative">
              <input
                v-model="formData.email"
                type="email"
                placeholder="Masukkan email kamu..."
                @focus="setEmailFocused(true)"
                @blur="setEmailFocused(false)"
                class="w-full px-4 py-3 bg-slate-700/60 text-slate-100 placeholder-slate-400 border-2 border-slate-600 rounded-lg transition-all duration-300 outline-none focus:border-blue-500 focus:bg-slate-700/80 focus:ring-4 focus:ring-blue-500/20 hover:border-slate-500"
                required
              />
              <div
                :class="[
                  'absolute inset-0 rounded-lg bg-gradient-to-r from-blue-500/5 via-slate-500/5 to-blue-600/5 pointer-events-none transition-opacity duration-300',
                  isEmailFocused ? 'opacity-100' : 'opacity-0',
                ]"
              ></div>
            </div>
          </div>

          <div class="relative group">
            <label class="block text-sm font-medium text-slate-300 mb-2">
              Password
            </label>
            <div class="relative">
              <input
                v-model="formData.password"
                :type="showPassword ? 'text' : 'password'"
                placeholder="Masukkan password kamu..."
                @focus="setPasswordFocused(true)"
                @blur="setPasswordFocused(false)"
                class="w-full px-4 py-3 bg-slate-700/60 text-slate-100 placeholder-slate-400 border-2 border-slate-600 rounded-lg transition-all duration-300 outline-none focus:border-blue-500 focus:bg-slate-700/80 focus:ring-4 focus:ring-blue-500/20 hover:border-slate-500"
                required
              />
              <button
                type="button"
                @click="showPassword = !showPassword"
                class="absolute inset-y-0 right-0 pr-3 flex items-center text-slate-400 hover:text-slate-200"
              >
                <svg
                  v-if="!showPassword"
                  xmlns="http://www.w3.org/2000/svg"
                  class="h-5 w-5"
                  viewBox="0 0 20 20"
                  fill="currentColor"
                >
                  <path d="M10 12a2 2 0 100-4 2 2 0 000 4z" />
                  <path
                    fill-rule="evenodd"
                    d="M.458 10C1.732 5.943 5.522 3 10 3s8.268 2.943 9.542 7c-1.274 4.057-5.022 7-9.542 7S1.732 14.057.458 10zM14 10a4 4 0 11-8 0 4 4 0 018 0z"
                    clip-rule="evenodd"
                  />
                </svg>
                <svg
                  v-else
                  xmlns="http://www.w3.org/2000/svg"
                  class="h-5 w-5"
                  viewBox="0 0 20 20"
                  fill="currentColor"
                >
                  <path
                    fill-rule="evenodd"
                    d="M3.707 2.293a1 1 0 00-1.414 1.414l14 14a1 1 0 001.414-1.414l-1.473-1.473A10.014 10.014 0 0019.542 10C18.268 5.943 14.478 3 10 3a9.958 9.958 0 00-4.512 1.074l-1.78-1.781zM10 12a2 2 0 100-4 2 2 0 000 4z"
                    clip-rule="evenodd"
                  />
                  <path
                    d="M2 10s3.939 7 8 7a9.958 9.958 0 004.512-1.074l-1.78-1.781a4 4 0 01-5.462-5.462L3.707 2.293A10.014 10.014 0 00.458 10z"
                  />
                </svg>
              </button>
              <div
                :class="[
                  'absolute inset-0 rounded-lg bg-gradient-to-r from-blue-500/5 via-slate-500/5 to-blue-600/5 pointer-events-none transition-opacity duration-300',
                  isPasswordFocused ? 'opacity-100' : 'opacity-0',
                ]"
              ></div>
            </div>
          </div>

          <button
            @click="handleLogin"
            :disabled="isLoading"
            class="w-full py-4 bg-gradient-to-r from-blue-600 via-blue-700 to-slate-700 text-white font-semibold rounded-lg shadow-lg transform transition-all duration-300 hover:scale-105 hover:from-blue-500 hover:via-blue-600 hover:to-slate-600 hover:shadow-2xl hover:shadow-blue-600/30 active:scale-95 focus:outline-none focus:ring-4 focus:ring-blue-500/40 disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:scale-100 disabled:hover:shadow-lg"
          >
            <span v-if="!isLoading" class="relative z-10">
              Masuk Sekarang
            </span>
            <span v-else class="relative z-10 flex items-center justify-center">
              <svg
                class="animate-spin -ml-1 mr-3 h-5 w-5 text-white"
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

          <div class="text-center">
            <a
              href="#"
              @click.prevent="handleForgotPassword"
              class="text-sm text-slate-400 hover:text-blue-400 transition-colors duration-300 hover:underline decoration-blue-400 focus:outline-none focus:ring-2 focus:ring-blue-500/40 focus:ring-offset-2 focus:ring-offset-slate-800 rounded px-1"
            >
              Lupa Password? 🤔
            </a>

            <router-link
              to="/register"
              class="block text-sm text-slate-400 hover:text-green-400 transition-colors duration-300 hover:underline decoration-green-400 focus:outline-none focus:ring-2 focus:ring-green-500/40 focus:ring-offset-2 focus:ring-offset-slate-800 rounded px-1"
            >
              Belum punya akun?
              <span class="font-medium text-green-500">Daftar di sini 🚀</span>
            </router-link>
          </div>
        </div>

        <!-- Subtle decorative elements -->
        <div
          class="absolute -top-3 -right-3 w-6 h-6 bg-gradient-to-r from-blue-600/20 to-slate-600/20 rounded-full blur-sm"
        ></div>
        <div
          class="absolute -bottom-3 -left-3 w-4 h-4 bg-gradient-to-r from-slate-600/20 to-blue-700/20 rounded-full blur-sm"
        ></div>
      </div>
    </div>

    <!-- Minimal floating elements -->
    <div
      class="absolute top-1/4 left-1/4 w-1 h-1 bg-blue-500/30 rounded-full animate-pulse"
    ></div>
    <div
      class="absolute top-3/4 right-1/4 w-2 h-2 bg-slate-500/20 rounded-full animate-pulse delay-1000"
    ></div>
    <div
      class="absolute bottom-1/3 left-1/2 w-1 h-1 bg-blue-600/25 rounded-full animate-pulse delay-500"
    ></div>
  </div>
</template>
