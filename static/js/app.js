function addTab() {
    console.log('addTab')
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
        duedate: dueDate
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
                console.log(msg);
            }
        }
    );


}

async function render(container, data) {

    let app = $(container);

    let content = "";
    content += '<div id="tabs">';
    content += ' <ul>';
    let tasksContent = '';
    let backlog = '<div class="backlog"><button class="opener" id="add_task">Add Task</button>';
    let completed = '<div class="completed">';
    let notUrgentNotImportant = '<div class="NotUrgentNotImportant">';
    let urgentNotImportant = '<div class="UrgentNotImportant">';
    let notUrgentImportant = '<div class="NotUrgentImportant">';
    let urgentImportant = '<div class="UrgentImportant">';

    for (let i = 0; i < data['categories'].length; i++) {
        const category = data.categories[i];
        content += '<li><a href="#tabs-' + category.ID + '">' + category.Name + '</a></li>';
        await fetch('/api/v1/tasks/' + category.ID).then(function (response) {
            return response.json();
        }).then(function (tasks) {
            tasksContent += '<div id="tabs-' + category.ID + '">';
            for (let taskNumber = 0; taskNumber < tasks.length; taskNumber++) {
                const task = tasks[taskNumber];
                console.log("Task: ", task)
                const taskDiv = '<div class="task" id="task-' + task.ID + '"> <span>' + task.Name + '</span></div>';

                if (task.IsBackLog) {
                    backlog += taskDiv
                } else if (task.IsComplete) {
                    completed += taskDiv
                } else if ( ! task.Priority && !task.Urgency ) {
                    notUrgentNotImportant += taskDiv;
                } else if (! task.Priority && task.Urgency) {
                    urgentNotImportant += taskDiv
                } else if (taskPriority && ! task.Urgency) {
                    notUrgentImportant += taskDiv
                } else if (task.Priority && task.Urgency) {
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
        var tabTitle = $( "#tab_title" ),
            tabContent = $( "#tab_content" ),
            tabTemplate = "<li><a href='#{href}'>#{label}</a> <span class='ui-icon ui-icon-close' role='presentation'>Remove Tab</span></li>",
            tabCounter = 2;

        var tabs = $( "#tabs" ).tabs();

        // Modal dialog init: custom buttons and a "close" callback resetting the form inside

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
});
}

async function main() {
    await home();
    await handleTabs();

    var addCategoryDialog = $("#add_category").dialog({
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
            form[ 0 ].reset();
        }
    });

    // AddTab form: calls addTab function on submit and closes the dialog
    var form = addCategoryDialog.find( "form" ).on( "submit", function( event ) {
        addTab();
        addCategoryDialog.dialog( "close" );
        event.preventDefault();
    });

    // AddTab button: just opens the dialog
    $( "#add_tab" )
        .button()
        .on( "click", function() {
            addCategoryDialog.dialog( "open" );
        });

    var addTaskDialog = $("#add_task").dialog({
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
            form[ 0 ].reset();
        }
    });

    // AddTab form: calls addTab function on submit and closes the dialog
    var form = addTaskDialog.find( "form" ).on( "submit", function( event ) {
        addTask();
        addTaskDialog.dialog( "close" );
        event.preventDefault();
    });

    // AddTab button: just opens the dialog
    $( "#add_task" )
        .button()
        .on( "click", function() {
            addTaskDialog.dialog( "open" );
        });

}

main();
