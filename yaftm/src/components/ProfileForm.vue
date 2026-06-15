<script setup lang="ts">
  import { useUserStore } from '@/stores/user'
  import { useAppStore } from '@/stores/app'
  import { useNotificationStore } from '@/stores/notification'
  import { useVuelidate } from '@vuelidate/core'
  import { minLength, required, sameAs } from '@vuelidate/validators'
  import type { User } from '@/types/user'
  const userStore = useUserStore()
  const appStore = useAppStore()
  const notifStore = useNotificationStore()

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

  const user = ref<User>()

  const fetchUser = async () => {
    try {
      isLoading.value = true
      user.value = await userStore.getUser()
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

  onMounted(async () => {
    await fetchUser()
    await notifStore.check()
  })

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

      <v-divider class="my-4" />

      <div class="text-subtitle-1 mb-2">Browser Notifications</div>
      <template v-if="!notifStore.isSupported">
        <v-alert type="warning" density="compact">
          Push notifications are not supported in this browser.
        </v-alert>
      </template>
      <template v-else-if="notifStore.permissionDenied">
        <v-alert type="error" density="compact">
          Notifications are blocked. Reset permissions in your browser settings and reload.
        </v-alert>
      </template>
      <template v-else>
        <v-switch
          :model-value="notifStore.isSubscribed"
          color="primary"
          :label="notifStore.isSubscribed ? 'Notifications enabled' : 'Notifications disabled'"
          hide-details
          @update:model-value="notifStore.toggle()"
        />
        <v-btn
          v-if="notifStore.isSubscribed"
          variant="outlined"
          size="small"
          class="mt-2"
          @click="notifStore.sendTest()"
        >
          Send test notification
        </v-btn>
      </template>
    </v-container>
  </div>
  <div v-else>User not found</div>
</template>
