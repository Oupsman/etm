<template>
  <v-app-bar
    v-model="showBar"
    @mouseenter="cancelHide"
    @mouseleave="scheduleHide"
  >
    <v-app-bar-title>
      <RouterLink to="/" style="text-decoration: none; color: inherit;">
        {{ t('nav.title') }}
      </RouterLink>
    </v-app-bar-title>
    <v-spacer />

    <!-- Language picker -->
    <v-menu>
      <template #activator="{ props: menuProps }">
        <v-btn v-bind="menuProps" variant="text" size="small" class="mr-1">
          <v-icon start>mdi-translate</v-icon>
          {{ localeStore.locale.toUpperCase() }}
        </v-btn>
      </template>
      <v-list density="compact">
        <v-list-item
          v-for="lang in languages"
          :key="lang.code"
          :active="localeStore.locale === lang.code"
          active-color="primary"
          @click="localeStore.setLocale(lang.code)"
        >
          <v-list-item-title>{{ lang.label }}</v-list-item-title>
        </v-list-item>
      </v-list>
    </v-menu>

    <!-- User menu -->
    <v-menu>
      <template #activator="{ props }">
        <v-btn v-bind="props" icon>
          <v-icon>mdi-account</v-icon>
        </v-btn>
      </template>
      <v-list>
        <v-list-item>
          <v-list-item-title>{{ t('nav.connectedAs', { user: authStore.userID }) }}</v-list-item-title>
        </v-list-item>
        <v-list-item :to="{ name: 'profile' }">
          <v-list-item-title>{{ t('nav.profile') }}</v-list-item-title>
        </v-list-item>
        <v-list-item :to="{ name: 'groups' }">
          <v-list-item-title>{{ t('nav.groups') }}</v-list-item-title>
        </v-list-item>
        <v-list-item @click="authStore.logout()">
          <v-list-item-title>{{ t('nav.logout') }}</v-list-item-title>
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
  import { useI18n } from 'vue-i18n'
  import { useAuthStore } from '@/stores/auth'
  import { useUserStore } from '@/stores/user'
  import { useLocaleStore } from '@/stores/locale'
  import type { User } from '@/types/user'

  const { t } = useI18n()
  const authStore = useAuthStore()
  const userStore = useUserStore()
  const localeStore = useLocaleStore()
  const currentUser = ref<User | null>(null)

  const languages = [
    { code: 'en', label: 'English' },
    { code: 'fr', label: 'Français' },
  ]

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
