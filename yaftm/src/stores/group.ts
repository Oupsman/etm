import { defineStore } from 'pinia'
import { ref } from 'vue'
import { axiosInstance } from '@/plugins/axios'
import type { Group, GroupDetail, UserSearchResult } from '@/types/group'

export const useGroupStore = defineStore('group', () => {
  const groups = ref<Group[]>([])

  const fetchGroups = async (adminView = false) => {
    const url = adminView ? '/api/v1/admin/groups' : '/api/v1/groups'
    const { data } = await axiosInstance.get(url)
    groups.value = data.groups
  }

  const createGroup = async (name: string): Promise<Group> => {
    const { data } = await axiosInstance.post('/api/v1/group', { name })
    return data.group
  }

  const deleteGroup = async (groupId: number) => {
    await axiosInstance.delete(`/api/v1/group/${groupId}`)
  }

  const getGroup = async (groupId: number): Promise<GroupDetail> => {
    const { data } = await axiosInstance.get(`/api/v1/group/${groupId}`)
    return { group: data.group, members: data.members }
  }

  const addMember = async (groupId: number, userId: number, role: string) => {
    await axiosInstance.post(`/api/v1/group/${groupId}/member`, { user_id: userId, role })
  }

  const updateMember = async (groupId: number, userId: number, role: string) => {
    await axiosInstance.put(`/api/v1/group/${groupId}/member/${userId}`, { user_id: userId, role })
  }

  const removeMember = async (groupId: number, userId: number) => {
    await axiosInstance.delete(`/api/v1/group/${groupId}/member/${userId}`)
  }

  const searchUsers = async (query: string): Promise<UserSearchResult[]> => {
    const { data } = await axiosInstance.get('/api/v1/users/search', { params: { q: query } })
    return data.users
  }

  return {
    groups,
    fetchGroups,
    createGroup,
    deleteGroup,
    getGroup,
    addMember,
    updateMember,
    removeMember,
    searchUsers,
  }
})
