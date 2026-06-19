import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Task } from '@/types/task'

export const useDragStore = defineStore('drag', () => {
  const draggingTask = ref<Task | null>(null)
  const draggingFromCategoryId = ref<number>(0)
  const refreshKey = ref(0)

  const startDrag = (task: Task, fromCategoryId: number) => {
    draggingTask.value = task
    draggingFromCategoryId.value = fromCategoryId
  }

  const endDrag = () => {
    draggingTask.value = null
    draggingFromCategoryId.value = 0
  }

  const notifyMoved = () => { refreshKey.value++ }

  return { draggingTask, draggingFromCategoryId, refreshKey, startDrag, endDrag, notifyMoved }
})
