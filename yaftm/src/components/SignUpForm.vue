<script setup lang="ts">
  import { ref } from 'vue'
  import type { Ref } from 'vue'
  import { useUserStore } from '@/stores/user'

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
      <v-card-title>SignUp</v-card-title>
      <v-card-text>
        <v-form @submit.prevent="signup">
          <v-text-field
            v-model="username"
            autofocus
            label="Username"
            type="text"
          />
          <v-text-field
            v-model="password"
            label="Password"
            type="password"
          />
          <v-text-field
            v-model="email"
            label="Email"
            type="email"
          />
          <v-btn color="primary" type="submit">Sign Up</v-btn>
        </v-form>
      </v-card-text>
    </v-card>

  </v-container>
</template>

<style scoped lang="sass">

</style>
