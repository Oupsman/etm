<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
onMounted(() => {
  const params = new URLSearchParams(window.location.hash.split('?')[1])
  const token = params.get('token')
  if (token) {
    localStorage.setItem('token', token) // use whatever key your app already uses
    router.push('/')
  } else {
    router.push('/login?error=oidc_failed')
  }
})
</script>
<template><div>Connexion en cours…</div></template>
