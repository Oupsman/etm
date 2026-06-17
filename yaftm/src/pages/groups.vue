<script setup lang="ts">
  import { ref, computed, onMounted } from 'vue'
  import { useGroupStore } from '@/stores/group'
  import { useUserStore } from '@/stores/user'
  import type { Group, GroupDetail, GroupMember, UserSearchResult } from '@/types/group'
  import type { User } from '@/types/user'

  const groupStore = useGroupStore()
  const userStore = useUserStore()

  const currentUser = ref<User | null>(null)
  const isAdmin = computed(() => currentUser.value?.isadmin === 'true')

  const selectedGroup = ref<GroupDetail | null>(null)
  const loading = ref(false)
  const groupLoading = ref(false)

  const createDialog = ref(false)
  const newGroupName = ref('')

  const addMemberDialog = ref(false)
  const searchQuery = ref('')
  const searchResults = ref<UserSearchResult[]>([])
  const searching = ref(false)
  const selectedRole = ref('reader')
  const selectedUser = ref<UserSearchResult | null>(null)

  const confirmDeleteDialog = ref(false)
  const groupToDelete = ref<Group | null>(null)

  const roles = ['reader', 'writer', 'owner']

  onMounted(async () => {
    currentUser.value = await userStore.getUser()
    await loadGroups()
  })

  const loadGroups = async () => {
    loading.value = true
    try {
      await groupStore.fetchGroups(isAdmin.value)
    } finally {
      loading.value = false
    }
  }

  const selectGroup = async (group: Group) => {
    groupLoading.value = true
    try {
      selectedGroup.value = await groupStore.getGroup(group.ID)
    } finally {
      groupLoading.value = false
    }
  }

  const isGroupOwner = computed(() => {
    if (!selectedGroup.value || !currentUser.value) return false
    return selectedGroup.value.group.owner_id === currentUser.value.ID || isAdmin.value
  })

  const createGroup = async () => {
    if (!newGroupName.value.trim()) return
    await groupStore.createGroup(newGroupName.value.trim())
    newGroupName.value = ''
    createDialog.value = false
    await loadGroups()
  }

  const confirmDelete = (group: Group) => {
    groupToDelete.value = group
    confirmDeleteDialog.value = true
  }

  const deleteGroup = async () => {
    if (!groupToDelete.value) return
    await groupStore.deleteGroup(groupToDelete.value.ID)
    if (selectedGroup.value?.group.ID === groupToDelete.value.ID) {
      selectedGroup.value = null
    }
    groupToDelete.value = null
    confirmDeleteDialog.value = false
    await loadGroups()
  }

  const searchUsers = async () => {
    if (searchQuery.value.length < 2) return
    searching.value = true
    try {
      searchResults.value = await groupStore.searchUsers(searchQuery.value)
    } finally {
      searching.value = false
    }
  }

  const addMember = async () => {
    if (!selectedUser.value || !selectedGroup.value) return
    await groupStore.addMember(selectedGroup.value.group.ID, selectedUser.value.ID, selectedRole.value)
    selectedGroup.value = await groupStore.getGroup(selectedGroup.value.group.ID)
    addMemberDialog.value = false
    searchQuery.value = ''
    searchResults.value = []
    selectedUser.value = null
    selectedRole.value = 'reader'
  }

  const updateRole = async (member: GroupMember, newRole: string) => {
    if (!selectedGroup.value) return
    await groupStore.updateMember(selectedGroup.value.group.ID, member.user_id, newRole)
    member.role = newRole
  }

  const removeMember = async (member: GroupMember) => {
    if (!selectedGroup.value) return
    await groupStore.removeMember(selectedGroup.value.group.ID, member.user_id)
    selectedGroup.value = await groupStore.getGroup(selectedGroup.value.group.ID)
  }
</script>

<template>
  <v-container fluid class="pa-4" style="height: 100%; overflow: auto;">
    <v-row>
      <!-- Group list -->
      <v-col cols="12" md="4">
        <v-card>
          <v-card-title class="d-flex align-center">
            Groups
            <v-spacer />
            <v-btn icon size="small" @click="createDialog = true">
              <v-icon>mdi-plus</v-icon>
            </v-btn>
          </v-card-title>
          <v-divider />
          <v-progress-linear v-if="loading" indeterminate />
          <v-list v-else>
            <v-list-item
              v-for="group in groupStore.groups"
              :key="group.ID"
              :active="selectedGroup?.group.ID === group.ID"
              active-color="primary"
              @click="selectGroup(group)"
            >
              <v-list-item-title>{{ group.name }}</v-list-item-title>
              <template #append>
                <v-btn
                  v-if="group.owner_id === currentUser?.ID || isAdmin"
                  icon
                  size="x-small"
                  variant="text"
                  color="error"
                  @click.stop="confirmDelete(group)"
                >
                  <v-icon>mdi-delete</v-icon>
                </v-btn>
              </template>
            </v-list-item>
            <v-list-item v-if="groupStore.groups.length === 0">
              <v-list-item-title class="text-disabled">No groups yet</v-list-item-title>
            </v-list-item>
          </v-list>
        </v-card>
      </v-col>

      <!-- Group detail -->
      <v-col cols="12" md="8">
        <v-card v-if="selectedGroup" min-height="300">
          <v-card-title class="d-flex align-center">
            {{ selectedGroup.group.name }}
            <v-chip class="ml-2" size="small" color="primary">
              {{ selectedGroup.members.length }} member{{ selectedGroup.members.length !== 1 ? 's' : '' }}
            </v-chip>
            <v-spacer />
            <v-btn
              v-if="isGroupOwner"
              variant="outlined"
              size="small"
              prepend-icon="mdi-account-plus"
              @click="addMemberDialog = true"
            >
              Add member
            </v-btn>
          </v-card-title>
          <v-divider />
          <v-progress-linear v-if="groupLoading" indeterminate />
          <v-list v-else>
            <v-list-item v-for="member in selectedGroup.members" :key="member.user_id">
              <template #prepend>
                <v-icon>mdi-account</v-icon>
              </template>
              <v-list-item-title>{{ member.username }}</v-list-item-title>
              <v-list-item-subtitle>{{ member.email }}</v-list-item-subtitle>
              <template #append>
                <div class="d-flex align-center gap-2">
                  <v-select
                    v-if="isGroupOwner && member.role !== 'owner'"
                    :model-value="member.role"
                    :items="roles"
                    density="compact"
                    hide-details
                    style="min-width: 110px;"
                    @update:model-value="updateRole(member, $event)"
                  />
                  <v-chip v-else :color="member.role === 'owner' ? 'warning' : 'default'" size="small">
                    {{ member.role }}
                  </v-chip>
                  <v-btn
                    v-if="isGroupOwner && member.role !== 'owner'"
                    icon
                    size="x-small"
                    variant="text"
                    color="error"
                    @click="removeMember(member)"
                  >
                    <v-icon>mdi-account-remove</v-icon>
                  </v-btn>
                </div>
              </template>
            </v-list-item>
          </v-list>
        </v-card>
        <v-card v-else min-height="300" class="d-flex align-center justify-center">
          <v-card-text class="text-center text-disabled">
            <v-icon size="64" class="mb-2">mdi-account-group</v-icon>
            <div>Select a group to see its members</div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <!-- Create group dialog -->
    <v-dialog v-model="createDialog" max-width="400">
      <v-card>
        <v-card-title>New group</v-card-title>
        <v-card-text>
          <v-text-field
            v-model="newGroupName"
            label="Group name"
            autofocus
            @keyup.enter="createGroup"
          />
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn @click="createDialog = false">Cancel</v-btn>
          <v-btn color="primary" :disabled="!newGroupName.trim()" @click="createGroup">Create</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Add member dialog -->
    <v-dialog v-model="addMemberDialog" max-width="500">
      <v-card>
        <v-card-title>Add member to {{ selectedGroup?.group.name }}</v-card-title>
        <v-card-text>
          <v-text-field
            v-model="searchQuery"
            label="Search by username or email"
            prepend-inner-icon="mdi-magnify"
            :loading="searching"
            @input="searchUsers"
          />
          <v-list v-if="searchResults.length > 0" density="compact" class="border rounded mb-3">
            <v-list-item
              v-for="u in searchResults"
              :key="u.ID"
              :active="selectedUser?.ID === u.ID"
              active-color="primary"
              @click="selectedUser = u"
            >
              <v-list-item-title>{{ u.username }}</v-list-item-title>
              <v-list-item-subtitle>{{ u.email }}</v-list-item-subtitle>
            </v-list-item>
          </v-list>
          <v-chip v-if="selectedUser" color="primary" class="mb-3" closable @click:close="selectedUser = null">
            {{ selectedUser.username }}
          </v-chip>
          <v-select
            v-model="selectedRole"
            :items="roles"
            label="Role"
          />
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn @click="addMemberDialog = false; searchQuery = ''; searchResults = []; selectedUser = null">Cancel</v-btn>
          <v-btn color="primary" :disabled="!selectedUser" @click="addMember">Add</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Delete confirmation -->
    <v-dialog v-model="confirmDeleteDialog" max-width="360">
      <v-card>
        <v-card-title>Delete group?</v-card-title>
        <v-card-text>
          Delete <strong>{{ groupToDelete?.name }}</strong>? This cannot be undone.
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn @click="confirmDeleteDialog = false">Cancel</v-btn>
          <v-btn color="error" @click="deleteGroup">Delete</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>
