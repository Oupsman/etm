function formatTask(task) {
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
}

function popupMessage(message, color) {
    const messagePopup = $("#message");
    messagePopup.html(message);
    messagePopup.css("background-color", color);
    messagePopup.css("color", "white");
    messagePopup.show();
    messagePopup.fadeOut(3000);
}

function addTab() {
    console.log('addTab')
}



function updateTaskPriority (taskID, category) {
    console.log("Task ID: ", taskID, "Category: ", category)
    let taskUrgency = category[0];
    let taskPriority = category[1];
    let taskName = '';
    let taskComment = '';
    let taskDueDate = '';
    let urgency = taskUrgency === '1';
    let priority = taskPriority === '1';
    let categoryclass = '';
    const iscompleted = category === 'completed';
    const isbacklog = category === 'backlog';

// get current task details
    $.ajax({
        url: '/api/v1/task/' + taskID,
        type: 'GET',
        beforeSend: function (xhr) {
            xhr.setRequestHeader("Content-Type", "application/json");
            xhr.setRequestHeader("Accept", "application/json");
        },
        dataType: 'json',
        async: false,
        success: function (msg) {
            console.log("UpdateTask : ", msg)
            const task = msg;
            taskName = task.name;
            taskComment = task.comment;
            taskDueDate = task.duedate;
            popupMessage("Task updated", "green");
        },
        error: function () {
            popupMessage("Error updating task", "red");
        }
    });

    if (iscompleted) {
        urgency = false
        priority = false
    }


    const body = {
        id: taskID.toString(),
        name: taskName,
        comment: taskComment,
        urgency: urgency,
        priority: priority,
        iscompleted: iscompleted,
        dueDate: taskDueDate,
        isbacklog: isbacklog,

    }
    switch (category) {
        case "11": categoryclass = ".UrgentImportant"
            break;
        case "10": categoryclass = ".UrgentNotImportant"
            break;
        case "01": categoryclass = ".NotUrgentImportant"
            break;
        case "00": categoryclass = ".NotUrgentNotImportant"
            break;
        completed: categoryclass = ".completed"
            break;
        backlog: categoryclass = ".backlog"
            break;



    }
    console.log("Body:", JSON.stringify(body))
    // API call to update the task
    $.ajax(
        {
            url: '/api/v1/task/' + taskID,
            type: 'POST',
            data: JSON.stringify(body),
            beforeSend: function(xhr){
                xhr.setRequestHeader("Content-Type","application/json");
                xhr.setRequestHeader("Accept","application/json");
            },
            dataType: 'json',
            async: false,
            success: function (msg) {
//                $('<p>Text</p>').appendTo('#Content');
                taskDiv = formatTask(msg);
                $(taskDiv).appendTo(categoryclass);
                $(taskDiv).parent().reload();

            }
        }
    );
}

function addTask () {
    const taskName = $("#taskName");
    const taskComment = $("#taskComment");
    const taskDueDate = $("#taskDueDate");

    const name = taskName.val()
    const comment = taskComment.val()
    const dueDate = taskDueDate.val()

    const body = {
        name: name,
        comment: comment,
        duedate: dueDate + "T00:00:00Z"
    }

    // API call to create the task
    // $.post("/api/v1/task", JSON.stringify(body), (data, status) => {
    //    console.log(data);
    // }, "json");
    $.ajax(
        {
            url: '/api/v1/task',
            type: 'POST',
            data: JSON.stringify(body),
            beforeSend: function(xhr){
                xhr.setRequestHeader("Content-Type","application/json");
                xhr.setRequestHeader("Accept","application/json");
            },
            dataType: 'json',
            async: false,
            success: function (msg) {
//                $('<p>Text</p>').appendTo('#Content');
                taskDiv = formatTask(msg);
                $('#details-task-*').hide();

                $(taskDiv).appendTo(".backlog");
                $( ".draggable" ).draggable({
                    snap: true,
                    // handle: "span.handle",
                    zIndex: 100,
                    stop: function () {
                        var task = $(this);
                        var taskID = task.attr('id').split('-')[1];
                        updateTaskPriority(taskID);
                    }
                });
                popupMessage("Task added", "green");
            },
            error: function () {
                popupMessage("Error Adding Task", "red");
            }
        }
    );


}

function editTask() {
    console.log('editTask')
    const taskID = $("#editTaskID").val();
    const taskName = $("#editTaskName").val();
    const taskComment = $("#editTaskComment").val();
    const taskDueDate = $("#editTaskDueDate").val();
    const taskPriority = $("#editTaskPriority").val();
    const taskUrgency = $("#editTaskUrgency").val();
    const taskIsBackLog = $("#editTaskIsBackLog").val();
    const taskIsComplete = $("#editTaskIsComplete").val();


    const body = {
        id: taskID,
        name: taskName,
        comment: taskComment,
        duedate: taskDueDate,
        priority: taskPriority === "true",
        urgency: taskUrgency === "true",
        isbacklog: taskIsBackLog === "true",
        iscomplete: taskIsComplete === "true",

    }

    $.ajax({
            url: '/api/v1/task/' + taskID,
            type: 'POST',
            data: JSON.stringify(body),
            success: function(data){
                //some logic to show that the data was updated
                //then close the window$
                popupMessage("Task Updated", "green");

            },
            error: function () {
                popupMessage("Error Updating Task", "red");
            }
    });

    $('#task-' + taskID).html(formatTask(body));
    // location.reload();
}

async function render(container, data) {

    let app = $(container);

    let content = "";
    content += '<div id="tabs">';
    content += ' <ul>';
    let tasksContent = '';
    let backlog = '<div class="backlog"><h3>Backlog</h3><button class="opener" id="add_task">Add Task</button>';
    let completed = '<div id="completed" class="completed droppable"><h3>Completed</h3>';
    let notUrgentNotImportant = '<div id="00" class="NotUrgentNotImportant droppable"><h3>Delegate</h3>';
    let urgentNotImportant = '<div id="10" class="UrgentNotImportant droppable"><h3>Urgent, Not Important</h3>';
    let notUrgentImportant = '<div id="01" class="NotUrgentImportant droppable"><h3>Important but not urgent</h3>';
    let urgentImportant = '<div id="11" class="UrgentImportant droppable"><h3>Urgent and Important</h3>';

    for (let i = 0; i < data['categories'].length; i++) {
        const category = data.categories[i];
        content += '<li><a href="#tabs-' + category.ID + '">' + category.Name + '</a></li>';
        await fetch('/api/v1/tasks/' + category.ID).then(function (response) {
            return response.json();
        }).then(function (tasks) {
            tasksContent += '<div id="tabs-' + category.ID + '">';
            for (let taskNumber = 0; taskNumber < tasks.length; taskNumber++) {
                const task = tasks[taskNumber];
                const taskDiv = formatTask(task);
                if (task.isbacklog) {
                    backlog += taskDiv
                } else if (task.iscomplete) {
                    completed += taskDiv
                } else if ( ! task.priority && !task.urgency ) {
                    notUrgentNotImportant += taskDiv;
                } else if (! task.priority && task.urgency) {
                    urgentNotImportant += taskDiv
                } else if (task.priority && ! task.urgency) {
                    notUrgentImportant += taskDiv
                } else if (task.priority && task.urgency) {
                    urgentImportant += taskDiv
                }
            }
            backlog += '</div>';
            completed += '</div>';
            notUrgentNotImportant += '</div>';
            urgentNotImportant += '</div>';
            notUrgentImportant += '</div>';
            urgentImportant += '</div>';
            tasksContent += backlog;
            tasksContent += completed;
            tasksContent += notUrgentImportant;
            tasksContent += urgentImportant;
            tasksContent += notUrgentNotImportant;
            tasksContent += urgentNotImportant;
            tasksContent += '</div>';
        });
    }
    content += '<li class="opener" id="add_tab"><a href="">Add</a></li>'
    content += '</ul>'
    content += tasksContent;
    content += '</div>';

    app.html(content)
    app.html();
    let addCategoryDialog = $("#add_category_dialog").dialog({
        height: 400,
        width: 550,
        modal: true,
        autoOpen: false,
        draggable: true,
        buttons: {
            Add: function() {
                addTab();
                $( this ).dialog( "close" );
            },
            Cancel: function() {
                $( this ).dialog( "close" );
            }
        },
        close: function() {
            categoryForm[0].reset();
        }
    });

    // AddTab form: calls addTab function on submit and closes the dialog
    let categoryForm = addCategoryDialog.find( "form" ).on( "submit", function( event ) {
        addTab();
        addCategoryDialog.dialog( "close" );
        event.preventDefault();
    });

    // AddTab button: just opens the dialog
    $( "#add_tab" ).button().on( "click", function() {
        addCategoryDialog.dialog( "open" );
    });
    var tabs = $( "#tabs" ).tabs();

    // Close icon: removing the tab on click
    tabs.on( "click", "span.ui-icon-close", function() {
        var panelId = $( this ).closest( "li" ).remove().attr( "aria-controls" );
        $( "#" + panelId ).remove();
        tabs.tabs( "refresh" );
    });

    tabs.on( "keyup", function( event ) {
        if ( event.altKey && event.keyCode === $.ui.keyCode.BACKSPACE ) {
            var panelId = tabs.find( ".ui-tabs-active" ).remove().attr( "aria-controls" );
            $( "#" + panelId ).remove();
            tabs.tabs( "refresh" );
        }
    });
    $("ul.tabs li").first().addClass("active");
}

async function home () {

        let data = {
            categories: []
        };

        await fetch('/api/v1/categories').then(function (response) {
            return response.json();
        }).then(function (categories) {
            data.categories = categories;
        });

        await render('#app', data);
}

async function handleTabs() {

    $( function() {


});
}

async function main() {
    await home();
    await handleTabs();


    let addTaskDom = $("#add_task_dialog")
    let addTaskDialog = addTaskDom.dialog({
        height: 400,
        width: 550,
        modal: true,
        autoOpen: false,
        draggable: true,
        buttons: {
            Add: function() {
                addTask();
                $( this ).dialog( "close" );
            },
            Cancel: function() {
                $( this ).dialog( "close" );
            }
        },
        close: function() {
            taskForm[0].reset();
        }
    });

    // AddTab form: calls addTab function on submit and closes the dialog
    var taskForm = addTaskDialog.find( "form" ).on( "submit", function( event ) {
        addTask();
        addTaskDialog.dialog( "close" );
        event.preventDefault();
    });

    // AddTab button: just opens the dialog
    $("#add_task").button().on( "click", function() {
            addTaskDialog.dialog( "open" );
        });

    // date picker for task addition
    $( "#taskDueDate" ).datepicker({
        dateFormat: "yy-mm-dd",
    });

    $( ".draggable" ).draggable({
        snap: true,
        // handle: "span.handle",
        zIndex: 100,
    });
    $( document ).tooltip();
    $('.droppable').droppable({
        drop: function (event, ui) {
            var task = ui.draggable;
            var taskID = task.attr('id').split('-')[1];
            var category = $(this).attr('id');
            updateTaskPriority(taskID, category);
        }
    });

    $('.edittask').click(function(e, ui){
        var task = $(this).parent();
        var taskID = task.attr('id').split('-')[1];
        console.log(taskID)
        e.preventDefault();
        $.ajax({
            url: '/api/v1/task/' + taskID,
            type: 'GET',
            success: function(data){
                let editTaskDialog = $('#editTaskDialog').dialog({
                    autoOpen: 'false',
                    modal: 'true',
                    minHeight: '400px',
                    minWidth: '400px',
                    buttons: {
                        'Save Changes': function () {
                            editTask()
                            $(this).dialog('close');
                        },

                        'Discard & Exit': function () {
                            $(this).dialog('close');
                        }
                    }
                });

                let html = "<p><form>"
                html += "<input type='hidden' name='taskID' id='editTaskID' value='" + data.ID + "'>"
                html += "<input type='hidden' name='taskPriority' id='editTaskPriority' value='" + data.priority + "'>"
                html += "<input type='hidden' name='taskUrgency' id='editTaskUrgency' value='" + data.urgency + "'>"
                html += "<input type='hidden' name='taskIsBackLog' id='editTaskIsBackLog' value='" + data.isbacklog + "'>"
                html += "<input type='hidden' name='taskIsCompleted' id='editTaskIsComplete' value='" + data.iscomplete + "'>"
                html += "<label for='taskName'>Name</label>"
                html += "<input type='text' name='taskName' id='editTaskName' value='" + data.name + "' class='text ui-widget-content ui-corner-all'>"
                html += "<label for='taskComment'>Comment</label>"
                html += "<input type='text' name='taskComment' id='editTaskComment' value='" + data.comment + "' class='text ui-widget-content ui-corner-all'>"
                html += "<label for='taskDueDate'>Due Date</label>"
                html += "<input type='text' name='taskDueDate' id='editTaskDueDate' value='" + data.duedate + "' class='text ui-widget-content ui-corner-all'>"
                html += "</form></p>"
                let editTaskDom = $('#editTaskDialog');
                editTaskDom.html(html);
                let editTaskForm = editTaskDom.find( "form" ).on( "submit", function( event ) {
                    editTask();
                    addTaskDialog.dialog( "close" );
                    event.preventDefault();
                });

                editTaskDialog.dialog('open');

                $( "#editTaskDueDate" ).datepicker({
                    dateFormat: "yy-mm-dd",
                });
            }
        });
    });
}

main();
