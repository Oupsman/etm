import { ref } from 'vue'
import type { Ref } from 'vue'
import { defineStore } from 'pinia'
import axios from 'axios'

import type { Category, NewCategory } from '@/types/category'


export const useCategoryStore = defineStore('category', () => {
  const categories = ref([] as Category[])
  let activeCategoryID = localStorage.getItem('etm-active-category-id') || '-1'
  const addCategory = (category: NewCategory): void => {
    const token = localStorage.getItem('etm-token')
    if (!token) {
      throw new Error('No token')
    }
    const request = axios.create({
      baseURL: import.meta.env.VITE_BACKEND_URL,
      timeout: 1000,
      headers: { Authorization: `Bearer ${token}` },
    })
    request.post(import.meta.env.VITE_BACKEND_URL + '/api/v1/category', {
      ...category,
      token,
    }).then(response => {
      console.log('categoryCard created:', response.data)

    }).catch(error => {
      console.error('Create categoryCard error:', error)
      throw new Error('Create category failed')
    })
  }

  const removeCategory = (categoryToDelete: Category): boolean => {
    categories.value = categories.value.filter(category => category.ID !== categoryToDelete.ID)

    console.log('Delete category - function')
    confirm('Are you sure you want to delete this category?')
    const token = localStorage.getItem('etm-token')
    if (!token) {
      throw new Error('No token')
    }
    const request = axios.create({
      baseURL: import.meta.env.VITE_BACKEND_URL,
      timeout: 1000,
      headers: { Authorization: `Bearer ${token}` },
    })
    request.delete(import.meta.env.VITE_BACKEND_URL + '/api/v1/category/' + categoryToDelete.ID).then(response => {
      console.log('category deleted:', response.data)
      return true
    }).catch(error => {
      console.error('Delete category error:', error)
      throw new Error('delete category')
    })
    return true
  }

  const updateCategory = (categoryId: number, updatedCategory: Category): boolean => {
    const index = categories.value.findIndex(category => category.ID === categoryId)
    if (index !== -1) {
      categories.value[index] = { ...categories.value[index], ...updatedCategory }
      console.log('categoryCard to save', updatedCategory)
      const token = localStorage.getItem('etm-token')
      if (!token) {
        throw new Error('No token')
      }
      const request = axios.create({
        baseURL: import.meta.env.VITE_BACKEND_URL,
        timeout: 1000,
        headers: { Authorization: `Bearer ${token}` },
      })
      request.post(import.meta.env.VITE_BACKEND_URL + '/api/v1/category/' + categoryId, {
        ...updatedCategory,
      }).then(response => {
        console.log('categoryCard updated:', response.data)
        return true
      }).catch(error => {
        console.error('Update categoryCard error:', error)
        throw new Error('update categoryCard')
      })
      return true
    }
    return false
  }

  const getCategories = async (): Promise<Category[]> => {
    console.log('Get categories - function')
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
      const response = await request.get('/api/v1/categories')
      categories.value = response.data.categories
      if (activeCategoryID === '-1') {
        localStorage.setItem('etm-active-category-id', categories.value[0].ID.toString())
        activeCategoryID = categories.value[0].ID.toString()
      }
      return response.data.categories
    } catch (error) {
      console.error('Get categories error:', error)
      throw new Error('get categories failed')
    }
  }
  return { categories, addCategory, removeCategory, updateCategory, getCategories }
})
