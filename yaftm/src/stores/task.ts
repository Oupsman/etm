import { ref } from 'vue'
import { defineStore } from 'pinia'
import { useSnackbarStore } from '@/stores/snackbar'
import { axiosInstance } from '@/plugins/axios'
import type { Task } from '@/types/task'

export const useTaskStore = defineStore('task', () => {
  const snackbar = useSnackbarStore()
  const tasks = ref([] as Task[])

  const addTask = async (task: Task): Promise<Task> => {
    try {
      const response = await axiosInstance.post('/api/v1/task', task)
      tasks.value.push(response.data.task)
      snackbar.showSnackbar({ message: 'Task added successfully.', color: 'success' })
      return response.data.task
    } catch (error: any) {
      snackbar.showSnackbar({ message: 'Unable to add task: ' + error.message, color: 'error' })
      throw error
    }
  }

  const deleteTask = async (taskToDelete: Task): Promise<boolean> => {
    tasks.value = tasks.value.filter(task => task.ID !== taskToDelete.ID)
    try {
      await axiosInstance.delete('/api/v1/task/' + taskToDelete.ID)
      snackbar.showSnackbar({ message: 'Task deleted successfully.', color: 'success' })
      return true
    } catch (error: any) {
      snackbar.showSnackbar({ message: 'Unable to delete task: ' + error.message, color: 'error' })
      throw error
    }
  }

  const updateTask = async (taskId: number, updatedTask: Task): Promise<boolean> => {
    try {
      await axiosInstance.post('/api/v1/task/' + taskId, updatedTask)
      const index = tasks.value.findIndex(task => task.ID === taskId)
      if (index !== -1) {
        tasks.value[index] = { ...tasks.value[index], ...updatedTask }
      }
      snackbar.showSnackbar({ message: 'Task updated successfully.', color: 'success' })
      return true
    } catch (error: any) {
      snackbar.showSnackbar({ message: 'Unable to update task: ' + error.message, color: 'error' })
      throw error
    }
  }

  const getTasks = async (categoryID: number): Promise<Task[]> => {
    const response = await axiosInstance.get('/api/v1/tasks/' + categoryID)
    tasks.value = response.data.tasks ?? []
    return tasks.value
  }

  return { tasks, addTask, deleteTask, updateTask, getTasks }
})
