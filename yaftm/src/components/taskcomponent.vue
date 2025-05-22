<script setup lang="ts">
  import type { Task } from '@/types/task.ts'
  import { useTaskStore } from '@/stores/task.ts'


  const emit = defineEmits(['updatecategory'])

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

  const onCompletedTask = (task: Task): void => {
    if (task.iscompleted) {
      task.iscompleted = true
      task.priority = false
      task.urgency = false
      task.isbacklog = false
    } else {
      task.iscompleted = false
      task.priority = false
      task.urgency = false
      task.isbacklog = true
    }
    taskStore.updateTask(task.ID, task)
    emit('updatecategory')
  }


  const saveTask = (task: Task): void => {
    console.log('Save task ', task)
    if (taskName.value && taskDescription.value && taskDueDate.value) {
      task.name = taskName.value
      task.comment = taskDescription.value
      task.duedate = taskDueDate.value.toISOString()

      if ( taskStore.updateTask(task.ID, task)) {
        triggerEditTask.value = false
        emit('updatecategory')

        console.log('Event emitted')
      }
    }
  }
  const deleteTask = (task: Task): void => {
    if (taskStore.deleteTask(task)) {
      triggerDeleteTask.value = false
      emit('updatecategory')
    }
  }
</script>

<template>
    <v-card class="task-card">
      <v-checkbox
        v-model="props.task.iscompleted"
        class="status-checkbox"
        @change="onCompletedTask(props.task)"
      ></v-checkbox>
      <div class="task-name">{{ props.task.name }}</div>
      <div class="task-actions">
        <v-btn class="edit-btn" icon="mdi-pencil" @click="onEditTask(props.task)" size="small"></v-btn>
        <v-btn class="delete-btn" icon="mdi-trash-can" @click="onDeleteTask(props.task)" size="small"></v-btn>
      </div>
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
.task-card
  width: 100%
  height: 100%
  margin: 0
  padding: 0
  border: none
  box-shadow: none
  display: flex
  align-items: center
  justify-content: center

.task-card:hover
  transform: translateY(-3px)
  box-shadow: 0 6px 15px rgba(0, 0, 0, 0.15)

.status-checkbox
  display: flex
  margin-right: 10px

.task-name
  font-family: 'Poppins', sans-serif
  font-size: 16px
  font-weight: 600
  color: #333
  flex-grow: 1  // Permet au nom de prendre tout l'espace disponible
  text-align: center  // Centre le texte
  display: flex
  align-items: center
  justify-content: center


.task-actions
  display: flex
  gap: 10px
  justify-content: flex-end  // Alignement à droite

.edit-btn, .delete-btn
  background: none
  border: none
  cursor: pointer
  font-size: 18px
  transition: color 0.3s ease

.edit-btn:hover
  color: #3498db

.delete-btn:hover
  color: #e74c3c

</style>
