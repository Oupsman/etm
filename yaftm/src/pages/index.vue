<script lang="ts" setup>
  import { useCategoryStore} from '@/stores/category.ts'
</script>

<template>
  <v-container class="margin-initial" style="position: relative; top: 0; left: 0; height: 100vh;">
    <Categories />
  </v-container>
</template>

<style scoped lang="sass">
  margin-initial
    margin: initial
</style>
