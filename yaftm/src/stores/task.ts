import { ref } from 'vue'
import type { Ref } from 'vue'
import { defineStore } from 'pinia'
import axios from 'axios'

import type { Task } from '@/types/task'


export const useTaskStore = defineStore('task', () => {
  const tasks = ref([] as Task[])
  const backlog = ref([] as Task[])

  const addTask = (task: Task): Task => {
    tasks.value.push(task)
    const token = localStorage.getItem('etm-token')
    if (!token) {
      throw new Error('No token')
    }
    const request = axios.create({
      baseURL: import.meta.env.VITE_BACKEND_URL,
      timeout: 1000,
      headers: { Authorization: `Bearer ${token}` },
    })
    request.post(import.meta.env.VITE_BACKEND_URL + '/api/v1/task', {
      ...task,
      token,
    }).then(response => {
      backlog.value.push(task)
      return response.data
    }).catch(error => {
      console.error('Create TaskCard error:', error)
      throw new Error('Create task failed')
    })
    return task
  }

  const removeTask = (taskToDelete: Task): boolean => {
    tasks.value = tasks.value.filter((task) => task.ID !== taskToDelete.ID)

    console.log('Delete task - function')
    confirm('Are you sure you want to delete this task?')
    const token = localStorage.getItem('etm-token')
    if (!token) {
      throw new Error('No token')
    }
    const request = axios.create({
      baseURL: import.meta.env.VITE_BACKEND_URL,
      timeout: 1000,
      headers: { Authorization: `Bearer ${token}` },
    })
    request.delete(import.meta.env.VITE_BACKEND_URL + '/api/v1/task/' + taskToDelete.ID).then(response => {
      console.log('task deleted:', response.data)
      return true
    }).catch(error => {
      console.error('Delete task error:', error)
      throw new Error('delete task')
    })
    return true
  }

  const updateTask = (taskId: number, updatedTask: Task): boolean => {
    console.log('Store UpdateTask')
    const index = tasks.value.findIndex((task) => task.ID === taskId)
    if (index !== -1) {
      tasks.value[index] = { ...tasks.value[index], ...updatedTask }
      console.log('TaskCard to save', updatedTask)
      const token = localStorage.getItem('etm-token')
      if (!token) {
        throw new Error('No token')
      }
      const request = axios.create({
        baseURL: import.meta.env.VITE_BACKEND_URL,
        timeout: 1000,
        headers: { Authorization: `Bearer ${token}` },
      })
      request.post(import.meta.env.VITE_BACKEND_URL + '/api/v1/task/' + taskId, {
        ...updatedTask,
      }).then(response => {
        console.log('TaskCard updated:', response.data)
        return true
      }).catch(error => {
        console.error('Update TaskCard error:', error)
        throw new Error('update TaskCard')
      })
      return true
    }
    return false
  }
  const getTasks = async (categoryID: number): Promise<Task[]> => {
    console.log('Get tasks - function')
    const token = localStorage.getItem('etm-token')
    if (!token) {
      throw new Error('No token')
    }
    const request = axios.create({
      baseURL: import.meta.env.VITE_BACKEND_URL,
      timeout: 1000,
      headers: { Authorization: `Bearer ${token}` },
    })
    try {
      const response = await request.get('/api/v1/tasks/' + categoryID)
      tasks.value = response.data
      return response.data
    } catch (error) {
      console.error('Get tasks error:', error)
      throw new Error('get tasks failed')
    }
  }
  return { tasks, addTask, removeTask, updateTask, getTasks }
})
