<script setup lang="ts">
  import { defineProps, onMounted, ref, watch } from 'vue'
  import { useI18n } from 'vue-i18n'
  import { useTaskStore } from '@/stores/task.ts'
  import { useDragStore } from '@/stores/drag'

  import type { Task } from '@/types/task.ts'
  import { VueDraggableNext as draggable } from 'vue-draggable-next'
  import TaskComponent from '@/components/taskcomponent.vue'
  import TaskFormDialog from '@/components/TaskFormDialog.vue'
  import { useSnackbarStore } from '@/stores/snackbar';

  interface DropEvent {
    added?: { element: Task; newIndex: number }
    removed?: { element: Task; oldIndex: number }
    moved?: { element: Task; newIndex: number; oldIndex: number }
  }

  // Maps each draggable list name to the flags it represents
  const FLAG_MAP: Record<string, Pick<Task, 'isbacklog' | 'iscompleted' | 'urgency' | 'priority'>> = {
    backlog:               { isbacklog: true, iscompleted: false, urgency: false, priority: false },
    completedTasks:        { isbacklog: false, iscompleted: true, urgency: false, priority: false },
    urgentImportant:       { isbacklog: false, iscompleted: false, urgency: true, priority: true },
    nonUrgentImportant:    { isbacklog: false, iscompleted: false, urgency: false, priority: true },
    urgentNonImportant:    { isbacklog: false, iscompleted: false, urgency: true, priority: false },
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

  const { t } = useI18n()
  const taskStore = useTaskStore()
  const dragStore = useDragStore()
  const snackbar = useSnackbarStore()

  const onTaskCreated = (task: Task) => {
    backlog.value.push(task)
  }

  const onListDragStart = (list: Task[], evt: { oldIndex: number }) => {
    const task = list[evt.oldIndex]
    if (task) dragStore.startDrag(task, props.categoryID)
  }

  const onDragEnd = () => dragStore.endDrag()

  // Re-fetch tasks when a cross-category move is signalled via the drag store.
  watch(() => dragStore.refreshKey, () => parseTasks())

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
      <h2>{{ t('matrix.backlog') }}</h2>
      <button v-if="!props.readonly" class="add-task-btn" @click="taskDialog = true">
        <v-icon class="mr-1" size="small">mdi-plus</v-icon>
        {{ t('matrix.addTask') }}
      </button>
      <v-chip
        v-else
        class="mb-1"
        color="info"
        size="small"
        variant="tonal"
      >{{ t('category.readOnly') }}</v-chip>
      <div style="flex: 1; min-height: 0; overflow-y: auto;">
        <draggable
          v-model="backlog"
          :disabled="props.readonly"
          group="tasks"
          itemkey="backlog"
          @change="(e) => onDrop(e, 'backlog')"
          @end="onDragEnd"
          @start="(e) => onListDragStart(backlog, e)"
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
        <h2>{{ t('matrix.urgentImportant') }}</h2>
        <div style="flex: 1; min-height: 0; overflow-y: auto;">
          <draggable
            v-model="urgentImportant"
            :disabled="props.readonly"
            group="tasks"
            itemkey="urgentImportant"
            @change="(e) => onDrop(e, 'urgentImportant')"
            @end="onDragEnd"
            @start="(e) => onListDragStart(urgentImportant, e)"
          >
            <v-card v-for="task in urgentImportant" :key="task.ID" class="mb-2 task">
              <TaskComponent :task="task" @updatecategory="parseTasks" />
            </v-card>
          </draggable>
        </div>
      </div>
      <div class="NotUrgentImportant" style="display: flex; flex-direction: column; overflow: hidden;">
        <h2>{{ t('matrix.notUrgentImportant') }}</h2>
        <div style="flex: 1; min-height: 0; overflow-y: auto;">
          <draggable
            v-model="nonUrgentImportant"
            :disabled="props.readonly"
            group="tasks"
            itemkey="nonUrgentImportant"
            @change="(e) => onDrop(e, 'nonUrgentImportant')"
            @end="onDragEnd"
            @start="(e) => onListDragStart(nonUrgentImportant, e)"
          >
            <v-card v-for="task in nonUrgentImportant" :key="task.ID" class="mb-2 task">
              <TaskComponent :task="task" @updatecategory="parseTasks" />
            </v-card>
          </draggable>
        </div>
      </div>
      <div class="UrgentNotImportant" style="display: flex; flex-direction: column; overflow: hidden;">
        <h2>{{ t('matrix.urgentNotImportant') }}</h2>
        <div style="flex: 1; min-height: 0; overflow-y: auto;">
          <draggable
            v-model="urgentNonImportant"
            :disabled="props.readonly"
            group="tasks"
            itemkey="urgentNonImportant"
            @change="(e) => onDrop(e, 'urgentNonImportant')"
            @end="onDragEnd"
            @start="(e) => onListDragStart(urgentNonImportant, e)"
          >
            <v-card v-for="task in urgentNonImportant" :key="task.ID" class="mb-2 task">
              <TaskComponent :task="task" @updatecategory="parseTasks" />
            </v-card>
          </draggable>
        </div>
      </div>
      <div class="NotUrgentNotImportant" style="display: flex; flex-direction: column; overflow: hidden;">
        <h2>{{ t('matrix.notUrgentNotImportant') }}</h2>
        <div style="flex: 1; min-height: 0; overflow-y: auto;">
          <draggable
            v-model="nonUrgentNonImportant"
            :disabled="props.readonly"
            group="tasks"
            itemkey="nonUrgentNonImportant"
            @change="(e) => onDrop(e, 'nonUrgentNonImportant')"
            @end="onDragEnd"
            @start="(e) => onListDragStart(nonUrgentNonImportant, e)"
          >
            <v-card v-for="task in nonUrgentNonImportant" :key="task.ID" class="mb-2 task">
              <TaskComponent :task="task" @updatecategory="parseTasks" />
            </v-card>
          </draggable>
        </div>
      </div>
    </div>

    <!-- Completed Tasks Column -->
    <div class="completed" style="flex: 0 0 25%; display: flex; flex-direction: column; overflow: hidden;">
      <h3>{{ t('matrix.completed') }}</h3>
      <div style="flex: 1; min-height: 0; overflow-y: auto;">
        <draggable
          v-model="completedTasks"
          :disabled="props.readonly"
          group="tasks"
          itemkey="completedTasks"
          @change="(e) => onDrop(e, 'completedTasks')"
          @end="onDragEnd"
          @start="(e) => onListDragStart(completedTasks, e)"
        >
          <v-card v-for="task in completedTasks" :key="task.ID" class="mb-2 task">
            <TaskComponent :task="task" @updatecategory="parseTasks" />
          </v-card>
        </draggable>
      </div>
    </div>

    <TaskFormDialog v-model="taskDialog" :category-id="props.categoryID" @saved="onTaskCreated" />

  </div>
</template>

<style scoped lang="sass">
@import url('https://fonts.googleapis.com/css2?family=Poppins:wght@300;400;600&display=swap')

.fill-height
  height: 100%

.add-task-btn
  width: 100%
  margin-bottom: 8px
  padding: 6px 0
  display: flex
  align-items: center
  justify-content: center
  background: transparent
  border: 1.5px dashed rgba(0, 0, 0, 0.25)
  border-radius: 6px
  font-size: 12px
  font-weight: 600
  color: rgba(0, 0, 0, 0.5)
  cursor: pointer
  transition: background 0.2s ease, color 0.2s ease, border-color 0.2s ease

  &:hover
    background: rgba(0, 0, 0, 0.06)
    border-color: rgba(0, 0, 0, 0.45)
    color: rgba(0, 0, 0, 0.75)

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
