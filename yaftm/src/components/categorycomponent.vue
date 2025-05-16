<script setup lang="ts">
  import { defineProps, ref, onMounted } from 'vue'
  import { useCategoryStore } from '@/stores/category.ts'
  import { useTaskStore } from '@/stores/task.ts'
  import { useAppStore } from '@/stores/app.ts'
  import { useUserStore } from '@/stores/user.ts'
  import type { Task, NewTask } from '@/types/task.ts'
  import { VueDraggableNext as draggable} from 'vue-draggable-next'

  const props = defineProps({
      categoryID: {
        type: Number,
        required: true
      },
    })
      const backlog = ref<Task[]>([])
      const urgentImportant = ref<Task[]>([])
      const nonUrgentImportant = ref<Task[]>([])
      const nonUrgentNonImportant = ref<Task[]>([])
      const urgentNonImportant = ref<Task[]>([])
      const completedTasks = ref<Task[]>([])

      const taskDialog = ref(false)
      const taskName = ref<string>('')
      const taskDescription = ref<string>('')
      const taskDueDate = ref<Date>(new Date())

      const message = ref<string>('')
      const displaySnack = ref(false)
      const formatTask = (task: Task) => {
        return '<div class="task draggable" id="task-' + task.ID + '">' +
          '<span class="ui-icon ui-icon-arrow-4" class="handle"></span> ' +
          '<span class="ui-icon ui-icon-newwin" class="view" title="Name: ' + task.name + ',Comment: ' + task.comment + ',Due Date: ' + task.duedate + '"></span> ' +
          '<span>' + task.name + '</span>' +
          '<button class="taskbutton deletetask"><span class="ui-icon ui-icon-trash"></span></button>' +
          '<button class="taskbutton edittask"><span class="ui-icon ui-icon-pencil"></span></button>' +
          '<div class="modal-task-display" id="details-task-' + task.ID + '"><label>Name: </label><p>' + task.name + '</p> ' +
          '<label>Comment: </label><p>' + task.comment + '</p>' +
          '<label>Due Date: </label><p>' + task.duedate + '</p>' +
          '</div>' +
          '</div>'
      }

      const categoryStore = useCategoryStore()
      const taskStore = useTaskStore()

      const triggerTaskDialog = () => {
        taskDialog.value = true
      }

      const addTask = () => {
        taskDialog.value = false
        console.log(typeof taskDueDate.value)
        if (taskName.value && taskDescription.value && taskDueDate.value) {
          const newTask: NewTask = {
            name: taskName.value,
            comment: taskDescription.value,
            duedate: taskDueDate.value.toISOString(),
            categoryid: props.categoryID,
            isbacklog: true,
          }
          const task: Task = {
            ID: 0,
            iscompleted: false,
            urgency: false,
            priority: false,

            ...newTask,
          }
          backlog.value.push(task)
        }
      }

      const onChange = (evt: any) => {
        console.log("onChange: ", evt)
      }

      const onMove = (evt: any) => {
        const task: Task = evt.draggedContext.element
        // const origin: String = evt.from.attributes.itemkey.nodeValue
        const destination: String = evt.to.attributes.itemkey.nodeValue
        if (destination === "backlog") {
          task.isbacklog = true
          task.iscompleted = false
          task.urgency = false
          task.priority = false
        } else if (destination === "completedTasks") {
          task.isbacklog = false
          task.iscompleted = true
          task.urgency = false
          task.priority = false
        } else if (destination ===  "urgentImportant") {
            task.isbacklog = false
            task.iscompleted = false
            task.urgency = true
            task.priority = true
        } else if (destination === "nonUrgentImportant") {
          task.isbacklog = false
          task.iscompleted = false
          task.urgency = false
          task.priority = true

        } else if (destination === "urgentNonImportant") {
          task.isbacklog = false
          task.iscompleted = false
          task.urgency = true
          task.priority = false
        } else if (destination === "nonUrgentNonImportant") {
          task.isbacklog = false
          task.iscompleted = false
          task.urgency = false
          task.priority = false
        }
         if (taskStore.updateTask(task.ID, task)) {

         } else {
            message.value = "Task update failed"
         }
          displaySnack.value = true
      }

      onMounted(async () => {
        // query tasks from the store
        const tasks = await taskStore.getTasks(props.categoryID)
        // Parse tasks and add them to the respective lists
        tasks.forEach((task: Task) => {
          if (task.isbacklog) {
            backlog.value.push(task)
          } else if (task.iscompleted) {
            completedTasks.value.push(task)
          } else if (task.urgency && task.priority) {
            urgentImportant.value.push(task)
          } else if (!task.urgency && task.priority) {
            nonUrgentImportant.value.push(task)
          } else if (task.urgency && !task.priority) {
            urgentNonImportant.value.push(task)
          } else {
            nonUrgentNonImportant.value.push(task)
          }
        })
      })
</script>

<template>
  <v-container class="fill-height" style="height: 90vh">
    <v-row class="fill-height">
      <!-- Backlog Column -->
      <v-col cols="3" class="d-flex flex-column backlog fill-height" style="position: absolute left: 0 height:80vh">
        <h2>Backlog</h2>
        <v-btn @click="triggerTaskDialog">Add task</v-btn>
        <draggable group="tasks"
                   v-model="backlog"
                   itemKey="backlog"
                   :move="onMove"
                   @change = "onChange">
          <v-card class="mb-2 task" v-for="task in backlog" :key="task.ID">
            <v-icon icon="mdi-checkbox-blank-outline" size="small"> </v-icon> {{ task.name }}
          </v-card>
        </draggable>

      </v-col>

      <!-- Eisenhower Matrix -->
      <v-col cols="6" class="mx-auto fill-height d-flex flex-column " >
        <v-row class="d-flex">
          <v-col class="mx-auto d-flex flex-column UrgentImportant" cols="6" style="position: relative left: 0 height: 40vh">
              <h2>Urgent et Important</h2>
              <draggable group="tasks"
                         v-model="urgentImportant"
                         itemKey="urgentImportant"
                         :move="onMove"
                         @change = "onChange">
                <v-card class="mb-2 task" v-for="task in urgentImportant" :key="task.ID">
                  <v-icon icon="mdi-checkbox-blank-outline" size="small"> </v-icon>  {{ task.name }}
                </v-card>
              </draggable>
          </v-col>
          <v-col cols="6" class="mx-auto d-flex flex-column NotUrgentImportant" style="position: relative left: 0 height: 40vh">
            <h2>Non Urgent et Important</h2>
            <draggable group="tasks"
                       v-model="nonUrgentImportant"
                       itemKey="nonUrgentImportant"
                       :move="onMove"
                       @change = "onChange">

                <v-card class="mb-2 task"  v-for="task in nonUrgentImportant" :key="task.ID">
                  <v-icon icon="mdi-checkbox-blank-outline" size="small"> </v-icon>  {{ task.name }}
                </v-card>
            </draggable>
          </v-col>
          <v-col cols="6" class="mx-auto d-flex flex-column UrgentNotImportant" style="position: relative left: 0 height: 40vh">
              <h2>Urgent et Non Important</h2>
              <draggable group="tasks"
                         v-model="urgentNonImportant"
                         itemKey="urgentNonImportant"
                         :move="onMove"
                         @change = "onChange">
                <v-card class="mb-2 task" v-for="task in urgentNonImportant" :key="task.ID">
                  <v-icon icon="mdi-checkbox-blank-outline" size="small"> </v-icon>  {{ task.name }}
                  </v-card>
              </draggable>
          </v-col>
          <v-col cols="6" class="mx-auto d-flex flex-column NotUrgentNotImportant" style="position: relative left: 0 height: 40vh">
            <h2>Non Urgent et Non Important</h2>
              <draggable group="tasks"
                         v-model="nonUrgentNonImportant"
                         itemKey="nonUrgentNonImportant"
                         :move="onMove"
                         @change = "onChange">

                <v-card class="mb-2 task" v-for="task in nonUrgentNonImportant" :key="task.ID">
                  <v-icon icon="mdi-checkbox-blank-outline" size="small"> </v-icon>  {{ task.name }}
                  </v-card>
              </draggable>
          </v-col>
        </v-row>
      </v-col>

      <!-- Completed Tasks Column -->
      <v-col cols="3" class="d-flex flex-column completed fill-height" style="position: absolute right: 0 height: 80vh">
          <v-card-title>Tâches Terminées</v-card-title>
          <draggable group="tasks"
                     v-model="completedTasks"
                     itemKey="completedTasks"
                     :move="onMove"
                     @change = "onChange">
            <v-card class="mb-2 task" v-for="task in completedTasks" :key="task.ID">
              <v-icon icon="mdi-checkbox-marked-outline" size="small"> </v-icon> {{ task.name }}
            </v-card>
          </draggable>
      </v-col>
    </v-row>
    <v-dialog v-model="taskDialog" persistent max-width="600px">
      <v-card>
        <v-card-title>
          <span class="headline">Add a new task</span>
        </v-card-title>
        <v-card-text>
          <v-container>
            <v-row>
              <v-col cols="12">
                <v-text-field
                  label="Name"
                  v-model="taskName"
                  required
                ></v-text-field>
              </v-col>
            </v-row>
            <v-row>
              <v-col cols="12">
                <v-text-field
                  label="Description"
                  v-model="taskDescription"
                  required
                ></v-text-field>
              </v-col>
            </v-row>

            <v-row>
              <v-col cols="12">
                <v-date-picker
                  label="Due Date"
                  v-model="taskDueDate"
                  required
                ></v-date-picker>
              </v-col>
            </v-row>
          </v-container>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="blue darken-1" text @click="taskDialog = false">
            Cancel
          </v-btn>
          <v-btn color="blue darken-1" text @click="addTask">
            Add
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <v-snackbar v-model="displaySnack" timeout="3000">
      {{ message }}
    </v-snackbar>
    </v-container>
</template>

<style scoped lang="sass">
  .fill-height
    height: 100%

  .backlog
    background-color: lightgrey

  .completed
    background-color: lightgreen

  .UrgentImportant
    background-color: lightsalmon

  .UrgentNotImportant
    background-color: lightyellow

  .NotUrgentNotImportant
    background-color: lightskyblue

  .NotUrgentImportant
    background-color: lightpink

  .task
    background-color: white
    width: 95%
    height: 40px
    margin: 2px auto
    padding: 2px
    position: relative
    text-align: left
    border: 1px darkgrey solid

</style>
