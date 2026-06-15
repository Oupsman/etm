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

  const message = ref<string>('')
  const displaySnack = ref(false)

  const taskStore = useTaskStore()
  const snackbar = useSnackbarStore();

  const triggerTaskDialog = () => {
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
      <v-btn @click="triggerTaskDialog">Add task</v-btn>
      <div style="flex: 1; min-height: 0; overflow-y: auto;">
        <draggable
          v-model="backlog"
          group="tasks"
          itemkey="backlog"
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
          <draggable v-model="urgentImportant" group="tasks" itemkey="urgentImportant" @change="(e) => onDrop(e, 'urgentImportant')">
            <v-card v-for="task in urgentImportant" :key="task.ID" class="mb-2 task">
              <TaskComponent :task="task" @updatecategory="parseTasks" />
            </v-card>
          </draggable>
        </div>
      </div>
      <div class="NotUrgentImportant" style="display: flex; flex-direction: column; overflow: hidden;">
        <h2>Non Urgent et Important</h2>
        <div style="flex: 1; min-height: 0; overflow-y: auto;">
          <draggable v-model="nonUrgentImportant" group="tasks" itemkey="nonUrgentImportant" @change="(e) => onDrop(e, 'nonUrgentImportant')">
            <v-card v-for="task in nonUrgentImportant" :key="task.ID" class="mb-2 task">
              <TaskComponent :task="task" @updatecategory="parseTasks" />
            </v-card>
          </draggable>
        </div>
      </div>
      <div class="UrgentNotImportant" style="display: flex; flex-direction: column; overflow: hidden;">
        <h2>Urgent et Non Important</h2>
        <div style="flex: 1; min-height: 0; overflow-y: auto;">
          <draggable v-model="urgentNonImportant" group="tasks" itemkey="urgentNonImportant" @change="(e) => onDrop(e, 'urgentNonImportant')">
            <v-card v-for="task in urgentNonImportant" :key="task.ID" class="mb-2 task">
              <TaskComponent :task="task" @updatecategory="parseTasks" />
            </v-card>
          </draggable>
        </div>
      </div>
      <div class="NotUrgentNotImportant" style="display: flex; flex-direction: column; overflow: hidden;">
        <h2>Non Urgent et Non Important</h2>
        <div style="flex: 1; min-height: 0; overflow-y: auto;">
          <draggable v-model="nonUrgentNonImportant" group="tasks" itemkey="nonUrgentNonImportant" @change="(e) => onDrop(e, 'nonUrgentNonImportant')">
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
          @change="(e) => onDrop(e, 'completedTasks')"
        >
          <v-card v-for="task in completedTasks" :key="task.ID" class="mb-2 task">
            <TaskComponent :task="task" @updatecategory="parseTasks" />
          </v-card>
        </draggable>
      </div>
    </div>

    <!-- Dialogs are teleported by Vuetify; placing them inside the root is fine -->
    <v-dialog v-model="taskDialog" max-width="600px" persistent>
      <v-card>
        <v-card-title>
          <span class="headline">Add a new task</span>
        </v-card-title>
        <v-card-text>
          <v-container>
            <v-row>
              <v-col cols="12">
                <v-text-field v-model="taskName" label="Name" required />
              </v-col>
            </v-row>
            <v-row>
              <v-col cols="12">
                <v-text-field v-model="taskDescription" label="Description" required />
              </v-col>
            </v-row>
            <v-row>
              <v-col cols="12">
                <v-date-picker v-model="taskDueDate" label="Due Date" required />
              </v-col>
            </v-row>
          </v-container>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn color="blue darken-1" text @click="taskDialog = false">Cancel</v-btn>
          <v-btn color="blue darken-1" text @click="addTask">Add</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <v-snackbar v-model="displaySnack" timeout="3000">{{ message }}</v-snackbar>

  </div>
</template>

<style scoped lang="sass">
@import url('https://fonts.googleapis.com/css2?family=Poppins:wght@300;400;600&display=swap')

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
