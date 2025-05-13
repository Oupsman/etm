<script setup lang="ts">
  import { useCategoryStore } from '@/stores/category.ts'
  import { useTaskStore} from '@/stores/task.ts'
  import { useAppStore } from '@/stores/app.ts'
  import { useUserStore} from '@/stores/user.ts'
  import type {Category} from '@/types/category.ts'
  import type {Task} from '@/types/task.ts'

  const backlog = ref<Task[]>([])
  const urgentImportant = ref<Task[]>([])
  const nonUrgentImportant = ref<Task[]>([])
  const nonUrgentNonImportant = ref<Task[]>([])
  const urgentNonImportant = ref<Task[]>([])
  const completedTasks = ref<Task[]>([])

  const taskDialog = ref(false)
  const taskName = ref<String>('')
  const taskDescription = ref<String>('')
  const taskDueDate = ref<Date>(new Date())


  const props = defineProps({
    category: {
      type: Object,
      required: true,
    },
  })

  const categoryStore = useCategoryStore()
  const taskStore = useTaskStore()

  const triggerTaskDialog = () => {
    taskDialog.value = true
  }

  const addTask = () => {
    taskDialog.value = false
    console.log(typeof(taskDueDate.value))
    if (taskName.value && taskDescription.value && taskDueDate.value) {
      const newTask: Task = {
        name: taskName.value,
        description: taskDescription.value,
        dueDate: taskDueDate.value.toISOString(),
        categoryid: props.category.ID,
        isbackLog: true
      }
      taskStore.addTask(newTask)
      backlog.value.push(newTask)
    }
  }

  onMounted(async () => {
    console.log("Category :" + JSON.stringify(props.category))
    // query tasks from the store
    const tasks = await taskStore.getTasks(props.category.ID)
    console.log("Tasks:", tasks)
    // Parse tasks and add them to the respective lists
    tasks.forEach((task: Task) => {
      if (task.isbacklog) {
        backlog.value.push(task)
      } else if (task.iscompleted) {
        completedTasks.value.push(task)
      } else if (task.urgency && task.priority) {
        urgentImportant.value.push(task)
      } else if (!task.urgency && task.priority) {
        nonUrgentImportant.value.push(task)
      } else if (task.urgency && !task.priority) {
        urgentNonImportant.value.push(task)
      } else {
        nonUrgentNonImportant.value.push(task)
      }
    })
  })
</script>
<template>
  <v-container class="fill-height" style="height: 90vh;">
    <v-row class="fill-height">
      <!-- Backlog Column -->
      <v-col cols="3" class="d-flex flex-column" style="position: absolute; left: 0; height:80vh;">
        <v-card class="backlog fill-height">
          <v-card-title>Backlog</v-card-title>
          <v-btn @click="triggerTaskDialog">Add task</v-btn>
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
    <v-dialog v-model="taskDialog" persistent max-width="600px">
      <v-card>
        <v-card-title>
          <span class="headline">Add a new task</span>
        </v-card-title>
        <v-card-text>
          <v-container>
            <v-row>
              <v-col cols="12">
                <v-text-field
                  label="Name"
                  v-model="taskName"
                  required
                ></v-text-field>
              </v-col>
            </v-row>
            <v-row>
              <v-col cols="12">
                <v-text-field
                  label="Description"
                  v-model="taskDescription"
                  required
                ></v-text-field>
              </v-col>
            </v-row>

            <v-row>
              <v-col cols="12">
                <v-date-picker
                  label="Due Date"
                  v-model="taskDueDate"
                  required
                ></v-date-picker>
              </v-col>
            </v-row>
          </v-container>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="blue darken-1" text @click="taskDialog = false">
            Cancel
          </v-btn>
          <v-btn color="blue darken-1" text @click="addTask">
            Add
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
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
