<template>
  <v-footer
    app
    height="40"
  >
    Yet Another F*cking Task Manager
  </v-footer>
</template>

<script setup lang="ts">

</script>

<style scoped lang="sass">
  .social-link :deep(.v-icon)
    color: rgba(var(--v-theme-on-background), var(--v-disabled-opacity))
    text-decoration: none
    transition: .2s ease-in-out

    &:hover
      color: rgba(25, 118, 210, 1)
</style>
