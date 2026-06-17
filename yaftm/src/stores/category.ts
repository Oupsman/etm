import { ref } from 'vue'
import { defineStore } from 'pinia'

import { axiosInstance } from '@/plugins/axios'
import type { Category } from '@/types/category'

export const useCategoryStore = defineStore('category', () => {
  const categories = ref([] as Category[])

  const addCategory = async (category: Category): Promise<void> => {
    await axiosInstance.post('/api/v1/category', category)
  }

  const removeCategory = async (categoryToDelete: Category): Promise<boolean> => {
    categories.value = categories.value.filter(category => category.ID !== categoryToDelete.ID)
    await axiosInstance.delete('/api/v1/category/' + categoryToDelete.ID)
    return true
  }

  const updateCategory = async (categoryId: number, updatedCategory: Category): Promise<boolean> => {
    const index = categories.value.findIndex(category => category.ID === categoryId)
    if (index === -1) return false
    categories.value[index] = { ...categories.value[index], ...updatedCategory }
    await axiosInstance.post('/api/v1/category/' + categoryId, updatedCategory)
    return true
  }

  const getCategories = async (): Promise<Category[]> => {
    const response = await axiosInstance.get('/api/v1/categories')
    categories.value = response.data.categories ?? []
    return categories.value
  }

  const getCategoryShares = async (categoryId: number) => {
    const { data } = await axiosInstance.get(`/api/v1/category/${categoryId}/shares`)
    return (data.groups ?? []) as { ID: number; name: string }[]
  }

  const shareCategory = async (categoryId: number, groupId: number): Promise<void> => {
    await axiosInstance.post(`/api/v1/category/${categoryId}/share`, { group_id: groupId })
  }

  const unshareCategory = async (categoryId: number, groupId: number): Promise<void> => {
    await axiosInstance.delete(`/api/v1/category/${categoryId}/share/${groupId}`)
  }

  return { categories, addCategory, removeCategory, updateCategory, getCategories, getCategoryShares, shareCategory, unshareCategory }
})
