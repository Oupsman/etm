<script setup lang="ts">
  import { useUserStore } from '@/stores/user'
  import { useAppStore } from '@/stores/app'
  import { useVuelidate } from '@vuelidate/core'
  import { minLength, required, sameAs } from '@vuelidate/validators'
  import type { User } from '@/types/user'
  const userStore = useUserStore()
  const appStore = useAppStore()

  const isLoading = ref(true)
  const showPassword = ref(false)

  const userForm = ref({
    oldPassword: '',
    newPassword: '',
    newPasswordConfirmation: '',
  })

  const rules = computed(() => ({
    oldPassword: { required },
    newPassword: { minLength: minLength(6) },
    newPasswordConfirmation: {
      sameAsPassword: sameAs(userForm.value.newPassword),
    },
  }))

  const v$ = useVuelidate(rules, userForm)

  const user = ref(<User>{})

  const fetchUser = async () => {
    try {
      isLoading.value = true
      await userStore.getUser()
      user.value = userStore.user
      if (user.value) {
        appStore.pageTitle = user.value.username
        Object.assign(userForm.value, user.value)
      }
    } catch (error) {
      console.error('Error while loading user', error)
    } finally {
      isLoading.value = false
    }
  }

  onMounted(fetchUser)

  const save = async () => {
    const isValid = await v$.value.$validate()
    if (!isValid) return

    const updatedUser = {
      ...userForm.value,
    }
    await userStore.updateUser(updatedUser)
  }
</script>

<template>
  <div v-if="isLoading">Loading...</div>
  <div v-else-if="user">
    <v-container>
      <v-form @submit.prevent="save">
        <v-text-field
          v-model="userForm.oldPassword"
          :append-icon="showPassword ? 'mdi-eye' : 'mdi-eye-off'"
          :error-messages="v$.oldPassword.$errors.map(e => unref(e.$message))"
          label="Current Password"
          prepend-icon="mdi-lock"
          required
          type="password"
        />
        <v-text-field
          v-model="userForm.newPassword"
          :append-icon="showPassword ? 'mdi-eye' : 'mdi-eye-off'"
          :error-messages="v$.newPassword.$errors.map(e => unref(e.$message))"
          label="New Password"
          prepend-icon="mdi-lock"
          type="password"
        />
        <v-text-field
          v-model="userForm.newPasswordConfirmation"
          :append-icon="showPassword ? 'mdi-eye' : 'mdi-eye-off'"
          :error-messages="v$.newPasswordConfirmation.$errors.map(e => unref(e.$message))"
          label="Confirm Password"
          prepend-icon="mdi-lock"
          type="password"
          @click:append="showPassword = !showPassword"
        />

        <v-btn color="primary" :disabled="v$.$invalid" type="submit">Save</v-btn>
      </v-form>
    </v-container>
  </div>
  <div v-else>User not found</div>
</template>
