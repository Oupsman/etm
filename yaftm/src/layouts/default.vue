<template>
  <v-app-bar
    v-model="showBar"
    @mouseenter="cancelHide"
    @mouseleave="scheduleHide"
  >
    <v-app-bar-title>
      <RouterLink to="/" style="text-decoration: none; color: inherit;">
        Eisenhower Matrix Task Manager
      </RouterLink>
    </v-app-bar-title>
    <v-spacer />
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
        <v-list-item :to="{ name: 'profile' }">
          <v-list-item-title>Profile</v-list-item-title>
        </v-list-item>
        <v-list-item :to="{ name: 'groups' }">
          <v-list-item-title>Groups</v-list-item-title>
        </v-list-item>
        <v-list-item @click="authStore.logout()">
          <v-list-item-title>Logout</v-list-item-title>
        </v-list-item>
      </v-list>
    </v-menu>
  </v-app-bar>
  <v-main>
    <router-view />
  </v-main>
</template>

<script lang="ts" setup>
  import { ref, onMounted, onUnmounted } from 'vue'
  import { useAuthStore } from '@/stores/auth'
  import { useUserStore } from '@/stores/user'
  import type { User } from '@/types/user'

  const authStore = useAuthStore()
  const userStore = useUserStore()
  const currentUser = ref<User | null>(null)

  onMounted(async () => {
    if (authStore.isAuthenticated) {
      try { currentUser.value = await userStore.getUser() } catch { /* ignore */ }
    }
  })
  const showBar = ref(false)
  let hideTimer: ReturnType<typeof setTimeout> | null = null

  const cancelHide = () => {
    if (hideTimer !== null) { clearTimeout(hideTimer); hideTimer = null }
  }

  const scheduleHide = () => {
    cancelHide()
    hideTimer = setTimeout(() => { showBar.value = false }, 500)
  }

  const onMouseMove = (e: MouseEvent) => {
    if (e.clientY < 8) {
      cancelHide()
      showBar.value = true
    }
  }

  onMounted(() => document.addEventListener('mousemove', onMouseMove))
  onUnmounted(() => document.removeEventListener('mousemove', onMouseMove))
</script>
