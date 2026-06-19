<script setup lang="ts">
  import { ref } from 'vue'
  import type { Ref } from 'vue'
  import { useI18n } from 'vue-i18n'
  import { useUserStore } from '@/stores/user'

  const { t } = useI18n()
  const userStore = useUserStore()

  const username: Ref<string> = ref('')
  const password: Ref<string> = ref('')
  const email: Ref<string> = ref('')

  const signup = () => {
    userStore.signup(username.value, password.value, email.value)
  }
</script>

<template>
  <v-container>
    <v-card>
      <v-card-title>{{ t('auth.signup') }}</v-card-title>
      <v-card-text>
        <v-form @submit.prevent="signup">
          <v-text-field
            v-model="username"
            autofocus
            :label="t('auth.username')"
            type="text"
          />
          <v-text-field
            v-model="password"
            :label="t('auth.password')"
            type="password"
          />
          <v-text-field
            v-model="email"
            :label="t('auth.email')"
            type="email"
          />
          <v-btn color="primary" type="submit">{{ t('auth.signup') }}</v-btn>
        </v-form>
      </v-card-text>
    </v-card>
  </v-container>
</template>
