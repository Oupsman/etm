<script setup lang="ts">
  import type { Task } from '@/types/task.ts'
  import { useI18n } from 'vue-i18n'
  import { useTaskStore } from '@/stores/task.ts'
  import TaskFormDialog from '@/components/TaskFormDialog.vue'

  const emit = defineEmits(['updatecategory'])

  const { t } = useI18n()
  const taskStore = useTaskStore()

  const props = defineProps({
    task: {
      type: Object as PropType<Task>,
      required: true,
    },
  })

  const dueDateInfo = computed(() => {
    if (!props.task.duedate || props.task.iscompleted) return null
    const days = Math.ceil((new Date(props.task.duedate).getTime() - Date.now()) / 86_400_000)
    if (days < 0) return { icon: 'mdi-fire', color: '#e74c3c', text: t('task.overdue') }
    if (days === 0) return { icon: 'mdi-alarm', color: '#e74c3c', text: t('task.dueToday') }
    if (days <= 3) return { icon: 'mdi-clock-alert-outline', color: '#f39c12', text: t('task.daysLeft', days, { named: { n: days } }) }
    if (days <= 7) return { icon: 'mdi-clock-outline', color: '#f1c40f', text: t('task.daysLeft', days, { named: { n: days } }) }
    return { icon: 'mdi-calendar-check', color: '#2ecc71', text: t('task.daysLeft', days, { named: { n: days } }) }
  })

  const showEditDialog = ref(false)
  const showDeleteDialog = ref(false)

  const onCompletedTask = async (task: Task, completed: boolean): Promise<void> => {
    const updated: Task = {
      ...task,
      iscompleted: completed,
      priority:    false,
      urgency:     false,
      isstarted:   false,
      isbacklog:   !completed,
    }
    await taskStore.updateTask(task.ID, updated)
    emit('updatecategory')
  }

  const onStartedTask = async (task: Task): Promise<void> => {
    await taskStore.updateTask(task.ID, { ...task, isstarted: !task.isstarted })
    emit('updatecategory')
  }

  const onSaved = (): void => {
    showEditDialog.value = false
    emit('updatecategory')
  }

  const deleteTask = async (): Promise<void> => {
    await taskStore.deleteTask(props.task)
    showDeleteDialog.value = false
    emit('updatecategory')
  }
</script>

<template>
  <v-card class="task-card" :class="{ 'task-started': props.task.isstarted }">
    <v-checkbox
      class="status-checkbox"
      :model-value="props.task.iscompleted"
      @update:model-value="(val) => onCompletedTask(props.task, val as boolean)"
    />
    <v-tooltip v-if="props.task.comment" location="top" max-width="320">
      <template #activator="{ props: tooltipProps }">
        <div v-bind="tooltipProps" class="task-name">{{ props.task.name }}</div>
      </template>
      <span class="task-notes-tooltip">{{ props.task.comment }}</span>
    </v-tooltip>
    <div v-else class="task-name">{{ props.task.name }}</div>
    <div class="task-actions">
      <v-tooltip v-if="props.task.link" location="top" :text="t('task.openLink')">
        <template #activator="{ props: tooltipProps }">
          <v-btn
            v-bind="tooltipProps"
            class="link-btn"
            :href="props.task.link"
            icon="mdi-link-variant"
            rel="noopener noreferrer"
            size="small"
            target="_blank"
          />
        </template>
      </v-tooltip>
      <v-tooltip v-if="dueDateInfo" location="top" :text="dueDateInfo.text">
        <template #activator="{ props: tooltipProps }">
          <v-icon v-bind="tooltipProps" :color="dueDateInfo.color" :icon="dueDateInfo.icon" size="small" />
        </template>
      </v-tooltip>
      <v-tooltip location="top" :text="props.task.isstarted ? t('task.markNotStarted') : t('task.markStarted')">
        <template #activator="{ props: tooltipProps }">
          <v-btn
            v-bind="tooltipProps"
            class="start-btn"
            :color="props.task.isstarted ? '#27ae60' : undefined"
            :icon="props.task.isstarted ? 'mdi-play-circle' : 'mdi-play-circle-outline'"
            size="small"
            @click="onStartedTask(props.task)"
          />
        </template>
      </v-tooltip>
      <v-btn class="edit-btn" icon="mdi-pencil" size="small" @click="showEditDialog = true" />
      <v-btn class="delete-btn" icon="mdi-trash-can" size="small" @click="showDeleteDialog = true" />
    </div>
  </v-card>

  <TaskFormDialog v-model="showEditDialog" :task="props.task" @saved="onSaved" />

  <v-dialog v-model="showDeleteDialog" max-width="600px" persistent>
    <v-card>
      <v-card-title>{{ t('task.deleteTitle') }}</v-card-title>
      <v-card-text>
        {{ t('task.deleteConfirmIntro') }}
        <div>{{ t('task.name') }}: {{ props.task.name }}</div>
        <div>{{ t('task.notes') }}: {{ props.task.comment }}</div>
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn @click="deleteTask">{{ t('task.yes') }}</v-btn>
        <v-btn @click="showDeleteDialog = false">{{ t('task.no') }}</v-btn>
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

.edit-btn, .delete-btn, .start-btn, .link-btn
  background: none
  border: none
  cursor: pointer
  font-size: 18px
  transition: color 0.3s ease

.edit-btn:hover
  color: #3498db

.delete-btn:hover
  color: #e74c3c

.task-notes-tooltip
  white-space: pre-wrap
  font-size: 12px
  line-height: 1.5
</style>
