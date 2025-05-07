<script setup lang="ts">
  import { useCategoryStore } from '@/stores/category.ts'
  import { useTaskStore} from '@/stores/task.ts'
  import { useAppStore } from '@/stores/app.ts'
  import { useUserStore} from '@/stores/user.ts'

  const backlog = ref()
  const urgentImportant = ref()
  const nonUrgentImportant = ref()
  const nonUrgentNonImportant = ref()
  const urgentNonImportant = ref()
  const completedTasks = ref()

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
    <v-container>
      <v-row class="flex-child text-subtitle-2" height="100%">
        <!-- Backlog Column -->
        <v-col class="d-flex" cols="1" width="25%" height="100%">
          <v-card class="d-flex" color="lightgrey" height="100%">
            <v-card-title>Backlog</v-card-title>
            <v-card-text>
              <v-list>
                <v-list-item v-for="(task, index) in backlog" :key="index">
                  {{ task.name }}
                </v-list-item>
              </v-list>
            </v-card-text>
          </v-card>
        </v-col>

        <!-- Eisenhower Matrix -->
        <v-col cols="2">
          <v-row>
            <v-col cols="6">
              <v-card class="UrgentImportant">
                <v-card-title>Urgent et Important</v-card-title>
                <v-card-text>
                  <v-list>
                    <v-list-item v-for="(task, index) in urgentImportant" :key="index">
                      {{ task.name }}
                    </v-list-item>
                  </v-list>
                </v-card-text>
              </v-card>
            </v-col>
            <v-col cols="6">
              <v-card class="NotUrgentImportant">
                <v-card-title>Non Urgent et Important</v-card-title>
                <v-card-text>
                  <v-list>
                    <v-list-item v-for="(task, index) in nonUrgentImportant" :key="index">
                      {{ task.name }}
                    </v-list-item>
                  </v-list>
                </v-card-text>
              </v-card>
            </v-col>
            <v-col cols="6">
              <v-card class="UrgentNotImportant">
                <v-card-title>Urgent et Non Important</v-card-title>
                <v-card-text>
                  <v-list>
                    <v-list-item v-for="(task, index) in urgentNonImportant" :key="index">
                      {{ task.name }}
                    </v-list-item>
                  </v-list>
                </v-card-text>
              </v-card>
            </v-col>
            <v-col cols="6">
              <v-card class="NotUrgentNotImportant">
                <v-card-title>Non Urgent et Non Important</v-card-title>
                <v-card-text>
                  <v-list>
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
        <v-col cols="3">
          <v-card class="completed">
            <v-card-title>Tâches Terminées</v-card-title>
            <v-card-text>
              <v-list>
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

  .backlog
    background-color: lightgrey
    width: 25%
    height: 100%
    float: left
    justify-content: center


  .completed
    background-color: lightgreen
    width: 25%
    height: 100%
    float: right
    justify-content: center


  .UrgentImportant
    background-color: lightsalmon
    width: 25%
    height: 50%
    position: relative
    float: left
    justify-content: center


  .UrgentNotImportant
    background-color: lightyellow
    width: 25%
    height: 50%
    position: relative
    float: left
    justify-content: center


  .NotUrgentNotImportant
    background-color: lightskyblue
    width: 25%
    height: 50%
    position: relative
    float: left
    justify-content: center


  .NotUrgentImportant
    background-color: lightpink
    width: 25%
    height: 50%
    position: relative
    float: left

</style>
