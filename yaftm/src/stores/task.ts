import { ref } from 'vue'
import { defineStore } from 'pinia'
import axios from 'axios'
import { useSnackbarStore } from '@/stores/snackbar';

import type { Task } from '@/types/task'


export const useTaskStore = defineStore('task', () => {
  const snackbar = useSnackbarStore();
  const tasks = ref([] as Task[])

  const addTask = (task: Task): Promise<Task> => {
    return new Promise((resolve, reject) => {
      const token = localStorage.getItem('etm-token');
      if (!token) {
        reject(new Error('No token'));
        return;
      }
      const request = axios.create({
        baseURL: import.meta.env.VITE_BACKEND_URL,
        timeout: 1000,
        headers: { Authorization: `Bearer ${token}` },
      });
      request.post(import.meta.env.VITE_BACKEND_URL + '/api/v1/task', {
        ...task,
        token,
      }).then(response => {
        snackbar.showSnackbar({
          message: 'Task added successfully.',
          color: 'success',
        })
        tasks.value.push(task)
        resolve(response.data);
      }).catch(error => {
        snackbar.showSnackbar({
          message: 'Unable to add a task ' + error.message,
          color: 'error',
        });
        reject(new Error('Create task failed'));
      });
    });
  };

  const deleteTask = (taskToDelete: Task): boolean => {
    tasks.value = tasks.value.filter(task => task.ID !== taskToDelete.ID)
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
      snackbar.showSnackbar({
        message: 'Task deleted successfully ',
        color: 'success',
      });
      return true
    }).catch(error => {
      snackbar.showSnackbar({
        message: 'Unable to delete a task ' + error.message,
        color: 'error',
      });
      throw new Error('delete task')
    })
    return true
  }

  const updateTask = (taskId: number, updatedTask: Task): boolean => {
    console.log('Store UpdateTask')
    const index = tasks.value.findIndex(task => task.ID === taskId)
    if (index !== -1) {
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
        tasks.value[index] = { ...tasks.value[index], ...updatedTask }
        snackbar.showSnackbar({
          message: 'Task updated successfully ',
          color: 'success',
        });
        return true
      }).catch(error => {
        snackbar.showSnackbar({
          message: 'Unable to update a task ' + error.message,
          color: 'error',
        });
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
      throw new Error('get tasks failed')
    }
  }
  return { tasks, addTask, deleteTask, updateTask, getTasks }
})
