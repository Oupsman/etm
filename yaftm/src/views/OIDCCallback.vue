<script setup lang="ts">
  import { onMounted } from 'vue'
  import { useRouter } from 'vue-router'
  import { useAuthStore } from '@/stores/auth'
  import { useSnackbarStore } from '@/stores/snackbar'

  const router = useRouter()
  const authStore = useAuthStore()
  const snackbar = useSnackbarStore()

  const ERROR_MESSAGES: Record<string, string> = {
    oidc_already_linked: 'This OIDC identity is already linked to another account.',
    invalid_link_state: 'Link session expired. Please try again.',
    user_not_found: 'Account not found. Please try again.',
    link_failed: 'Failed to link account. Please try again.',
    oidc_failed: 'OIDC login failed.',
  }

  onMounted(() => {
    const params = new URLSearchParams(globalThis.location.search)
    const token = params.get('token')
    const linked = params.get('linked')
    const error = params.get('error')

    if (token) {
      authStore.setTokens(token, '')
      router.push('/')
    } else if (linked === 'true') {
      snackbar.showSnackbar({ message: 'OIDC account linked successfully.', color: 'success' })
      router.push('/profile')
    } else {
      const msg = ERROR_MESSAGES[error ?? ''] ?? 'Authentication failed.'
      snackbar.showSnackbar({ message: msg, color: 'error' })
      router.push('/login')
    }
  })
</script>
<template><div>Connexion en cours…</div></template>
