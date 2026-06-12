<script setup lang="ts">
  import type { Task } from '@/types/task.ts'
  import { useTaskStore } from '@/stores/task.ts'


  const emit = defineEmits(['updatecategory'])

  const taskStore = useTaskStore()

  const props = defineProps({
    task: {
      type: Object as PropType<Task>,
      required: true,
    },
  })

  const dueDateStyle = computed(() => {
    if (!props.task.duedate || props.task.iscompleted) return {}
    const days = Math.ceil((new Date(props.task.duedate).getTime() - Date.now()) / 86_400_000)
    if (days < 0)  return { borderLeft: '4px solid #e74c3c' }  // overdue
    if (days <= 3) return { borderLeft: '4px solid #f39c12' }  // due within 3 days
    if (days <= 7) return { borderLeft: '4px solid #f1c40f' }  // due this week
    return {}
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

  const onCompletedTask = async (task: Task): Promise<void> => {
    const updated: Task = {
      ...task,
      // v-model already flipped iscompleted; derive the rest from that new value
      priority:    false,
      urgency:     false,
      isbacklog:   !task.iscompleted,
    }
    await taskStore.updateTask(task.ID, updated)
    emit('updatecategory')
  }

  const saveTask = async (): Promise<void> => {
    if (!taskName.value || !taskDueDate.value) return
    const updated: Task = {
      ...props.task,
      name:    taskName.value,
      comment: taskDescription.value,
      duedate: taskDueDate.value.toISOString(),
    }
    await taskStore.updateTask(props.task.ID, updated)
    triggerEditTask.value = false
    emit('updatecategory')
  }

  const deleteTask = async (): Promise<void> => {
    await taskStore.deleteTask(props.task)
    triggerDeleteTask.value = false
    emit('updatecategory')
  }
</script>

<template>
  <v-card class="task-card" :style="dueDateStyle">
    <v-checkbox
      v-model="props.task.iscompleted"
      class="status-checkbox"
      @change="onCompletedTask(props.task)"
    />
    <div class="task-name">{{ props.task.name }}</div>
    <div class="task-actions">
      <v-btn class="edit-btn" icon="mdi-pencil" size="small" @click="onEditTask(props.task)" />
      <v-btn class="delete-btn" icon="mdi-trash-can" size="small" @click="onDeleteTask(props.task)" />
    </div>
  </v-card>
  <v-dialog v-model="triggerEditTask" max-width="600px" persistent>
    <v-card>
      <v-card-title>
        <span class="headline">Edit task</span>
      </v-card-title>
      <v-card-text>
        <v-container>
          <v-row>
            <v-col cols="12">
              <v-text-field
                v-model="taskName"
                label="Name"
                required
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="12">
              <v-text-field
                v-model="taskDescription"
                label="Description"
                required
              />
            </v-col>
          </v-row>

          <v-row>
            <v-col cols="12">
              <v-date-picker
                v-model="taskDueDate"
                label="Due Date"
                required
              />
            </v-col>
          </v-row>
        </v-container>
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn @click="saveTask()">Save</v-btn>
        <v-btn @click="triggerEditTask = false">Cancel</v-btn>
      </v-card-actions>

    </v-card>
  </v-dialog>
  <v-dialog v-model="triggerDeleteTask" max-width="600px" persistent>
    <v-card>
      <v-card-title>Are you sure ?</v-card-title>
      <v-card-text>Do you really want to delete this task ?
        Name: {{ taskName }}
        Description: {{ taskDescription }}</v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn @click="deleteTask()">YES</v-btn>
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
  font-size: 12px
  font-weight: 600
  color: #333
  flex-grow: 1
  text-align: center
  display: flex
  align-items: center
  justify-content: center


.task-actions
  display: flex
  gap: 10px
  justify-content: flex-end

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
