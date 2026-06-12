<script lang="ts" setup>
  import { useUserStore } from '@/stores/user'
  import { useAppStore } from '@/stores/app'

  const userStore = useUserStore()
  const appStore = useAppStore()
  if (!userStore.userIsLoggedIn) {
    // redirect to the login page
    appStore.navigateToPage('/login')
  }
</script>

<template>
  <div style="width: 100%; height: 100%; display: flex; flex-direction: column;">
    <Categories />
  </div>
</template>
