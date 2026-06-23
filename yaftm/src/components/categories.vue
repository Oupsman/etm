<script setup lang="ts">
  import { ref, watch, computed } from 'vue'
  import { useI18n } from 'vue-i18n'
  import CategoryComponent from '@/components/categorycomponent.vue'
  import { useCategoryStore } from '@/stores/category.ts'
  import { useGroupStore } from '@/stores/group'
  import { useAuthStore } from '@/stores/auth'
  import { useUserStore } from '@/stores/user'
  import { useDragStore } from '@/stores/drag'
  import { useTaskStore } from '@/stores/task'
  import type { Category } from '@/types/category.ts'
  import type { Task } from '@/types/task'

  const { t } = useI18n()
  const categoryStore = useCategoryStore()
  const groupStore = useGroupStore()
  const authStore = useAuthStore()
  const userStore = useUserStore()
  const dragStore = useDragStore()
  const taskStore = useTaskStore()

  const isDragging = computed(() => dragStore.draggingTask !== null)

  const onTaskDropToCategory = async (targetCategory: Category) => {
    const task = dragStore.draggingTask
    dragStore.endDrag()
    if (!task) return
    if (task.categoryid === targetCategory.ID) return
    const updated: Task = { ...task, categoryid: targetCategory.ID }
    await taskStore.updateTask(task.ID, updated)
    dragStore.notifyMoved()
  }

  let categories: Category[] = []
  const categoriesDisplay = ref<Category[]>([])
  const dialog = ref(false)
  const categoryName = ref('')
  const categoryColor = ref('#EE2222')
  const categoryGroupID = ref<number | null>(null)
  const activeTab = ref(0)

  const currentUserID = ref(0)

  const shareDialog = ref(false)
  const sharingCategory = ref<Category | null>(null)
  const currentShares = ref<{ ID: number; name: string }[]>([])
  const sharingGroupID = ref<number | null>(null)

  const deleteDialog = ref(false)
  const categoryToDelete = ref<Category | null>(null)

  const userGroups = computed(() =>
    groupStore.groups.filter(g => !currentShares.value.some(s => s.ID === g.ID))
  )

  const loadCategories = async () => {
    try {
      categories = await categoryStore.getCategories()
      categoriesDisplay.value = categories
      if (categories.length === 0) return

      const userData = await userStore.getUser()
      currentUserID.value = userData.ID
      const savedId = userData.active_category_id
      const exists = savedId > 0 && categories.some(c => c.ID === savedId)
      setActiveTab(exists ? savedId : categories[0].ID, false)
    } catch (error) {
      console.log('Error fetching categories: ', error)
    }
  }

  watch(() => authStore.isAuthenticated, (isAuth) => {
    if (isAuth) loadCategories()
  }, { immediate: true })

  const triggerDialogCategory = async () => {
    categoryName.value = ''
    categoryColor.value = '#EE2222'
    categoryGroupID.value = null
    await groupStore.fetchGroups(false)
    dialog.value = true
  }

  const addCategory = async () => {
    dialog.value = false
    if (categoryName.value && categoryColor.value) {
      const newCategory: Category = {
        ID: 0,
        name: categoryName.value,
        color: categoryColor.value,
        group_id: categoryGroupID.value ?? undefined,
      }
      try {
        await categoryStore.addCategory(newCategory)
        await loadCategories()
      } catch (error) {
        console.log('Error adding category: ', error)
      }
    }
  }

  const setActiveTab = (categoryId: number, persist = true) => {
    activeTab.value = categoryId
    if (persist) userStore.saveActiveCategory(categoryId)
  }

  const isReadOnly = (cat: Category): boolean => {
    if (!cat.shared) return false
    return cat.user_role === 'reader'
  }

  const canShare = (cat: Category): boolean => {
    return !cat.shared && cat.userid === currentUserID.value
  }

  const openShareDialog = async (cat: Category) => {
    sharingCategory.value = cat
    sharingGroupID.value = null
    currentShares.value = await categoryStore.getCategoryShares(cat.ID)
    await groupStore.fetchGroups(false)
    shareDialog.value = true
  }

  const addShare = async () => {
    if (!sharingCategory.value || !sharingGroupID.value) return
    await categoryStore.shareCategory(sharingCategory.value.ID, sharingGroupID.value)
    currentShares.value = await categoryStore.getCategoryShares(sharingCategory.value.ID)
    sharingGroupID.value = null
    await groupStore.fetchGroups(false)
  }

  const removeShare = async (groupId: number) => {
    if (!sharingCategory.value) return
    await categoryStore.unshareCategory(sharingCategory.value.ID, groupId)
    currentShares.value = currentShares.value.filter(s => s.ID !== groupId)
    await groupStore.fetchGroups(false)
  }

  const confirmDeleteCategory = (cat: Category) => {
    categoryToDelete.value = cat
    deleteDialog.value = true
  }

  const deleteCategory = async () => {
    if (!categoryToDelete.value) return
    deleteDialog.value = false
    try {
      await categoryStore.removeCategory(categoryToDelete.value)
      await loadCategories()
    } catch (error) {
      console.log('Error deleting category: ', error)
    }
    categoryToDelete.value = null
  }
</script>

<template>
  <div style="display: flex; flex-direction: column; height: 100%; overflow: hidden;">
    <div style="flex: 0 0 auto;">
      <v-tabs v-model="activeTab">
        <v-tab
          v-for="category in categoriesDisplay"
          :key="category.ID"
          :class="['category-tab', { 'drop-zone': isDragging && category.ID !== dragStore.draggingFromCategoryId }]"
          :style="{
            '--cat-bg': category.color,
            opacity: activeTab === category.ID ? 1 : 0.55,
            fontWeight: activeTab === category.ID ? '700' : '400',
            borderBottom: activeTab === category.ID ? '3px solid rgba(0,0,0,0.4)' : '3px solid transparent',
          }"
          :value="category.ID"
          @click="setActiveTab(category.ID)"
          @dragover.prevent
          @drop.prevent="onTaskDropToCategory(category)"
        >
          <span class="mr-1">{{ category.name }}</span>
          <v-icon v-if="category.shared" class="mr-1" size="x-small" :title="t('category.sharedCategory')">mdi-account-group</v-icon>
          <v-icon v-if="isReadOnly(category)" size="x-small" :title="t('category.readOnly')">mdi-eye</v-icon>
          <v-btn
            v-if="canShare(category)"
            class="ml-1"
            icon
            size="x-small"
            :title="t('category.manageSharing')"
            variant="text"
            @click.stop="openShareDialog(category)"
          >
            <v-icon size="x-small">mdi-share-variant</v-icon>
          </v-btn>
          <v-btn
            v-if="canShare(category)"
            icon
            size="x-small"
            variant="text"
            class="ml-1"
            :title="t('category.deleteTooltip')"
            @click.stop="confirmDeleteCategory(category)"
          >
            <v-icon size="x-small">mdi-delete</v-icon>
          </v-btn>
        </v-tab>
        <div class="add-tab" :title="t('category.addTooltip')" @click="triggerDialogCategory">
          <v-icon size="small">mdi-plus</v-icon>
        </div>
      </v-tabs>
    </div>
    <div style="flex: 1; min-height: 0; overflow: hidden; position: relative;">
      <template v-for="category in categoriesDisplay" :key="category.ID">
        <CategoryComponent
          v-if="activeTab === category.ID"
          :category-i-d="category.ID"
          :readonly="isReadOnly(category)"
          style="height: 100%;"
        />
      </template>
    </div>

    <!-- Add category dialog -->
    <v-dialog v-model="dialog" max-width="600px" persistent>
      <v-card>
        <v-card-title>
          <span class="headline">{{ t('category.addTitle') }}</span>
        </v-card-title>
        <v-card-text>
          <v-container>
            <v-row>
              <v-col cols="12">
                <v-text-field v-model="categoryName" :label="t('category.name')" required />
              </v-col>
              <v-col cols="12">
                <v-text-field v-model="categoryColor" :label="t('category.color')" required type="color" />
              </v-col>
              <v-col cols="12">
                <v-select
                  v-model="categoryGroupID"
                  :items="groupStore.groups"
                  item-title="name"
                  item-value="ID"
                  :label="t('category.shareWithGroup')"
                  clearable
                  :no-data-text="groupStore.groups.length === 0 ? t('category.noGroups') : t('category.noGroupsData')"
                />
              </v-col>
            </v-row>
          </v-container>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn color="blue darken-1" text @click="dialog = false">{{ t('category.cancel') }}</v-btn>
          <v-btn color="blue darken-1" text @click="addCategory">{{ t('category.add') }}</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Delete category confirmation dialog -->
    <v-dialog v-model="deleteDialog" max-width="400">
      <v-card>
        <v-card-title>{{ t('category.deleteTitle', { name: categoryToDelete?.name }) }}</v-card-title>
        <v-card-text>{{ t('category.deleteConfirm') }}</v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn @click="deleteDialog = false">{{ t('category.cancel') }}</v-btn>
          <v-btn color="error" @click="deleteCategory">{{ t('category.delete') }}</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Share management dialog -->
    <v-dialog v-model="shareDialog" max-width="500">
      <v-card>
        <v-card-title>{{ t('category.shareTitle', { name: sharingCategory?.name }) }}</v-card-title>
        <v-card-text>
          <div v-if="currentShares.length > 0">
            <div class="text-subtitle-2 mb-2">{{ t('category.sharedWith') }}</div>
            <v-chip
              v-for="share in currentShares"
              :key="share.ID"
              class="mr-2 mb-2"
              closable
              @click:close="removeShare(share.ID)"
            >
              <v-icon start>mdi-account-group</v-icon>
              {{ share.name }}
            </v-chip>
            <v-divider class="my-3" />
          </div>
          <div class="text-subtitle-2 mb-2">{{ t('category.addGroup') }}</div>
          <v-select
            v-model="sharingGroupID"
            :items="userGroups"
            item-title="name"
            item-value="ID"
            :label="t('category.selectGroup')"
            :no-data-text="groupStore.groups.length === 0 ? t('category.noGroups') : t('category.allGroupsAdded')"
            clearable
          />
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn @click="shareDialog = false">{{ t('category.close') }}</v-btn>
          <v-btn color="primary" :disabled="!sharingGroupID" @click="addShare">{{ t('category.share') }}</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<style scoped lang="sass">
  fill-height
    height: 100%

  .drop-zone
    outline: 2px dashed rgba(0, 0, 0, 0.5)
    outline-offset: -3px
    opacity: 0.85 !important

  .add-tab
    display: flex
    align-items: center
    justify-content: center
    align-self: stretch
    padding: 0 14px
    background-color: rgba(0, 0, 0, 0.08)
    border-bottom: 3px solid transparent
    cursor: pointer
    transition: background-color 0.2s ease

    &:hover
      background-color: rgba(0, 0, 0, 0.15)
</style>
