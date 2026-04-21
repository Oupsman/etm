<!-- App.vue -->
<template>
  <v-app>
    <v-menu>
      <template #activator="{ props }">
        <v-btn v-bind="props" icon>
          <v-icon>mdi-account</v-icon>
        </v-btn>
      </template>
      <v-list>
        <v-list-item>
          <v-list-item-title>Connected as {{ authStore.userID }}</v-list-item-title>
        </v-list-item>
        <v-list-item @click="authStore.logout()">
          <v-list-item-title>Logout</v-list-item-title>
        </v-list-item>
      </v-list>
    </v-menu>
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
        <v-btn
          variant="text"
          @click="snackbar.closeSnackbar"
        >
          Close
        </v-btn>
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
  if (!authStore.isAuthenticated) {
    appStore.navigateToPage('/login')
  }
</script>
