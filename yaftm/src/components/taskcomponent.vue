<script setup lang="ts">
  import type { Task } from '@/types/task.ts'
  import { useTaskStore } from '@/stores/task.ts'

  const taskStore = useTaskStore()

  const props = defineProps({
    task: {
      type: Object as PropType<Task>,
      required: true
    },
  })

  const triggerEditTask = ref(false)
  const taskName = ref('')
  const taskDescription = ref('')
  const taskDueDate = ref<Date>()
  const triggerDeleteTask = ref(false)


  const onEditTask = (task: Task): void => {

    console.log('Edit task ', task)
    taskName.value = task.name
    taskDescription.value = task.comment
    taskDueDate.value = new Date(task.duedate)
    triggerEditTask.value = true
  }

  const onDeleteTask = (task: Task): void => {
    taskName.value = task.name
    taskDescription.value = task.comment
    triggerDeleteTask.value = true
  }

  const saveTask = (task: Task): void => {
    console.log('Save task ', task)
    if (taskName.value && taskDescription.value && taskDueDate.value) {
      task.name = taskName.value
      task.comment = taskDescription.value
      task.duedate = taskDueDate.value.toISOString()
      if (taskStore.updateTask(task.ID, task)) {
        triggerEditTask.value = false
      }
    }
  }
  const deleteTask = (task: Task): void => {
    console.log('Delete task: ', task)
    if (taskStore.deleteTask(task)) {
      triggerDeleteTask.value = false
    }
  }
</script>

<template>
  <v-card class="mb-2 task" style="margin: 0;">
    <v-icon icon="mdi-checkbox-marked-outline" size="small"> </v-icon> {{ props.task.name }}
    <v-btn icon="mdi-pencil" @click="onEditTask(props.task)" size="small" density="compact" :right="true" :absolute="true"> </v-btn>
    <v-btn icon="mdi-trash-can" @click="onDeleteTask(props.task)" size="small" density="compact" :right="true" :absolute="true"> </v-btn>
  </v-card>
  <v-dialog v-model="triggerEditTask" persistent max-width="600px">
    <v-card>
      <v-card-title>
        <span class="headline">Edit task</span>
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
    <v-btn @click="saveTask(props.task)">Save</v-btn>
    <v-btn @click="triggerEditTask = false">Cancel</v-btn>
      </v-card-actions>

    </v-card>
  </v-dialog>
  <v-dialog v-model="triggerDeleteTask" persistent max-width="600px">
    <v-card>
      <v-card-title>Are you sure ?</v-card-title>
      <v-card-text>Do you really want to delete this task ?
        Name: {{ taskName }}
      Description: {{ taskDescription }}</v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn @click="deleteTask(props.task)">YES</v-btn>
        <v-btn @click="triggerDeleteTask = false">NO</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<style scoped lang="sass">

</style>
