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
      <v-row>
        <!-- Backlog Column -->
        <v-col cols="3">
          <v-card>
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
        <v-col cols="6">
          <v-row>
            <v-col cols="6">
              <v-card>
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
              <v-card>
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
              <v-card>
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
              <v-card>
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
          <v-card>
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

</style>
