<script setup lang="ts">
  import { computed, ref, watch } from 'vue'
  import { useI18n } from 'vue-i18n'
  import { useTaskStore } from '@/stores/task'
  import { useCategoryStore } from '@/stores/category'
  import type { Task } from '@/types/task'

  const props = defineProps<{
    modelValue: boolean
    task?: Task
    categoryId?: number
  }>()

  const emit = defineEmits<{
    'update:model-value': [value: boolean]
    'saved': [task: Task]
    'deleted': []
  }>()

  const { t } = useI18n()
  const taskStore = useTaskStore()
  const categoryStore = useCategoryStore()

  const isEditMode = computed(() => !!props.task)

  const taskName = ref('')
  const taskDescription = ref('')
  const taskLink = ref('')
  const taskDueDate = ref<Date>(new Date())
  const taskCategoryId = ref<number>(0)
  const showDatePicker = ref(false)

  watch(() => props.modelValue, open => {
    if (!open) return
    showDatePicker.value = false
    if (props.task) {
      taskName.value = props.task.name
      taskDescription.value = props.task.comment
      taskLink.value = props.task.link ?? ''
      taskDueDate.value = new Date(props.task.duedate)
      taskCategoryId.value = props.task.categoryid
    } else {
      taskName.value = ''
      taskDescription.value = ''
      taskLink.value = ''
      taskDueDate.value = new Date()
      taskCategoryId.value = props.categoryId ?? 0
    }
  })

  const formatDate = (iso?: string): string => {
    if (!iso) return '—'
    return new Date(iso).toLocaleDateString(undefined, { year: 'numeric', month: 'short', day: 'numeric' })
  }

  const statusChips = computed(() => props.task ? [
    { label: t('task.priority'), active: props.task.priority, color: '#e67e22', icon: 'mdi-alert-circle-outline' },
    { label: t('task.urgent'), active: props.task.urgency, color: '#e74c3c', icon: 'mdi-lightning-bolt' },
    { label: t('task.started'), active: props.task.isstarted, color: '#27ae60', icon: 'mdi-play-circle-outline' },
    { label: t('task.backlog'), active: props.task.isbacklog, color: '#7f8c8d', icon: 'mdi-tray-full' },
    { label: t('task.completed'), active: props.task.iscompleted, color: '#2980b9', icon: 'mdi-check-circle-outline' },
  ] : [])

  const close = () => emit('update:model-value', false)

  const save = async () => {
    if (!taskName.value || !taskDueDate.value) return
    if (isEditMode.value && props.task) {
      const updated: Task = {
        ...props.task,
        name:       taskName.value,
        comment:    taskDescription.value,
        link:       taskLink.value,
        duedate:    taskDueDate.value.toISOString(),
        categoryid: taskCategoryId.value,
      }
      await taskStore.updateTask(props.task.ID, updated)
      emit('saved', updated)
    } else {
      const created: Task = {
        ID:          0,
        name:        taskName.value,
        comment:     taskDescription.value,
        link:        taskLink.value,
        duedate:     taskDueDate.value.toISOString(),
        categoryid:  props.categoryId ?? 0,
        isbacklog:   true,
        iscompleted: false,
        isstarted:   false,
        priority:    false,
        urgency:     false,
      }
      const saved = await taskStore.addTask(created)
      emit('saved', saved)
    }
    close()
  }
</script>

<template>
  <v-dialog
    max-width="720px"
    :model-value="modelValue"
    persistent
    scrollable
    @update:model-value="close"
  >
    <v-card class="form-dialog-card" rounded="lg">

      <v-card-title class="dialog-header">
        <v-icon class="mr-2" :icon="isEditMode ? 'mdi-pencil-outline' : 'mdi-plus-circle-outline'" size="small" />
        {{ isEditMode ? t('task.editTitle') : t('task.newTask') }}
        <v-chip
          v-if="isEditMode && task"
          class="ml-2"
          color="grey"
          size="x-small"
          variant="tonal"
        >
          #{{ task.ID }}
        </v-chip>
      </v-card-title>

      <v-divider />

      <v-card-text class="dialog-body">

        <v-text-field
          v-model="taskName"
          class="mb-3"
          density="comfortable"
          :label="t('task.name')"
          prepend-inner-icon="mdi-format-title"
          variant="outlined"
        />

        <v-text-field
          v-model="taskLink"
          class="mb-3"
          clearable
          density="comfortable"
          :label="t('task.link')"
          :placeholder="t('task.linkPlaceholder')"
          prepend-inner-icon="mdi-link-variant"
          variant="outlined"
        />

        <div class="notes-label">
          <v-icon class="mr-1" icon="mdi-note-text-outline" size="small" />
          {{ t('task.notes') }}
        </div>
        <v-textarea
          v-model="taskDescription"
          auto-grow
          class="notes-textarea"
          hide-details
          :placeholder="t('task.notesPlaceholder')"
          rows="14"
          variant="outlined"
        />

        <v-select
          v-if="isEditMode"
          v-model="taskCategoryId"
          class="mb-3 mt-3"
          density="comfortable"
          item-title="name"
          item-value="ID"
          :items="categoryStore.categories"
          :label="t('task.category')"
          prepend-inner-icon="mdi-folder-outline"
          variant="outlined"
        />

        <div class="due-date-row mt-3">
          <v-btn
            :color="showDatePicker ? 'primary' : 'default'"
            prepend-icon="mdi-calendar-outline"
            size="small"
            variant="tonal"
            @click="showDatePicker = !showDatePicker"
          >
            {{ taskDueDate ? formatDate(taskDueDate.toISOString()) : t('task.setDueDate') }}
          </v-btn>
        </div>

        <v-expand-transition>
          <div v-if="showDatePicker" class="date-picker-wrapper mt-2">
            <v-date-picker
              v-model="taskDueDate"
              border
              elevation="0"
              @update:model-value="showDatePicker = false"
            />
          </div>
        </v-expand-transition>

        <v-expansion-panels v-if="isEditMode && task" class="mt-4 properties-panel" variant="accordion">
          <v-expansion-panel>
            <v-expansion-panel-title class="properties-title">
              <v-icon class="mr-2" icon="mdi-information-outline" size="small" />
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
                <span><strong>{{ t('task.created') }}</strong> {{ formatDate(task.CreatedAt) }}</span>
                <span><strong>{{ t('task.updated') }}</strong> {{ formatDate(task.UpdatedAt) }}</span>
              </div>
            </v-expansion-panel-text>
          </v-expansion-panel>
        </v-expansion-panels>

      </v-card-text>

      <v-divider />

      <v-card-actions class="dialog-actions">
        <v-btn variant="text" @click="close">{{ t('task.cancel') }}</v-btn>
        <v-btn color="primary" variant="flat" @click="save">
          {{ isEditMode ? t('task.save') : t('task.add') }}
        </v-btn>
      </v-card-actions>

    </v-card>
  </v-dialog>
</template>

<style scoped lang="sass">
.form-dialog-card
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
</style>
