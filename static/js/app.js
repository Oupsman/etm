
async function render(container, data) {

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

    app.html(content)
    app.html();
    $("ul.tabs li").first().addClass("active");
}

async function home () {

        let data = {
            categories: []
        };

        fetch('/api/v1/categories').then(function (response) {
            return response.json();
        }).then(function (categories) {
            data.categories = categories;
        });

        render('#app', data);
        $( "#tabs" ).tabs();


}

function handleTabs() {

    $( function() {
        var tabTitle = $( "#tab_title" ),
            tabContent = $( "#tab_content" ),
            tabTemplate = "<li><a href='#{href}'>#{label}</a> <span class='ui-icon ui-icon-close' role='presentation'>Remove Tab</span></li>",
            tabCounter = 2;

        var tabs = $( "#tabs" ).tabs();

        // Modal dialog init: custom buttons and a "close" callback resetting the form inside
        var dialog = $( "#dialog" ).dialog({
            autoOpen: false,
            modal: true,
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
        var form = dialog.find( "form" ).on( "submit", function( event ) {
            addTab();
            dialog.dialog( "close" );
            event.preventDefault();
        });

        // Actual addTab function: adds new tab using the input from the form above
        function addTab() {
            var label = tabTitle.val() || "Tab " + tabCounter,
                id = "tabs-" + tabCounter,
                li = $( tabTemplate.replace( /#\{href\}/g, "#" + id ).replace( /#\{label\}/g, label ) ),
                tabContentHtml = tabContent.val() || "Tab " + tabCounter + " content.";

            tabs.find( ".ui-tabs-nav" ).append( li );
            tabs.append( "<div id='" + id + "'><p>" + tabContentHtml + "</p></div>" );
            tabs.tabs( "refresh" );
            tabCounter++;
        }

        // AddTab button: just opens the dialog
        $( "#add_tab" )
            .button()
            .on( "click", function() {
                dialog.dialog( "open" );
            });

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

await home();
handleTabs();