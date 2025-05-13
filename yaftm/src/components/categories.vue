<script setup lang="ts">
  import CategoryComponent from '@/components/categorycomponent.vue'
  import { useCategoryStore } from '@/stores/category.ts'
  import type { Category, NewCategory } from '@/types/category.ts'
  const categoryStore = useCategoryStore()

  let categories:Category[] = []
  const categoriesDisplay = ref<Category[]>([])
  const dialog = ref(false)
  const categoryName = ref('')
  const categoryColor = ref('#000000')
  const activeTab = ref(0)
  const category = ref({} as Category)
  onMounted(async () => {
    try {
      categories = await categoryStore.getCategories()
      categoriesDisplay.value = categories
      console.log("Categories: ", categories)
      setActiveTab(categories[0].ID)
    } catch (error) {
      console.log('Error fetching categories')
    }

  })
  const triggerDialogCategory = () => {
    dialog.value = true
  }

  const addCategory = () => {
    dialog.value = false
    if (categoryName.value && categoryColor.value) {
      const newCategory:NewCategory = {
        name: categoryName.value,
        color: categoryColor.value,
      }
      categoryStore.addCategory(newCategory)
    }
  }

  const setActiveTab = (categoryId: number) => {
    console.log("Switching active tab to category ID: " + categoryId)
    activeTab.value = categoryId
  }


</script>

<template>
  <v-container style="position: relative; top: 0; left: 0; height: 100vh;">
    <v-row style="position: relative; top: 0; left: 0; height: 5vh;">
      <v-col>
        <v-tabs v-model="activeTab">
          <v-tab
            v-for="category in categoriesDisplay"
            :key="category.ID"
            :style="{ backgroundColor: category.color }"
            @click="setActiveTab(category.ID)"
          >
            {{ category.name }}
          </v-tab>
          <v-btn @click="triggerDialogCategory">
            Add
          </v-btn>
        </v-tabs>
      </v-col>
    </v-row>
    <v-row style="position: relative; top: 0; left: 0; height: 65vh;">
      <v-col>
        <v-tabs-items v-model="activeTab">
          <v-tab-item
            v-for="category in categoriesDisplay"
            :key="category.ID"
          >
            <CategoryComponent v-if="activeTab === category.ID" :category="category" />
          </v-tab-item>
        </v-tabs-items>
      </v-col>
    </v-row>

    <v-dialog v-model="dialog" persistent max-width="600px">
      <v-card>
        <v-card-title>
          <span class="headline">Add a new category</span>
        </v-card-title>
        <v-card-text>
          <v-container>
            <v-row>
              <v-col cols="12">
                <v-text-field
                  label="Name"
                  v-model="categoryName"
                  required
                ></v-text-field>
              </v-col>
              <v-col cols="12">
                <v-text-field
                  label="Color"
                  v-model="categoryColor"
                  type="color"
                  required
                ></v-text-field>
              </v-col>
            </v-row>
          </v-container>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="blue darken-1" text @click="dialog = false">
            Cancel
          </v-btn>
          <v-btn color="blue darken-1" text @click="addCategory">
            Add
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>

<style scoped lang="sass">
  fill-height
    height: 100%
</style>
