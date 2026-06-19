<script setup lang="ts">
  import { useI18n } from 'vue-i18n'
  import { useUserStore } from '@/stores/user'

  const { t } = useI18n()
  const userStore = useUserStore()
  const username: Ref<string> = ref('')
  const password: Ref<string> = ref('')

  const login = () => {
    userStore.login(username.value, password.value)
  }
</script>

<template>
  <v-container>
    <v-card>
      <v-card-title>{{ t('auth.login') }}</v-card-title>
      <v-card-text>
        <v-form @submit.prevent="login">
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
          <v-btn color="primary" type="submit">{{ t('auth.login') }}</v-btn>
        </v-form>
      </v-card-text>
    </v-card>
  </v-container>
</template>
