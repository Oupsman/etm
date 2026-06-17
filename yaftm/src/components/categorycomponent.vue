<script setup lang="ts">
  import { defineProps, onMounted, ref } from 'vue'
  import { useTaskStore } from '@/stores/task.ts'

  import type { NewTask, Task } from '@/types/task.ts'
  import { VueDraggableNext as draggable } from 'vue-draggable-next'
  import TaskComponent from '@/components/taskcomponent.vue'
  import { useSnackbarStore } from '@/stores/snackbar';

  interface DropEvent {
    added?: { element: Task; newIndex: number }
    removed?: { element: Task; oldIndex: number }
    moved?: { element: Task; newIndex: number; oldIndex: number }
  }

  // Maps each draggable list name to the flags it represents
  const FLAG_MAP: Record<string, Pick<Task, 'isbacklog' | 'iscompleted' | 'urgency' | 'priority'>> = {
    backlog:               { isbacklog: true,  iscompleted: false, urgency: false, priority: false },
    completedTasks:        { isbacklog: false, iscompleted: true,  urgency: false, priority: false },
    urgentImportant:       { isbacklog: false, iscompleted: false, urgency: true,  priority: true  },
    nonUrgentImportant:    { isbacklog: false, iscompleted: false, urgency: false, priority: true  },
    urgentNonImportant:    { isbacklog: false, iscompleted: false, urgency: true,  priority: false },
    nonUrgentNonImportant: { isbacklog: false, iscompleted: false, urgency: false, priority: false },
  }

  const props = defineProps({
    categoryID: {
      type: Number,
      required: true,
    },
    readonly: {
      type: Boolean,
      default: false,
    },
  })
  const backlog = ref<Task[]>([])
  const urgentImportant = ref<Task[]>([])
  const nonUrgentImportant = ref<Task[]>([])
  const nonUrgentNonImportant = ref<Task[]>([])
  const urgentNonImportant = ref<Task[]>([])
  const completedTasks = ref<Task[]>([])

  const taskDialog = ref(false)
  const taskName = ref<string>('')
  const taskDescription = ref<string>('')
  const taskDueDate = ref<Date>(new Date())
  const showDatePicker = ref(false)

  const message = ref<string>('')
  const displaySnack = ref(false)

  const taskStore = useTaskStore()
  const snackbar = useSnackbarStore();

  const formatDate = (iso?: string): string => {
    if (!iso) return '—'
    return new Date(iso).toLocaleDateString(undefined, { year: 'numeric', month: 'short', day: 'numeric' })
  }

  const triggerTaskDialog = () => {
    taskName.value = ''
    taskDescription.value = ''
    taskDueDate.value = new Date()
    showDatePicker.value = false
    taskDialog.value = true
  }

  const addTask = () => {
    console.log('addTask')
    taskDialog.value = false
    if (taskName.value && taskDueDate.value) {
      const newTask: NewTask = {
        name: taskName.value,
        comment: taskDescription.value,
        duedate: taskDueDate.value.toISOString(),
        categoryid: props.categoryID,
        isbacklog: true,
      }
      const task: Task = {
        ID: 0,
        iscompleted: false,
        isstarted: false,
        urgency: false,
        priority: false,

        ...newTask,
      }
      taskStore.addTask(task).then(responseTask => {
        console.log('Response from server: ', responseTask)
        backlog.value.push(responseTask)

      }).catch (error => {
        console.log(error)
      })
    }
  }

  // Called once after each drop. Only cross-list moves (evt.added) change flags.
  const onDrop = (evt: DropEvent, destination: string) => {
    if (!evt.added) return
    const updated: Task = { ...evt.added.element, ...FLAG_MAP[destination] }
    taskStore.updateTask(updated.ID, updated)
  }

  const parseTasks = async () => {
    // query tasks from the store then parse tasks and add them to the respective lists
    backlog.value = []
    completedTasks.value = []
    urgentImportant.value = []
    nonUrgentImportant.value = []
    nonUrgentNonImportant.value = []
    urgentNonImportant.value = []
    await taskStore.getTasks(props.categoryID).then(tasks => {
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
    }
    ).catch(error => {
      snackbar.showSnackbar({
        message: 'Unable to get the task list ' + error.message,
        color: 'error',
      });
    })
  }

  onMounted(async () => {
    parseTasks()
  })
</script>

<template>
  <!-- Single root so that style="height:100%" from the parent is applied here -->
  <div style="height: 100%; display: flex; overflow: hidden;">

    <!-- Backlog Column -->
    <div class="backlog" style="flex: 0 0 25%; display: flex; flex-direction: column; overflow: hidden;">
      <h2>Backlog</h2>
      <v-btn v-if="!props.readonly" @click="triggerTaskDialog">Add task</v-btn>
      <v-chip v-else size="small" color="info" variant="tonal" class="mb-1">Read-only</v-chip>
      <div style="flex: 1; min-height: 0; overflow-y: auto;">
        <draggable
          v-model="backlog"
          group="tasks"
          itemkey="backlog"
          :disabled="props.readonly"
          @change="(e) => onDrop(e, 'backlog')"
        >
          <v-card v-for="task in backlog" :key="task.ID" class="mb-2 task">
            <TaskComponent :task="task" @updatecategory="parseTasks" />
          </v-card>
        </draggable>
      </div>
    </div>

    <!-- Eisenhower Matrix: 2×2 grid that fills remaining width and full height -->
    <div style="flex: 1; min-width: 0; display: grid; grid-template-columns: 1fr 1fr; grid-template-rows: 1fr 1fr; overflow: hidden;">
      <div class="UrgentImportant" style="display: flex; flex-direction: column; overflow: hidden;">
        <h2>Urgent et Important</h2>
        <div style="flex: 1; min-height: 0; overflow-y: auto;">
          <draggable v-model="urgentImportant" group="tasks" itemkey="urgentImportant" :disabled="props.readonly" @change="(e) => onDrop(e, 'urgentImportant')">
            <v-card v-for="task in urgentImportant" :key="task.ID" class="mb-2 task">
              <TaskComponent :task="task" @updatecategory="parseTasks" />
            </v-card>
          </draggable>
        </div>
      </div>
      <div class="NotUrgentImportant" style="display: flex; flex-direction: column; overflow: hidden;">
        <h2>Non Urgent et Important</h2>
        <div style="flex: 1; min-height: 0; overflow-y: auto;">
          <draggable v-model="nonUrgentImportant" group="tasks" itemkey="nonUrgentImportant" :disabled="props.readonly" @change="(e) => onDrop(e, 'nonUrgentImportant')">
            <v-card v-for="task in nonUrgentImportant" :key="task.ID" class="mb-2 task">
              <TaskComponent :task="task" @updatecategory="parseTasks" />
            </v-card>
          </draggable>
        </div>
      </div>
      <div class="UrgentNotImportant" style="display: flex; flex-direction: column; overflow: hidden;">
        <h2>Urgent et Non Important</h2>
        <div style="flex: 1; min-height: 0; overflow-y: auto;">
          <draggable v-model="urgentNonImportant" group="tasks" itemkey="urgentNonImportant" :disabled="props.readonly" @change="(e) => onDrop(e, 'urgentNonImportant')">
            <v-card v-for="task in urgentNonImportant" :key="task.ID" class="mb-2 task">
              <TaskComponent :task="task" @updatecategory="parseTasks" />
            </v-card>
          </draggable>
        </div>
      </div>
      <div class="NotUrgentNotImportant" style="display: flex; flex-direction: column; overflow: hidden;">
        <h2>Non Urgent et Non Important</h2>
        <div style="flex: 1; min-height: 0; overflow-y: auto;">
          <draggable v-model="nonUrgentNonImportant" group="tasks" itemkey="nonUrgentNonImportant" :disabled="props.readonly" @change="(e) => onDrop(e, 'nonUrgentNonImportant')">
            <v-card v-for="task in nonUrgentNonImportant" :key="task.ID" class="mb-2 task">
              <TaskComponent :task="task" @updatecategory="parseTasks" />
            </v-card>
          </draggable>
        </div>
      </div>
    </div>

    <!-- Completed Tasks Column -->
    <div class="completed" style="flex: 0 0 25%; display: flex; flex-direction: column; overflow: hidden;">
      <h3>Tâches Terminées</h3>
      <div style="flex: 1; min-height: 0; overflow-y: auto;">
        <draggable
          v-model="completedTasks"
          group="tasks"
          itemkey="completedTasks"
          :disabled="props.readonly"
          @change="(e) => onDrop(e, 'completedTasks')"
        >
          <v-card v-for="task in completedTasks" :key="task.ID" class="mb-2 task">
            <TaskComponent :task="task" @updatecategory="parseTasks" />
          </v-card>
        </draggable>
      </div>
    </div>

    <!-- Dialogs are teleported by Vuetify; placing them inside the root is fine -->
    <v-dialog v-model="taskDialog" max-width="720px" persistent scrollable>
      <v-card rounded="lg" class="task-dialog-card">

        <v-card-title class="dialog-header">
          <v-icon icon="mdi-plus-circle-outline" size="small" class="mr-2" />
          New task
        </v-card-title>

        <v-divider />

        <v-card-text class="dialog-body">

          <v-text-field
            v-model="taskName"
            label="Name"
            variant="outlined"
            density="comfortable"
            prepend-inner-icon="mdi-format-title"
            class="mb-3"
          />

          <div class="notes-label">
            <v-icon icon="mdi-note-text-outline" size="small" class="mr-1" />
            Notes
          </div>
          <v-textarea
            v-model="taskDescription"
            variant="outlined"
            auto-grow
            rows="14"
            placeholder="Add notes…"
            class="notes-textarea"
            hide-details
          />

          <div class="due-date-row mt-3">
            <v-btn
              variant="tonal"
              size="small"
              :color="showDatePicker ? 'primary' : 'default'"
              prepend-icon="mdi-calendar-outline"
              @click="showDatePicker = !showDatePicker"
            >
              {{ taskDueDate ? formatDate(taskDueDate.toISOString()) : 'Set due date' }}
            </v-btn>
          </div>

          <v-expand-transition>
            <div v-if="showDatePicker" class="date-picker-wrapper mt-2">
              <v-date-picker
                v-model="taskDueDate"
                elevation="0"
                border
                @update:model-value="showDatePicker = false"
              />
            </div>
          </v-expand-transition>

        </v-card-text>

        <v-divider />

        <v-card-actions class="dialog-actions">
          <v-btn variant="text" @click="taskDialog = false">Cancel</v-btn>
          <v-btn variant="flat" color="primary" @click="addTask">Add</v-btn>
        </v-card-actions>

      </v-card>
    </v-dialog>
    <v-snackbar v-model="displaySnack" timeout="3000">{{ message }}</v-snackbar>

  </div>
</template>

<style scoped lang="sass">
@import url('https://fonts.googleapis.com/css2?family=Poppins:wght@300;400;600&display=swap')

// ── Task creation dialog (matches edit dialog in taskcomponent.vue) ────────────
.task-dialog-card
  display: flex
  flex-direction: column
  max-height: 90vh

.dialog-header
  display: flex
  align-items: center
  padding: 16px 20px
  font-size: 16px
  font-weight: 600
  flex-shrink: 0

.dialog-body
  padding: 20px 24px
  display: flex
  flex-direction: column
  flex: 1
  overflow-y: auto

.notes-label
  font-size: 12px
  font-weight: 600
  color: #666
  text-transform: uppercase
  letter-spacing: 0.05em
  display: flex
  align-items: center
  margin-bottom: 6px

.notes-textarea
  font-family: 'Poppins', sans-serif
  font-size: 14px
  line-height: 1.6
  flex: 1

.due-date-row
  display: flex
  align-items: center

.date-picker-wrapper
  display: flex
  justify-content: flex-start

.dialog-actions
  padding: 12px 16px
  justify-content: flex-end
  gap: 8px
  flex-shrink: 0

.fill-height
  height: 100%

.backlog
  background: linear-gradient(135deg, #bdc3c7, #eef2f7)
  border-radius: 8px
  padding: 15px
  box-shadow: 0 5px 10px rgba(0, 0, 0, 0.1)

.completed
  background: linear-gradient(135deg, #2ecc71, #7ed56f)
  border-radius: 8px
  padding: 15px
  box-shadow: 0 5px 10px rgba(0, 0, 0, 0.1)

.UrgentImportant
  background: linear-gradient(135deg, #e74c3c, #ff6b6b)
  border-radius: 8px
  padding: 15px
  transition: transform 0.3s ease, box-shadow 0.3s ease

.UrgentNotImportant
  background: linear-gradient(135deg, #f39c12, #ffba49)
  border-radius: 8px
  padding: 15px
  transition: transform 0.3s ease, box-shadow 0.3s ease

.NotUrgentNotImportant
  background: linear-gradient(135deg, #3498db, #6ab0f3)
  border-radius: 8px
  padding: 15px
  transition: transform 0.3s ease, box-shadow 0.3s ease

.NotUrgentImportant
  background: linear-gradient(135deg, #9b59b6, #bb6bd9)
  border-radius: 8px
  padding: 15px
  transition: transform 0.3s ease, box-shadow 0.3s ease

.task
  width: 100%
  height: 50px
  min-height: 50px
  margin: 0
  padding: 0
  border: none
  box-shadow: none
  display: flex
  align-items: center
  justify-content: center

.task:hover
  transform: translateY(-5px)
  box-shadow: 0 8px 15px rgba(0, 0, 0, 0.2)

</style>
