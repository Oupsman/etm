<script setup lang="ts">
  import { useCategoryStore } from '@/stores/category.ts'
  import { useTaskStore} from '@/stores/task.ts'
  import { useAppStore } from '@/stores/app.ts'
  import { useUserStore} from '@/stores/user.ts'
  import type {Category} from "@/types/category.ts";

  const backlog = ref<Category[]>([])
  const urgentImportant = ref<Category[]>([])
  const nonUrgentImportant = ref<Category[]>([])
  const nonUrgentNonImportant = ref<Category[]>([])
  const urgentNonImportant = ref<Category[]>([])
  const completedTasks = ref<Category[]>([])

  const props = defineProps({
    category: {
      type: Object,
      required: true,
    },
  })

  const categoryStore = useCategoryStore()
  const taskStore = useTaskStore()



</script>
<template>
  <v-container class="fill-height" style="height: 90vh;">
    <v-row class="fill-height">
      <!-- Backlog Column -->
      <v-col cols="3" class="d-flex flex-column" style="position: absolute; left: 0; height:80vh;">
        <v-card class="backlog fill-height">
          <v-card-title>Backlog</v-card-title>
          <v-card-text>
            <v-list v-if="backlog.length > 0">
              <v-list-item v-for="(task, index) in backlog" :key="index">
                {{ task.name }}
              </v-list-item>
            </v-list>
          </v-card-text>
        </v-card>
      </v-col>

      <!-- Eisenhower Matrix -->
      <v-col cols="6" class="mx-auto fill-height d-flex flex-column " >
        <v-row class="d-flex">
          <v-col class="mx-auto d-flex flex-column " cols="6" style="position: relative; left: 0; height: 40vh;">
            <v-card class="UrgentImportant fill-height">
              <v-card-title>Urgent et Important</v-card-title>
              <v-card-text>
                <v-list v-if="urgentImportant.length > 0">
                  <v-list-item v-for="(task, index) in urgentImportant" :key="index">
                    {{ task.name }}
                  </v-list-item>
                </v-list>
              </v-card-text>
            </v-card>
          </v-col>
          <v-col cols="6" class="mx-auto d-flex flex-column " style="position: relative; left: 0; height: 40vh;">
            <v-card class="NotUrgentImportant fill-height">
              <v-card-title>Non Urgent et Important</v-card-title>
              <v-card-text>
                <v-list v-if="nonUrgentImportant.length > 0">
                  <v-list-item v-for="(task, index) in nonUrgentImportant" :key="index">
                    {{ task.name }}
                  </v-list-item>
                </v-list>
              </v-card-text>
            </v-card>
          </v-col>
          <v-col cols="6" class="mx-auto d-flex flex-column " style="position: relative; left: 0; height: 40vh;">
            <v-card class="UrgentNotImportant fill-height">
              <v-card-title>Urgent et Non Important</v-card-title>
              <v-card-text>
                <v-list v-if="urgentNonImportant.length > 0">
                  <v-list-item v-for="(task, index) in urgentNonImportant" :key="index">
                    {{ task.name }}
                  </v-list-item>
                </v-list>
              </v-card-text>
            </v-card>
          </v-col>
          <v-col cols="6" class="mx-auto d-flex flex-column " style="position: relative; left: 0; height: 40vh;">
            <v-card class="NotUrgentNotImportant fill-height">
              <v-card-title>Non Urgent et Non Important</v-card-title>
              <v-card-text >
                <v-list v-if="nonUrgentNonImportant.length > 0">
                  <v-list-item v-for="(task, index) in nonUrgentNonImportant" :key="index">
                    {{ task.name }}
                  </v-list-item>
                </v-list>
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>
      </v-col>

      <!-- Completed Tasks Column -->
      <v-col cols="3" class="d-flex flex-column" style="position: absolute; right: 0; height: 80vh;">
        <v-card class="completed fill-height">
          <v-card-title>Tâches Terminées</v-card-title>
          <v-card-text >
            <v-list v-if="completedTasks.length > 0">
              <v-list-item v-for="(task, index) in completedTasks" :key="index">
                {{ task.name }}
              </v-list-item>
            </v-list>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<style scoped lang="sass">
  .fill-height
    height: 100%

  .backlog
    background-color: lightgrey

  .completed
    background-color: lightgreen

  .UrgentImportant
    background-color: lightsalmon

  .UrgentNotImportant
    background-color: lightyellow

  .NotUrgentNotImportant
    background-color: lightskyblue

  .NotUrgentImportant
    background-color: lightpink

</style>
