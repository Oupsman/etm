let apiToken = localStorage.getItem("EtmToken")
let loggingin = false
var worker
var tokenWorker

function LoginDialog() {
    loggingin = true
    let content = "<div id='logindiv'>\n"
    content += "<form id='loginform'>\n"
    content += "<input type='text' name='username' id='username' placeholder='Username' required>"
    content += "<input type='password' name='password' id='password' placeholder='Password' required>"
    content += "</form>"
    content += "<p>No account ? Create an account <a href='/signup'>Here</a> </p>"
    content += "</div>"

    let loginDOM =  $("#logindialog")
    let loginDialog = loginDOM.dialog({
        autoOpen: 'false',
        modal: 'true',
        width: '800',
        height: '600',
        buttons: [
            {
                text: 'Exit',
                click: function () {
                    loginDialog.dialog('close');
                }},
            {
                text: 'Login',
                click : function () {
                    let username = $("#username").val();
                    let password = $("#password").val();
                    let data = {
                        name: username,
                        password: password
                    }
                    fetch('/api/v1/user/login', {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json',

                        },
                        body: JSON.stringify(data),
                    }).then(function (response) {
                        return response.json();
                    }).then(function (response) {
                        if (response.token != null) {
                            apiToken = response.token;
                            // save token in local storage
                            localStorage.setItem("EtmToken", apiToken);
                            worker.postMessage(apiToken)
                            loginDialog.dialog('close');
                            main();
                        } else {
                            popupMessage("Error logging in", "red")
                        }
                    });
                }
            }]
    });

    loginDialog.html(content)
}

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

function parseJwt (token) {
    var base64Url = token.split('.')[1];
    var base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    var jsonPayload = decodeURIComponent(atob(base64).split('').map(function(c) {
        return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
    }).join(''));

    return JSON.parse(jsonPayload);
}

function addCategory() {
    console.log('addCategory')
    const categoryName = $("#tabName").val();
    const categoryColor = $("#tabColor").val();
    const token = parseJwt(apiToken)
    const body = {
        userid: token.sub,
        name: categoryName,
        color: categoryColor
    }
    $.ajax(
        {
            url: '/api/v1/categories',
            type: 'POST',
            data: JSON.stringify(body),
            beforeSend: function(xhr){
                xhr.setRequestHeader("Content-Type","application/json");
                xhr.setRequestHeader("Accept","application/json");
            },
            dataType: 'json',
            async: false,
            success: function (msg) {
                popupMessage("Category added", "green");
                location.reload();
            },
            error: function () {
                popupMessage("Error Adding Category", "red");
            }
        }
    );
}



function updateTaskPriority (taskID, category) {
    console.log("TaskCard ID: ", taskID, "Category: ", category)
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
        headers: {
            'Authorization': 'Bearer ' + apiToken
        },
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
            popupMessage("TaskCard updated", "green");
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
        id: Number(taskID),
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
        case "completed": categoryclass = ".completed"
            break;
        case "backlog": categoryclass = ".backlog"
            break;



    }
    console.log("Body:", JSON.stringify(body))
    // API call to update the task
    $.ajax(
        {
            url: '/api/v1/task/' + taskID,
            type: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': 'Bearer ' + apiToken
            },
            data: JSON.stringify(body),
            beforeSend: function(xhr){
                xhr.setRequestHeader("Content-Type","application/json");
                xhr.setRequestHeader("Accept","application/json");
            },
            dataType: 'json',
            async: false,
            success: function (msg) {
                let taskDiv = formatTask(msg.task);
                $(taskDiv).appendTo(categoryclass);
                $(taskDiv).parent().reload();
                bindAll()
            }
        }
    );
}
Date.prototype.addDays = function(days) {
    var date = new Date(this.valueOf());
    date.setDate(date.getDate() + days);
    return date;
}

function addTask () {
    const taskName = $("#taskName");
    const taskComment = $("#taskComment");
    const taskDueDate = $("#taskDueDate");

    const name = taskName.val()
    const comment = taskComment.val()
    const dueDate = taskDueDate.val()


    let curTab = $('.ui-tabs-active');
    let curTabID = curTab.attr('aria-controls').split('-')[1];
    console.log("Current Tab: ", curTabID)

    var realDueDate

    if (dueDate === "") {
        let date = new Date
        realDueDate = date.addDays(2)
    } else {
        realDueDate = dueDate + "T00:00:00Z"
    }
    const token = parseJwt(apiToken)
    const body = {
        userid: token.sub,
        name: name,
        comment: comment,
        duedate: realDueDate,
        categoryid: curTabID,
    }

    // API call to create the task
    // $.post("/api/v1/task", JSON.stringify(body), (data, status) => {
    //    console.log(data);
    // }, "json");
    $.ajax(
        {
            url: '/api/v1/task',
            type: 'POST',
            headers: {
                'Authorization': 'Bearer ' + apiToken
            },
            data: JSON.stringify(body),
            beforeSend: function(xhr){
                xhr.setRequestHeader("Content-Type","application/json");
                xhr.setRequestHeader("Accept","application/json");
            },
            dataType: 'json',
            async: false,
            success: function (msg) {
                console.log(msg)
                let taskDiv = formatTask(msg.task);
                $('#details-task-*').hide();

                $(taskDiv).appendTo(".backlog");
                $( ".draggable" ).draggable({
                    snap: true,
                    // handle: "span.handle",
                    zIndex: 100,
                    stop: function () {
                        const task = $(this);
                        const taskID = task.attr('id').split('-')[1];
                        updateTaskPriority(taskID);
                    }
                });
                popupMessage("TaskCard added", "green");
                bindAll()
            },
            error: function () {
                popupMessage("Error Adding TaskCard", "red");
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

    const token = parseJwt(apiToken)
    const body = {
        userid: token.sub,
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
            headers: {
                'Content-Type': 'application/json',
                'Authorization': 'Bearer ' + apiToken
            },
            data: JSON.stringify(body),
            success: function(data){
                //some logic to show that the data was updated
                //then close the window$
                popupMessage("TaskCard Updated", "green");

            },
            error: function () {
                popupMessage("Error Updating TaskCard", "red");
            }
    });

    $('#task-' + taskID).html(formatTask(body));
}

function deleteTask() {
    console.log('deleteTask')
    const taskID = $("#deleteTaskID").val();

    $.ajax({
        url: '/api/v1/task/' + taskID,
        type: 'DELETE',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer ' + apiToken
        },
        success: function(data){
            //some logic to show that the data was updated
            //then close the window$
            popupMessage("TaskCard delete", "green");

            },
        error: function () {
            popupMessage("Error deleting TaskCard", "red");
        }
    });

    location.reload();
}


async function render(container, data) {
    console.log('render() called')
    let app = $(container);

    let content = "";
    content += '<div id="tabs">';
    content += ' <ul>';
    let tasksContent = '';

    for (let i = 0; i < data['categories'].length; i++) {
        const category = data.categories[i];
        console.log(category.name)
        let backlog = '<div class="backlog"><h3>Backlog</h3><button class="opener" id="add_task">Add TaskCard</button>';
        let completed = '<div id="completed" class="completed droppable"><h3>Completed</h3>';
        let notUrgentNotImportant = '<div id="00" class="NotUrgentNotImportant droppable"><h3>Delegate</h3>';
        let urgentNotImportant = '<div id="10" class="UrgentNotImportant droppable"><h3>Urgent, Not Important</h3>';
        let notUrgentImportant = '<div id="01" class="NotUrgentImportant droppable"><h3>Important but not urgent</h3>';
        let urgentImportant = '<div id="11" class="UrgentImportant droppable"><h3>Urgent and Important</h3>';

        content += '<li style="background-color:' + category.color + '; "><a href="#tabs-' + category.ID + '">' + category.name + '</a></li>';

        await fetch('/api/v1/tasks/' + category.ID, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': 'Bearer ' + apiToken
            }
        }).then(function (response) {
            return response.json();
        }).then(function (tasks) {
            tasksContent += '<div id="tabs-' + category.ID + '">';
                for (let taskNumber = 0 ; taskNumber < tasks.length ; taskNumber++ ) {
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
    let tabs = $( "#tabs" ).tabs();
    tabs.on( "tabsactivate", function( event, ui ) {
        console.log("Tab switch");
        bindAll();
    });
    // Close icon: removing the tab on click
    tabs.on( "click", "span.ui-icon-close", function() {
        let panelId = $( this ).closest( "li" ).remove().attr( "aria-controls" );
        $( "#" + panelId ).remove();
        tabs.tabs( "refresh" );
    });

    tabs.on( "keyup", function( event ) {
        if ( event.altKey && event.keyCode === $.ui.keyCode.BACKSPACE ) {
            let panelId = tabs.find( ".ui-tabs-active" ).remove().attr( "aria-controls" );
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

        await fetch('/api/v1/categories', {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': 'Bearer ' + apiToken
            }
        }).then(function (response) {
            return response.json();
        }).then(function (categories) {
            data.categories = categories;
        });

        await render('#app', data);
}

function bindAll() {
    console.log("bindAll called")
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
                addTaskDialog.dialog( "close" );
            },
            Cancel: function() {
                addTaskDialog.dialog( "close" );
            }
        },
        close: function() {
            taskForm[0].reset();
        }
    });

    // AddTab form: calls addCategory function on submit and closes the dialog
    let taskForm = addTaskDialog.find( "form" ).on( "submit", function( event ) {
        addTask();
        addTaskDialog.dialog( "close" );
        event.preventDefault();
    });

    // AddTab button: just opens the dialog
    $(".opener").button().on( "click", function() {
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
            const task = ui.draggable;
            const taskID = task.attr('id').split('-')[1];
            const category = $(this).attr('id');
            updateTaskPriority(taskID, category);
        }
    });

    $('.edittask').click(function(e, ui){
        const task = $(this).parent();
        const taskID = task.attr('id').split('-')[1];
        e.preventDefault();
        $.ajax({
            url: '/api/v1/task/' + taskID,
            type: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': "Bearer " + apiToken,
            },
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
                    editTaskDialog.dialog( "close" );
                    event.preventDefault();
                });

                editTaskDialog.dialog('open');

                $( "#editTaskDueDate" ).datepicker({
                    dateFormat: "yy-mm-dd",
                });
            }
        });
    });
    $('.deletetask').click(function(e, ui){
        const task = $(this).parent();
        const taskID = task.attr('id').split('-')[1];
        e.preventDefault();
        $.ajax({
            url: '/api/v1/task/' + taskID,
            type: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': "Bearer " + apiToken,
            },
            success: function(data){
                let deleteTaskDialog = $('#deleteTaskDialog').dialog({
                    autoOpen: 'false',
                    modal: 'true',
                    minHeight: '400px',
                    minWidth: '400px',
                    buttons: {
                        'Delete task': function () {
                            deleteTask()
                            $(this).dialog('close');
                        },

                        'Exit': function () {
                            $(this).dialog('close');
                        }
                    }
                });

                let html = "<p><form>"
                html += "<input type='hidden' name='taskID' id='deleteTaskID' value='" + data.ID + "'>"
                html += "<input type='hidden' name='taskPriority' id='deleteTaskPriority' value='" + data.priority + "'>"
                html += "<input type='hidden' name='taskUrgency' id='deleteTaskUrgency' value='" + data.urgency + "'>"
                html += "<input type='hidden' name='taskIsBackLog' id='deleteTaskIsBackLog' value='" + data.isbacklog + "'>"
                html += "<input type='hidden' name='taskIsCompleted' id='deleteTaskIsComplete' value='" + data.iscomplete + "'>"
                html += "<label for='taskName'>Name</label>"
                html += "<div id='deleteTaskName'>" + data.name +"</div>"
                html += "<label for='taskComment'>Comment</label>"
                html += "<div id='deleteTaskComment'>" + data.comment +"</div>"
                html += "<label for='taskDueDate'>Due Date</label>"
                html += "<div id='deleteTaskDueDate'>" + data.duedate +"</div>"
                html += "</form></p>"
                let deleteTaskDom = $('#deleteTaskDialog');
                deleteTaskDom.html(html);
                let deleteTaskForm = deleteTaskDom.find( "form" ).on( "submit", function( event ) {
                    deleteTask();
                    deleteTaskDialog.dialog( "close" );
                    event.preventDefault();
                });

                deleteTaskDialog.dialog('open');

            }
        });
    });

}

const check = () => {
    if (!('serviceWorker' in navigator)) {
        throw new Error('No Service Worker support!')
    }
    if (!('PushManager' in window)) {
        throw new Error('No Push API Support!')
    }
}

async function registerServiceWorker  () {
    const swRegistration = await navigator.serviceWorker.register('service-worker.js')
    navigator.serviceWorker.ready.then((registration) => {
        while (apiToken === "") {
            new Promise(r => setTimeout(r, 2000));
        }
        registration.active.postMessage(
            apiToken,
        );
    });
    return swRegistration
}

async function requestNotificationPermission() {
    const permission = await window.Notification.requestPermission()
    // value of permission can be 'granted', 'default', 'denied'
    // granted: user has accepted the request
    // default: user has dismissed the notification permission popup by clicking on x
    // denied: user has denied the request.
    if (permission !== 'granted') {
        throw new Error('Permission not granted for Notification')
    }
}

const push = async () => {
    check()
    const swRegistration = await registerServiceWorker()

    const permission = await requestNotificationPermission()
}

async function main() {
    if (apiToken === "" ) {
        LoginDialog();
    }

    await home();



    bindAll();

    if (window.Worker) {
        worker = new Worker('/static/js/webworker.js');
        worker.postMessage(apiToken)
        worker.onmessage = function (event) {
            console.log("Worker message received")
            if (event.data === "login" && !loggingin) {
                LoginDialog()
            }
        }
        tokenWorker = new Worker('/static/js/renewtoken.js');
        tokenWorker.postMessage(apiToken)
        tokenWorker.onmessage = function (event) {
            apiToken = event.data
            console.debug("New token received", apiToken)
            worker.postMessage(apiToken)
        }

    } else {
        console.log("Web worker not supported")
    }

    $( "#add_tab" ).button().on( "click", function(e) {
        e.preventDefault();
        let addCategoryDOM = $("#add_category_dialog")
        let addCategoryDialog = addCategoryDOM.dialog({
            autoOpen: 'false',
            modal: 'true',
            minHeight: '400px',
            minWidth: '400px',
            buttons: {
                'Add category': function () {
                    addCategory()
                    addCategoryDialog.dialog('close');
                },

                'Exit': function () {
                    addCategoryDialog.dialog('close');
                }
            }
        });

        let html = "<p><form>"
        html += "<label for='tabName'>Name</label>"
        html += "<input type='text' name='tabName' id='tabName' value='New Tab' class='text ui-widget-content ui-corner-all'>"
        html += "<label for='tabColor'>Color</label>\n"
        html += "<input type='color' name='tabColor' id='tabColor' value='#cccccc' class='text ui-widget-content ui-corner-all'>"
        html += "</form></p>"

        addCategoryDOM.html(html);

        let addCategoryForm = addCategoryDOM.find( "form" ).on( "submit", function( event ) {
            addCategory();
            addCategoryDialog.dialog( "close" );
            event.preventDefault();
        });

        addCategoryDialog.dialog('open');
    });

    // AddTab button: just opens the dialog


}

push();

main();
