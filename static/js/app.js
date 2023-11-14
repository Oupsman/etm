
async function render(container, data) {
    console.log("In render")
    console.log(JSON.stringify(data))

    let app = $(container);

    let content = "";
    content += '<div id="tabs">';
    content += ' <ul>';
    let tasksContent = '';
    let backlog = '<div class="backlog">';
    let completed = '<div class="completed">';
    let notUrgentNotImportant = '<div class="NotUrgentNotImportant">';
    let urgentNotImportant = '<div class="UrgentNotImportant">';
    let notUrgentImportant = '<div class="NotUrgentImportant">';
    let urgentImportant = '<div class="UrgentImportant">';

    for (let i = 0; i < data['categories'].length; i++) {
        const category = data.categories[i];
        console.log("i: " + i)
        content += '<li><a href="#tabs-' + category.ID + '">' + category.Name + '</a></li>';
        await fetch('/api/v1/tasks/' + category.ID).then(function (response) {
            return response.json();
        }).then(function (tasks) {
            tasksContent += '<div id="tabs-' + category.ID + '">';
            for (let taskNumber = 0; taskNumber < tasks.length; taskNumber++) {
                const task = tasks[taskNumber];
                const taskDiv = '<div id="task-' + task.ID + '">' + task.Name + '</div>';
                if (task.isBackLog) {
                    backlog += taskDiv
                } else if (task.isComplete) {
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

    content += '</ul>'
    content += tasksContent;
    content += '</div>';
    console.log("content")
    console.log(content)
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

    console.log("Before render")
    console.log(data)
    render('#app', data);
}

home();

