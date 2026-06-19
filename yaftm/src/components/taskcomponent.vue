<script setup lang="ts">
  import type { Task } from '@/types/task.ts'
  import { useI18n } from 'vue-i18n'
  import { useTaskStore } from '@/stores/task.ts'
  import { useCategoryStore } from '@/stores/category.ts'


  const emit = defineEmits(['updatecategory'])

  const { t } = useI18n()
  const taskStore = useTaskStore()
  const categoryStore = useCategoryStore()

  const props = defineProps({
    task: {
      type: Object as PropType<Task>,
      required: true,
    },
  })

  const dueDateInfo = computed(() => {
    if (!props.task.duedate || props.task.iscompleted) return null
    const days = Math.ceil((new Date(props.task.duedate).getTime() - Date.now()) / 86_400_000)
    if (days < 0)   return { icon: 'mdi-fire',                color: '#e74c3c', text: t('task.overdue') }
    if (days === 0) return { icon: 'mdi-alarm',               color: '#e74c3c', text: t('task.dueToday') }
    if (days <= 3)  return { icon: 'mdi-clock-alert-outline', color: '#f39c12', text: t('task.daysLeft', days, { named: { n: days } }) }
    if (days <= 7)  return { icon: 'mdi-clock-outline',       color: '#f1c40f', text: t('task.daysLeft', days, { named: { n: days } }) }
    return          { icon: 'mdi-calendar-check',             color: '#2ecc71', text: t('task.daysLeft', days, { named: { n: days } }) }
  })

  const triggerEditTask = ref(false)
  const showDatePicker = ref(false)
  const taskName = ref('')
  const taskDescription = ref('')
  const taskDueDate = ref<Date>()
  const taskCategoryId = ref<number>(0)
  const triggerDeleteTask = ref(false)

  const formatDate = (iso?: string): string => {
    if (!iso) return '—'
    return new Date(iso).toLocaleDateString(undefined, { year: 'numeric', month: 'short', day: 'numeric' })
  }

  const statusChips = computed(() => [
    { label: t('task.priority'),  active: props.task.priority,    color: '#e67e22', icon: 'mdi-alert-circle-outline' },
    { label: t('task.urgent'),    active: props.task.urgency,     color: '#e74c3c', icon: 'mdi-lightning-bolt' },
    { label: t('task.started'),   active: props.task.isstarted,   color: '#27ae60', icon: 'mdi-play-circle-outline' },
    { label: t('task.backlog'),   active: props.task.isbacklog,   color: '#7f8c8d', icon: 'mdi-tray-full' },
    { label: t('task.completed'), active: props.task.iscompleted, color: '#2980b9', icon: 'mdi-check-circle-outline' },
  ])

  const onEditTask = (task: Task): void => {
    taskName.value = task.name
    taskDescription.value = task.comment
    taskDueDate.value = new Date(task.duedate)
    taskCategoryId.value = task.categoryid
    showDatePicker.value = false
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
      isstarted:   false,
      isbacklog:   !task.iscompleted,
    }
    await taskStore.updateTask(task.ID, updated)
    emit('updatecategory')
  }

  const onStartedTask = async (task: Task): Promise<void> => {
    const updated: Task = {
      ...task,
      isstarted: !task.isstarted,
    }
    await taskStore.updateTask(task.ID, updated)
    emit('updatecategory')
  }

  const saveTask = async (): Promise<void> => {
    if (!taskName.value || !taskDueDate.value) return
    const updated: Task = {
      ...props.task,
      name:       taskName.value,
      comment:    taskDescription.value,
      duedate:    taskDueDate.value.toISOString(),
      categoryid: taskCategoryId.value,
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
  <v-card class="task-card" :class="{ 'task-started': props.task.isstarted }">
    <v-checkbox
      v-model="props.task.iscompleted"
      class="status-checkbox"
      @change="onCompletedTask(props.task)"
    />
    <v-tooltip v-if="props.task.comment" location="top" max-width="320">
      <template #activator="{ props: tooltipProps }">
        <div v-bind="tooltipProps" class="task-name">{{ props.task.name }}</div>
      </template>
      <span class="task-notes-tooltip">{{ props.task.comment }}</span>
    </v-tooltip>
    <div v-else class="task-name">{{ props.task.name }}</div>
    <div class="task-actions">
      <v-tooltip v-if="dueDateInfo" :text="dueDateInfo.text" location="top">
        <template #activator="{ props: tooltipProps }">
          <v-icon v-bind="tooltipProps" :icon="dueDateInfo.icon" :color="dueDateInfo.color" size="small" />
        </template>
      </v-tooltip>
      <v-tooltip :text="props.task.isstarted ? t('task.markNotStarted') : t('task.markStarted')" location="top">
        <template #activator="{ props: tooltipProps }">
          <v-btn
            v-bind="tooltipProps"
            class="start-btn"
            :icon="props.task.isstarted ? 'mdi-play-circle' : 'mdi-play-circle-outline'"
            :color="props.task.isstarted ? '#27ae60' : undefined"
            size="small"
            @click="onStartedTask(props.task)"
          />
        </template>
      </v-tooltip>
      <v-btn class="edit-btn" icon="mdi-pencil" size="small" @click="onEditTask(props.task)" />
      <v-btn class="delete-btn" icon="mdi-trash-can" size="small" @click="onDeleteTask(props.task)" />
    </div>
  </v-card>
  <v-dialog v-model="triggerEditTask" max-width="720px" persistent scrollable>
    <v-card rounded="lg" class="edit-dialog-card">

      <v-card-title class="dialog-header">
        <v-icon icon="mdi-pencil-outline" size="small" class="mr-2" />
        {{ t('task.editTitle') }}
        <v-chip size="x-small" variant="tonal" color="grey" class="ml-2">#{{ props.task.ID }}</v-chip>
      </v-card-title>

      <v-divider />

      <v-card-text class="dialog-body">

        <v-text-field
          v-model="taskName"
          :label="t('task.name')"
          variant="outlined"
          density="comfortable"
          prepend-inner-icon="mdi-format-title"
          class="mb-3"
        />

        <div class="notes-label">
          <v-icon icon="mdi-note-text-outline" size="small" class="mr-1" />
          {{ t('task.notes') }}
        </div>
        <v-textarea
          v-model="taskDescription"
          variant="outlined"
          auto-grow
          rows="14"
          :placeholder="t('task.notesPlaceholder')"
          class="notes-textarea"
          hide-details
        />

        <v-select
          v-model="taskCategoryId"
          :items="categoryStore.categories"
          item-title="name"
          item-value="ID"
          :label="t('task.category')"
          variant="outlined"
          density="comfortable"
          prepend-inner-icon="mdi-folder-outline"
          class="mb-3"
        />

        <div class="due-date-row mt-3">
          <v-btn
            variant="tonal"
            size="small"
            :color="showDatePicker ? 'primary' : 'default'"
            prepend-icon="mdi-calendar-outline"
            @click="showDatePicker = !showDatePicker"
          >
            {{ taskDueDate ? formatDate(taskDueDate.toISOString()) : t('task.setDueDate') }}
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

        <v-expansion-panels variant="accordion" class="mt-4 properties-panel">
          <v-expansion-panel>
            <v-expansion-panel-title class="properties-title">
              <v-icon icon="mdi-information-outline" size="small" class="mr-2" />
              {{ t('task.properties') }}
            </v-expansion-panel-title>
            <v-expansion-panel-text>
              <div class="status-chips">
                <v-chip
                  v-for="chip in statusChips"
                  :key="chip.label"
                  size="small"
                  :style="chip.active ? `background:${chip.color};color:#fff` : ''"
                  :variant="chip.active ? 'flat' : 'tonal'"
                >
                  <v-icon :icon="chip.icon" size="x-small" start />
                  {{ chip.label }}
                </v-chip>
              </div>
              <div class="meta-dates">
                <span><strong>{{ t('task.created') }}</strong> {{ formatDate(props.task.CreatedAt) }}</span>
                <span><strong>{{ t('task.updated') }}</strong> {{ formatDate(props.task.UpdatedAt) }}</span>
              </div>
            </v-expansion-panel-text>
          </v-expansion-panel>
        </v-expansion-panels>

      </v-card-text>

      <v-divider />

      <v-card-actions class="dialog-actions">
        <v-btn variant="text" @click="triggerEditTask = false">{{ t('task.cancel') }}</v-btn>
        <v-btn variant="flat" color="primary" @click="saveTask()">{{ t('task.save') }}</v-btn>
      </v-card-actions>

    </v-card>
  </v-dialog>
  <v-dialog v-model="triggerDeleteTask" max-width="600px" persistent>
    <v-card>
      <v-card-title>{{ t('task.deleteTitle') }}</v-card-title>
      <v-card-text>
        {{ t('task.deleteConfirmIntro') }}
        <div>{{ t('task.name') }}: {{ taskName }}</div>
        <div>{{ t('task.notes') }}: {{ taskDescription }}</div>
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn @click="deleteTask()">{{ t('task.yes') }}</v-btn>
        <v-btn @click="triggerDeleteTask = false">{{ t('task.no') }}</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<style scoped lang="sass">
// ── Task card ─────────────────────────────────────────────────────────────────
.task-card
  width: 100%
  height: 100%
  margin: 0
  padding: 0
  border: none
  box-shadow: none
  display: flex
  align-items: center

.task-card:hover
  transform: translateY(-3px)
  box-shadow: 0 6px 15px rgba(0, 0, 0, 0.15)

.task-started
  border-left: 3px solid #27ae60

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

.edit-btn, .delete-btn, .start-btn
  background: none
  border: none
  cursor: pointer
  font-size: 18px
  transition: color 0.3s ease

.edit-btn:hover
  color: #3498db

.delete-btn:hover
  color: #e74c3c

// ── Edit dialog ───────────────────────────────────────────────────────────────
.edit-dialog-card
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

.properties-panel
  border-radius: 8px
  overflow: hidden
  flex-shrink: 0

.properties-title
  font-size: 13px
  font-weight: 600
  min-height: 40px !important
  color: #555

.status-chips
  display: flex
  flex-wrap: wrap
  gap: 8px
  padding: 4px 0 10px

.meta-dates
  display: flex
  gap: 20px
  font-size: 12px
  color: #777
  padding-top: 4px

.dialog-actions
  padding: 12px 16px
  justify-content: flex-end
  gap: 8px
  flex-shrink: 0

.task-notes-tooltip
  white-space: pre-wrap
  font-size: 12px
  line-height: 1.5
</style>
