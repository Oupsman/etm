<script setup lang="ts">
  import { onMounted } from 'vue'
  import { useRouter } from 'vue-router'
  import { useAuthStore } from '@/stores/auth'

  const router = useRouter()
  const authStore = useAuthStore()

  onMounted(() => {
    const params = new URLSearchParams(globalThis.location.search)
    const token = params.get('token')
    if (token) {
      authStore.setTokens(token, '')
      router.push('/')
    } else {
      router.push('/login?error=oidc_failed')
    }
  })
</script>
<template><div>Connexion en cours…</div></template>
