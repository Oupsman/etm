<script lang="ts" setup>
  import { onMounted, ref } from 'vue'
  import axios from 'axios'
  import FormLogin from '@/components/LoginForm.vue'

  const oidcEnabled = ref(false)

  onMounted(async () => {
    try {
      const { data } = await axios.get('/api/v1/auth/oidc/status')
      oidcEnabled.value = data.enabled
    } catch {
      oidcEnabled.value = false
    }
  })

  function loginWithOIDC () {
    globalThis.location.href = '/api/v1/auth/oidc/login'
  }
</script>

<template>
  <div class="login">
    <h1>ETM Login</h1>
    <p>Welcome to the ETM login page. Please enter your username and password to login.</p>

    <form-login />

    <v-container v-if="oidcEnabled">
      <v-divider class="my-4" />
      <v-btn
        block
        color="secondary"
        prepend-icon="mdi-login-variant"
        @click="loginWithOIDC"
      >
        Login with OIDC
      </v-btn>
    </v-container>

    <p>No account yet ? Signup <router-link to="/signup">Here</router-link></p>
  </div>
</template>
