<script setup lang="ts">
  import { useI18n } from 'vue-i18n'
  import { useDeviceStore } from '@/stores/device'
  import { useSnackbarStore } from '@/stores/snackbar'

  const { t } = useI18n()
  const deviceStore = useDeviceStore()
  const snackbar = useSnackbarStore()

  const registerDialog = ref(false)
  const keyDialog = ref(false)
  const revokeDialog = ref(false)
  const newDeviceName = ref('')
  const newApiKey = ref('')
  const copied = ref(false)
  const revokeTarget = ref<{ id: number; name: string } | null>(null)
  const registering = ref(false)

  onMounted(() => deviceStore.fetchDevices())

  const openRegister = () => {
    newDeviceName.value = ''
    registerDialog.value = true
  }

  const submitRegister = async () => {
    if (!newDeviceName.value.trim()) return
    registering.value = true
    try {
      const key = await deviceStore.registerDevice(newDeviceName.value.trim())
      newApiKey.value = key
      registerDialog.value = false
      keyDialog.value = true
    } catch (e) {
      snackbar.showSnackbar({ message: e instanceof Error ? e.message : 'Registration failed', color: 'error' })
    } finally {
      registering.value = false
    }
  }

  const copyKey = async () => {
    await navigator.clipboard.writeText(newApiKey.value)
    copied.value = true
    setTimeout(() => { copied.value = false }, 2000)
  }

  const confirmRevoke = (id: number, name: string) => {
    revokeTarget.value = { id, name }
    revokeDialog.value = true
  }

  const doRevoke = async () => {
    if (!revokeTarget.value) return
    try {
      await deviceStore.deleteDevice(revokeTarget.value.id)
      snackbar.showSnackbar({ message: 'Device revoked', color: 'success' })
    } catch (e) {
      snackbar.showSnackbar({ message: e instanceof Error ? e.message : 'Revoke failed', color: 'error' })
    } finally {
      revokeDialog.value = false
      revokeTarget.value = null
    }
  }

  const formatDate = (iso: string): string => {
    if (!iso || iso.startsWith('0001-')) return t('devices.never')
    return new Date(iso).toLocaleString()
  }
</script>

<template>
  <div class="text-subtitle-1 mb-2">{{ t('devices.title') }}</div>

  <v-btn
    class="mb-3"
    prepend-icon="mdi-plus"
    size="small"
    variant="outlined"
    @click="openRegister"
  >
    {{ t('devices.register') }}
  </v-btn>

  <template v-if="deviceStore.loading">
    <v-progress-linear class="mb-2" color="primary" indeterminate />
  </template>
  <template v-else-if="deviceStore.devices.length === 0">
    <p class="text-body-2 text-medium-emphasis">{{ t('devices.noDevices') }}</p>
  </template>
  <template v-else>
    <v-list density="compact" lines="two">
      <v-list-item
        v-for="device in deviceStore.devices"
        :key="device.ID"
        :subtitle="`${device.api_key_prefix}… · ${t('devices.lastUsed')}: ${formatDate(device.lastusedat)}`"
      >
        <template #title>
          <span class="font-weight-medium">{{ device.devicename }}</span>
        </template>
        <template #append>
          <v-btn
            color="error"
            icon="mdi-delete-outline"
            size="small"
            variant="text"
            @click="confirmRevoke(device.ID, device.devicename)"
          />
        </template>
      </v-list-item>
    </v-list>
  </template>

  <!-- Register dialog -->
  <v-dialog v-model="registerDialog" max-width="420">
    <v-card>
      <v-card-title>{{ t('devices.dialogTitle') }}</v-card-title>
      <v-card-text>
        <v-text-field
          v-model="newDeviceName"
          autofocus
          :label="t('devices.deviceName')"
          @keyup.enter="submitRegister"
        />
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn variant="text" @click="registerDialog = false">{{ t('devices.cancel') }}</v-btn>
        <v-btn
          color="primary"
          :disabled="!newDeviceName.trim() || registering"
          :loading="registering"
          variant="flat"
          @click="submitRegister"
        >
          {{ t('devices.submit') }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <!-- API key display dialog -->
  <v-dialog v-model="keyDialog" max-width="500" persistent>
    <v-card>
      <v-card-title>{{ t('devices.keyTitle') }}</v-card-title>
      <v-card-text>
        <v-alert class="mb-3" density="compact" type="warning">
          {{ t('devices.keyWarning') }}
        </v-alert>
        <v-text-field
          append-inner-icon="mdi-content-copy"
          density="compact"
          font-family="monospace"
          :model-value="newApiKey"
          readonly
          variant="outlined"
          @click:append-inner="copyKey"
        />
        <v-fade-transition>
          <p v-if="copied" class="text-success text-body-2 mt-1">{{ t('devices.copied') }}</p>
        </v-fade-transition>
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn color="primary" variant="flat" @click="keyDialog = false">{{ t('devices.close') }}</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <!-- Revoke confirmation dialog -->
  <v-dialog v-model="revokeDialog" max-width="400">
    <v-card>
      <v-card-title>{{ t('devices.revoke') }}</v-card-title>
      <v-card-text>
        {{ t('devices.revokeConfirm', { name: revokeTarget?.name ?? '' }) }}
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn variant="text" @click="revokeDialog = false">{{ t('devices.cancel') }}</v-btn>
        <v-btn color="error" variant="flat" @click="doRevoke">{{ t('devices.revoke') }}</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
