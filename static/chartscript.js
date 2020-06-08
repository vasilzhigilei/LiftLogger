var chartElement;
var ctx;
var chart;
var jsondata;

if(loggedin) {
    chartElement = document.getElementById('mainchart');
    ctx = chartElement.getContext('2d');
    setupchart()
}

function setupchart() {
    $.ajax({
        url: "/getlifts" || window.location.pathname,
        type: "POST",
        success: function (data) {
            jsondata = JSON.parse(data)
            jsondatalocal = JSON.parse(data)
            chart = new Chart(ctx, {
                // The type of chart we want to create
                type: 'line',

                // The data for our dataset
                data: {
                    labels: jsondata["Date"],
                    datasets: [{
                        label: 'Weight',
                        backgroundColor: 'rgba(255,96,23,0.3)',
                        borderColor: 'rgb(255,96,23)',
                        data: jsondatalocal["Weight"]
                    },
                        {
                            label: 'Deadlift',
                            backgroundColor: 'rgba(0, 0, 0, 0)',
                            borderColor: 'rgb(255,23,58)',
                            data: jsondatalocal["Deadlift"]
                        },
                        {
                            label: 'Squat',
                            backgroundColor: 'rgba(0, 0, 0, 0)',
                            borderColor: 'rgb(255,23,243)',
                            data: jsondatalocal["Squat"]
                        },
                        {
                            label: 'Bench',
                            backgroundColor: 'rgba(0, 0, 0, 0)',
                            borderColor: 'rgb(23,62,255)',
                            data: jsondatalocal["Bench"]
                        },
                        {
                            label: 'Overhead',
                            backgroundColor: 'rgba(0, 0, 0, 0)',
                            borderColor: 'rgb(23,255,247)',
                            data: jsondatalocal["Overhead"]
                        }]
                },

                // Configuration options go here
                options: {maintainAspectRatio: false,}
            });
            console.log("lifts fetched for chart")
            if(loggedin) {
                window.addEventListener('load', function () {
                    repschange();
                });
            }
        },
        error: function (jXHR, textStatus, errorThrown) {
            alert(errorThrown);
        }
    });
}

var reppref = document.getElementById("repsdisplay")
repval = $.cookie("repspreference");
if(repval > 0) {
    reppref.value = repval;
}else {
    reppref.value = 1;
}

if(!loggedin){
    reppref.setAttribute("disabled", "disabled")
    var repcont = document.getElementById("repscontainer")
    repcont.setAttribute("data-toggle", "tooltip")
    repcont.setAttribute("data-placement", "top")
    repcont.setAttribute("title", "Log in to use!")
    $(function () {
        $('[data-toggle="tooltip"]').tooltip()
    })
}

function repschange() {
    updateReps(reppref.value)
    date = new Date(Date.now() + 365*24*60*60*1000)
    document.cookie = "repspreference=" + reppref.value + "; expires=" +
        date + "; path=/"
}
function updateReps(reps) {
    inputblocks = document.querySelectorAll(".inputblock")
    for(var i = 0; i < inputblocks.length; i++){
        if(inputblocks[i].id != "Personal") {
            fields = inputblocks[i].getElementsByClassName("form-control")
            arr = jsondata[inputblocks[i].id]
            fields[0].value = Math.round(generateXRM(arr[arr.length - 1], reps))
            fields[1].value = reps
        }
    }
    if(loggedin) {
        chart.data.datasets.forEach((dataset) => {
            if (dataset.label == "Weight") {
                return
            }
            for (i = 0; i < dataset.data.length; i++) {
                dataset.data[i] = Math.round(generateXRM(jsondata[dataset.label][i], reps))
            }
        });
        chart.update();
    }
}