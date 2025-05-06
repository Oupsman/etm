<script setup lang="ts">
  const hover = ref(false)
  const props = defineProps({
    task: {
      type: Object,
      required: true,
    },
  })
</script>

<draggable>
  <template>

  </template>
</draggable>

<style scoped lang="sass">

</style>
