<script setup lang="ts">
  import { ref, watch } from 'vue'
  import CategoryComponent from '@/components/categorycomponent.vue'
  import { useCategoryStore } from '@/stores/category.ts'
  import { useAuthStore } from '@/stores/auth'
  import type { Category } from '@/types/category.ts'

  const categoryStore = useCategoryStore()
  const authStore = useAuthStore()

  let categories:Category[] = []
  const categoriesDisplay = ref<Category[]>([])
  const dialog = ref(false)
  const categoryName = ref('')
  const categoryColor = ref('#EE2222')
  const activeTab = ref(0)

  const loadCategories = async () => {
    try {
      categories = await categoryStore.getCategories()
      categoriesDisplay.value = categories
      if (categories.length > 0) setActiveTab(categories[0].ID)
    } catch (error) {
      console.log('Error fetching categories: ', error)
    }
  }

  watch(() => authStore.isAuthenticated, (isAuth) => {
    if (isAuth) loadCategories()
  }, { immediate: true })
  const triggerDialogCategory = () => {
    dialog.value = true
  }

  const addCategory = async () => {
    dialog.value = false
    if (categoryName.value && categoryColor.value) {
      const newCategory:Category = {
        ID:0,
        name: categoryName.value,
        color: categoryColor.value,
      }
      try {
        await categoryStore.addCategory(newCategory)
        await loadCategories()
      } catch(error) {
        console.log('Error adding category: ', error)
      }

    }
  }

  const setActiveTab = (categoryId: number) => {
    console.log('Switching active tab to category ID: ' + categoryId)
    activeTab.value = categoryId
  }


</script>

<template>
  <div style="display: flex; flex-direction: column; height: 100%; overflow: hidden;">
    <div style="flex: 0 0 auto;">
      <v-tabs v-model="activeTab">
          <v-tab
            v-for="category in categoriesDisplay"
            :key="category.ID"
            :value="category.ID"
            :style="{
              backgroundColor: category.color,
              opacity: activeTab === category.ID ? 1 : 0.55,
              fontWeight: activeTab === category.ID ? '700' : '400',
              borderBottom: activeTab === category.ID ? '3px solid rgba(0,0,0,0.4)' : '3px solid transparent',
            }"
            @click="setActiveTab(category.ID)"
          >
            {{ category.name }}
          </v-tab>
          <v-btn @click="triggerDialogCategory">
            Add
          </v-btn>
        </v-tabs>
    </div>
    <div style="flex: 1; min-height: 0; overflow: hidden; position: relative;">
      <template v-for="category in categoriesDisplay" :key="category.ID">
        <CategoryComponent
          v-if="activeTab === category.ID"
          :category-i-d="category.ID"
          style="height: 100%;"
        />
      </template>
    </div>

    <v-dialog v-model="dialog" max-width="600px" persistent>
      <v-card>
        <v-card-title>
          <span class="headline">Add a new category</span>
        </v-card-title>
        <v-card-text>
          <v-container>
            <v-row>
              <v-col cols="12">
                <v-text-field
                  v-model="categoryName"
                  label="Name"
                  required
                />
              </v-col>
              <v-col cols="12">
                <v-text-field
                  v-model="categoryColor"
                  label="Color"
                  required
                  type="color"
                />
              </v-col>
            </v-row>
          </v-container>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn color="blue darken-1" text @click="dialog = false">
            Cancel
          </v-btn>
          <v-btn color="blue darken-1" text @click="addCategory">
            Add
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<style scoped lang="sass">
  fill-height
    height: 100%
</style>
