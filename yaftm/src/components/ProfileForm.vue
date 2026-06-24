<script setup lang="ts">
  import { useI18n } from 'vue-i18n'
  import { useUserStore } from '@/stores/user'
  import { useAppStore } from '@/stores/app'
  import { useNotificationStore } from '@/stores/notification'
  import { useAuthStore } from '@/stores/auth'
  import { useSnackbarStore } from '@/stores/snackbar'
  import { axiosInstance } from '@/plugins/axios'
  import { useVuelidate } from '@vuelidate/core'
  import { minLength, required, sameAs } from '@vuelidate/validators'
  import type { User } from '@/types/user'

  const { t } = useI18n()
  const userStore = useUserStore()
  const appStore = useAppStore()
  const notifStore = useNotificationStore()
  const authStore = useAuthStore()
  const snackbar = useSnackbarStore()

  const isLoading = ref(true)
  const showPassword = ref(false)
  const oidcEnabled = ref(false)

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

  const linkOIDC = () => {
    window.location.href = `/api/v1/auth/oidc/link?token=${authStore.token}`
  }

  const unlinkOIDC = async () => {
    try {
      await axiosInstance.delete('/api/v1/auth/oidc/link')
      snackbar.showSnackbar({ message: 'OIDC account unlinked.', color: 'success' })
      await fetchUser()
    } catch (e: any) {
      snackbar.showSnackbar({ message: 'Unlink failed: ' + e.message, color: 'error' })
    }
  }

  onMounted(async () => {
    await fetchUser()
    await notifStore.check()
    try {
      const { data } = await axiosInstance.get('/api/v1/auth/oidc/status')
      oidcEnabled.value = data.enabled === true
    } catch { /* OIDC status fetch is best-effort */ }
  })

  const save = async () => {
    const isValid = await v$.value.$validate()
    if (!isValid) return
    await userStore.updateUser({ ...userForm.value })
  }
</script>

<template>
  <div v-if="isLoading">{{ t('profile.loading') }}</div>
  <div v-else-if="user">
    <v-container>
      <v-form @submit.prevent="save">
        <v-text-field
          v-model="userForm.oldPassword"
          :append-icon="showPassword ? 'mdi-eye' : 'mdi-eye-off'"
          :error-messages="v$.oldPassword.$errors.map(e => unref(e.$message))"
          :label="t('profile.currentPassword')"
          prepend-icon="mdi-lock"
          required
          type="password"
        />
        <v-text-field
          v-model="userForm.newPassword"
          :append-icon="showPassword ? 'mdi-eye' : 'mdi-eye-off'"
          :error-messages="v$.newPassword.$errors.map(e => unref(e.$message))"
          :label="t('profile.newPassword')"
          prepend-icon="mdi-lock"
          type="password"
        />
        <v-text-field
          v-model="userForm.newPasswordConfirmation"
          :append-icon="showPassword ? 'mdi-eye' : 'mdi-eye-off'"
          :error-messages="v$.newPasswordConfirmation.$errors.map(e => unref(e.$message))"
          :label="t('profile.confirmPassword')"
          prepend-icon="mdi-lock"
          type="password"
          @click:append="showPassword = !showPassword"
        />

        <v-btn color="primary" :disabled="v$.$invalid" type="submit">{{ t('profile.save') }}</v-btn>
      </v-form>

      <v-divider class="my-4" />

      <div class="text-subtitle-1 mb-2">{{ t('profile.notifications') }}</div>
      <template v-if="!notifStore.isSupported">
        <v-alert density="compact" type="warning">
          {{ t('profile.notSupported') }}
        </v-alert>
      </template>
      <template v-else-if="notifStore.permissionDenied">
        <v-alert density="compact" type="error">
          {{ t('profile.blocked') }}
        </v-alert>
      </template>
      <template v-else>
        <v-switch
          color="primary"
          hide-details
          :label="notifStore.isSubscribed ? t('profile.enabled') : t('profile.disabled')"
          :model-value="notifStore.isSubscribed"
          @update:model-value="notifStore.toggle()"
        />
        <v-btn
          v-if="notifStore.isSubscribed"
          class="mt-2"
          size="small"
          variant="outlined"
          @click="notifStore.sendTest()"
        >
          {{ t('profile.sendTest') }}
        </v-btn>
      </template>

      <v-divider class="my-4" />
      <DevicesPanel />

      <template v-if="oidcEnabled">
        <v-divider class="my-4" />
        <div class="text-subtitle-1 mb-2">{{ t('profile.oidcTitle') }}</div>
        <template v-if="user?.oidc_subject">
          <v-alert class="mb-3" density="compact" type="success">
            {{ t('profile.linked', { provider: user.oidc_provider }) }}
          </v-alert>
          <v-btn color="error" size="small" variant="outlined" @click="unlinkOIDC">
            {{ t('profile.unlink') }}
          </v-btn>
        </template>
        <template v-else>
          <v-alert class="mb-3" density="compact" type="info">
            {{ t('profile.notLinked') }}
          </v-alert>
          <v-btn color="primary" size="small" variant="outlined" @click="linkOIDC">
            {{ t('profile.link') }}
          </v-btn>
        </template>
      </template>
    </v-container>
  </div>
  <div v-else>{{ t('profile.notFound') }}</div>
</template>
