<script setup lang="ts">

  const hover = ref(false)
  const props = defineProps({
    task: {
      type: Object,
      required: true,
    },
  })

</script>

<template>
  <v-container class="ma-2 pa-1" style="height: 250px; width:300px;" @mouseleave="hover = false" @mouseover="hover = true">
    <v-card v-if="!hover" style="height: 250px; width:300px;">
      <v-card-title>
        {{ props.task.name }}
      </v-card-title>
    </v-card>
  </v-container>
</template>

<style scoped lang="sass">

</style>
