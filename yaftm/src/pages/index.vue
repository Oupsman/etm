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
  <v-container class="margin-initial" style="position: relative; top: 0; left: 0; height: 100vh;">
    <Categories />
  </v-container>
</template>

<style scoped lang="sass">
  margin-initial
    margin: initial
</style>
