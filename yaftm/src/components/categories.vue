<script setup lang="ts">
  import CategoryComponent from '@/components/categorycomponent.vue'
  import { useCategoryStore } from '@/stores/category.ts'
  import type { Category } from '@/types/category.ts'
  const categoryStore = useCategoryStore()

  let categories:Category[] = []
  const categoriesDisplay = ref<Category[]>([])
  const dialog = ref(false)

  onMounted(async () => {
    try {
      categories = await categoryStore.getCategories()
      categoriesDisplay.value = categories
    } catch (error) {
      console.log('Error fetching categories')
    }
  })
  const addCategory = () => {
    dialog.value = true
  }

</script>

<template>
  <v-tabs>
    <v-tab v-for="category in categories" :key="category.id">
      {{ category.name }}
    </v-tab>
    <v-tab @click="addCategory">
      Add
    </v-tab>
    <v-tab-item v-for="category in categories" :key="category.id"><
      CategoryComponent :category="category"/>
    </v-tab-item>
  </v-tabs>
  <v-dialog v-model="dialog" persistent max-width="600px">
    <template v-slot:activator="{ on, attrs }">
      <v-btn color="primary" dark v-bind="attrs" v-on="on">
       Add a category
      </v-btn>
    </template>
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
</template>

<style scoped lang="sass">

</style>
