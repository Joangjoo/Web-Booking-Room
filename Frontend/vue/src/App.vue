<script setup lang="ts">
import { RouterView } from 'vue-router'
import ToastContainer from '../src/components/ToastContainer.vue' 
</script>

<template>
  <div>
    <RouterView />
    <ToastContainer /> </div>
</template>