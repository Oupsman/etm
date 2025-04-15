import { ref } from 'vue'
import type { Ref } from 'vue'
import { defineStore } from 'pinia'
import axios from 'axios'

import type { Task } from '@/types/task'


export const useTaskStore = defineStore('task', () => {
  const tasks = ref([] as Task[])

  const addTask = (task: Task): void => {
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
      console.log('Task created:', response.data)
    }).catch(error => {
      console.error('Create Task error:', error)
      throw new Error('Create task failed')
    })
  }

  const removeTask = (taskId: number): boolean => {
    tasks.value = tasks.value.filter((task) => task.id !== taskId)
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
    request.delete(import.meta.env.VITE_BACKEND_URL + '/api/v1/task/' + taskId).then(response => {
      console.log('task deleted:', response.data)
      return true
    }).catch(error => {
      console.error('Delete task error:', error)
      throw new Error('delete task')
    })
    return true
  }

  const updateTask = (taskId: number, updatedTask: Task): boolean => {
    const index = tasks.value.findIndex((task) => task.id === taskId)
    if (index !== -1) {
      tasks.value[index] = { ...tasks.value[index], ...updatedTask }
      console.log('Task to save', updatedTask)
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
        console.log('Task updated:', response.data)
        return true
      }).catch(error => {
        console.error('Update Task error:', error)
        throw new Error('update Task')
      })
      return true
    }
    return false
  }
  const getTasks = async (): Promise<Task[]> => {
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
      const response = await request.get('/api/v1/tasks')
      tasks.value = response.data
      return response.data
    } catch (error) {
      console.error('Get tasks error:', error)
      throw new Error('get tasks failed')
    }
  }
  return { tasks, addTask, removeTask, updateTask, getTasks }
})
