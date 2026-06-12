<!-- App.vue -->
<template>
  <v-app>
    <v-snackbar
      v-model="snackbar.show"
      :color="snackbar.color"
      :timeout="snackbar.timeout"
      @update:modelValue="snackbar.closeSnackbar"
    >
      {{ snackbar.message }}
      <template v-slot:actions>
        <v-btn
          v-if="snackbar.actionText"
          variant="text"
          @click="snackbar.actionCallback"
        >
          {{ snackbar.actionText }}
        </v-btn>
        <v-btn variant="text" @click="snackbar.closeSnackbar">Close</v-btn>
      </template>
    </v-snackbar>
    <router-view />
  </v-app>
</template>

<script setup lang="ts">
  import { useSnackbarStore } from '@/stores/snackbar';
  import { useAuthStore } from '@/stores/auth';
  import { useAppStore } from '@/stores/app';

  const snackbar = useSnackbarStore();
  const authStore = useAuthStore();
  const appStore = useAppStore();
  if (!authStore.isAuthenticated && !globalThis.location.pathname.startsWith('/auth/')) {
    appStore.navigateToPage('/login')
  }
</script>
