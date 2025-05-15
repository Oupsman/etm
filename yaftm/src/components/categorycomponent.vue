<script lang="ts">
  import { defineComponent, ref, onMounted } from 'vue';
  import { useCategoryStore } from '@/stores/category.ts';
  import { useTaskStore } from '@/stores/task.ts';
  import { useAppStore } from '@/stores/app.ts';
  import { useUserStore } from '@/stores/user.ts';
  import type { Category } from '@/types/category.ts';
  import type { Task } from '@/types/task.ts';
  import { VueDraggableNext} from 'vue-draggable-next'

  export default defineComponent({
    name: 'CategoryComponent',

    components: {
      draggable: VueDraggableNext,
    },
    props: {
    category: {
      type: Object as () => Category,
      required: true,
    },
  },
  setup(props) {
    const backlog = ref<Task[]>([]);
    const urgentImportant = ref<Task[]>([]);
    const nonUrgentImportant = ref<Task[]>([]);
    const nonUrgentNonImportant = ref<Task[]>([]);
    const urgentNonImportant = ref<Task[]>([]);
    const completedTasks = ref<Task[]>([]);

    const taskDialog = ref(false);
    const taskName = ref<String>('');
    const taskDescription = ref<String>('');
    const taskDueDate = ref<Date>(new Date());

    const message = ref<String>('');
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
        '</div>';
    };

    const categoryStore = useCategoryStore();
    const taskStore = useTaskStore();

    const triggerTaskDialog = () => {
      taskDialog.value = true;
    };

    const addTask = () => {
      taskDialog.value = false;
      console.log(typeof taskDueDate.value);
      if (taskName.value && taskDescription.value && taskDueDate.value) {
        const newTask: Task = {
          name: taskName.value,
          description: taskDescription.value,
          dueDate: taskDueDate.value.toISOString(),
          categoryid: props.category.ID,
          isbackLog: true,
        };
        taskStore.addTask(newTask);
        backlog.value.push(newTask);
      }
    };
    const onMove = (evt: any) => {
      const task: Task = evt.draggedContext.element

      const destination: String = evt.to.attributes.itemkey.nodeValue
      console.log('To :', destination)
      if (destination === "backlog") {
        console.log('Backlog');
        task.isbacklog = true;
        task.iscompleted = false;
        task.urgency = false;
        task.priority = false;
      } else if (destination === "completedTasks") {

        console.log("Completed Tasks");
        task.isbacklog = false;
        task.iscompleted = true;
        task.urgency = false;
        task.priority = false;
      } else if (destination ===  "urgentImportant") {
          console.log('Urgent Important');
          task.isbacklog = false;
          task.iscompleted = false;
          task.urgency = true;
          task.priority = true;
      } else if (destination === "nonUrgentImportant") {

        console.log("Non Urgent Important");
        task.isbacklog = false;
        task.iscompleted = false;
        task.urgency = false;
        task.priority = true;

      } else if (destination === "urgentNonImportant") {

        console.log('Urgent Non Important');
        task.isbacklog = false;
        task.iscompleted = false;
        task.urgency = true;
        task.priority = false;
      } else if (destination === "nonUrgentNonImportant") {

        console.log('Non Urgent Non Important');
        task.isbacklog = false;
        task.iscompleted = false;
        task.urgency = false;
        task.priority = false;
      } else {
        console.log('Unknown destination', typeof(destination));

      }
       if (taskStore.updateTask(task.ID, task)) {
         message.value = "Task updated successfully"
       } else {
          message.value = "Task update failed"
       }
        displaySnack.value = true;
      // Add logic to handle the move event if needed
    };
    const onDragEnd = (evt: any) => {
      console.log('Drag ended ', evt, typeof evt);
    };

    onMounted(async () => {
      console.log("Category :" + JSON.stringify(props.category));
      // query tasks from the store
      const tasks = await taskStore.getTasks(props.category.ID);
      console.log("Tasks:", tasks);
      // Parse tasks and add them to the respective lists
      tasks.forEach((task: Task) => {
        if (task.isbacklog) {
          backlog.value.push(task);
        } else if (task.iscompleted) {
          completedTasks.value.push(task);
        } else if (task.urgency && task.priority) {
          urgentImportant.value.push(task);
        } else if (!task.urgency && task.priority) {
          nonUrgentImportant.value.push(task);
        } else if (task.urgency && !task.priority) {
          urgentNonImportant.value.push(task);
        } else {
          nonUrgentNonImportant.value.push(task);
        }
      });
    });

    return {
      message,
      displaySnack,
      backlog,
      urgentImportant,
      nonUrgentImportant,
      nonUrgentNonImportant,
      urgentNonImportant,
      completedTasks,
      taskDialog,
      taskName,
      taskDescription,
      taskDueDate,
      formatTask,
      triggerTaskDialog,
      addTask,
      onDragEnd,
      onMove,
    };
  },
})
</script>

<template>
  <v-container class="fill-height" style="height: 90vh;">
    <v-row class="fill-height">
      <!-- Backlog Column -->
      <v-col cols="3" class="d-flex flex-column backlog fill-height" style="position: absolute; left: 0; height:80vh;">
        <h2>Backlog</h2>
        <v-btn @click="triggerTaskDialog">Add task</v-btn>
        <draggable group="tasks"
                   :list="backlog"
                   itemKey="backlog"
                   @change="onDragEnd"
                   :move="onMove">
          <v-card class="mb-2" v-if="backlog.length > 0" v-for="task in backlog" :key="task.ID">
            {{ task.name }}
          </v-card>
        </draggable>

      </v-col>

      <!-- Eisenhower Matrix -->
      <v-col cols="6" class="mx-auto fill-height d-flex flex-column " >
        <v-row class="d-flex">
          <v-col class="mx-auto d-flex flex-column UrgentImportant" cols="6" style="position: relative; left: 0; height: 40vh;">
              <h2>Urgent et Important</h2>
              <draggable group="tasks"
                         :list="urgentImportant"
                         itemKey="urgentImportant"
                         @change="onDragEnd"
                         :move="onMove">
                <v-list v-if="urgentImportant.length > 0">
                  <v-list-item v-for="(task, index) in urgentImportant" :key="index">
                    {{ task.name }}
                  </v-list-item>
                </v-list>
              </draggable>
          </v-col>
          <v-col cols="6" class="mx-auto d-flex flex-column NotUrgentImportant" style="position: relative; left: 0; height: 40vh;">
            <h2>Non Urgent et Important</h2>
            <draggable group="tasks"
                       :list="nonUrgentImportant"
                       itemKey="nonUrgentImportant"
                       @change="onDragEnd"
                       :move="onMove">
                <v-list v-if="nonUrgentImportant.length > 0">
                  <v-list-item v-for="(task, index) in nonUrgentImportant" :key="index">
                    {{ task.name }}
                  </v-list-item>
                </v-list>
            </draggable>
          </v-col>
          <v-col cols="6" class="mx-auto d-flex flex-column UrgentNotImportant" style="position: relative; left: 0; height: 40vh;">
              <h2>Urgent et Non Important</h2>
              <draggable group="tasks"
                          :list="urgentNonImportant"
                         itemKey="urgentNonImportant"
                          @change="onDragEnd"
                         :move="onMove">
                <v-list v-if="urgentNonImportant.length > 0">
                  <v-list-item v-for="(task, index) in urgentNonImportant" :key="index">
                    {{ task.name }}
                  </v-list-item>
                </v-list>
              </draggable>
          </v-col>
          <v-col cols="6" class="mx-auto d-flex flex-column " style="position: relative; left: 0; height: 40vh;">
            <h2>Non Urgent et Non Important</h2>
              <draggable group="tasks"
                         :list="nonUrgentNonImportant"
                         itemKey="nonUrgentNonImportant"
                         @change="onDragEnd"
                         :move="onMove">
                <v-list v-if="nonUrgentNonImportant.length > 0">
                  <v-list-item v-for="(task, index) in nonUrgentNonImportant" :key="index">
                    {{ task.name }}
                  </v-list-item>
                </v-list>
              </draggable>
          </v-col>
        </v-row>
      </v-col>

      <!-- Completed Tasks Column -->
      <v-col cols="3" class="d-flex flex-column completed fill-height" style="position: absolute; right: 0; height: 80vh;">
          <v-card-title>Tâches Terminées</v-card-title>
          <draggable group="tasks"
                     :list="completedTasks"
                     itemKey="completedTasks"
                     @change="onDragEnd"
                     :move="onMove">
            <v-list v-if="completedTasks.length > 0">
              <v-list-item v-for="(task, index) in completedTasks" :key="index">

                {{ task.name }}
              </v-list-item>
            </v-list>
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

</style>
