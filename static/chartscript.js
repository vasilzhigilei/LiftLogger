var chartElement = document.getElementById('mainchart');
var ctx = chartElement.getContext('2d');
var chart;
var jsondata;
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
    },
    error: function (jXHR, textStatus, errorThrown) {
        alert(errorThrown);
    }
});

document.getElementById("repsdisplay").onchange = function () {
    updateReps(this.value)
    date = new Date(Date.now() + 365*24*60*60*1000)
    document.cookie = "repspreference=" + this.value + "; expires=" +
        date + "; path=/"
}
function updateReps(reps) {
    inputblocks = document.querySelectorAll(".inputblock")
    for(var i = 0; i < inputblocks.length; i++){
        if(inputblocks[i].id != "Personal") {
            fields = inputblocks[i].getElementsByClassName("form-control")
            arr = jsondata[inputblocks[i].id]
            console.log(inputblocks[i].id)
            fields[0].value = Math.round(generateXRM(arr[arr.length - 1], reps))
            fields[1].value = reps
        }
    }
    chart.data.datasets.forEach((dataset) => {
        if(dataset.label == "Weight"){
            return
        }
        for(i = 0; i < dataset.data.length; i++){
            dataset.data[i] = Math.round(generateXRM(jsondata[dataset.label][i], reps))
        }
    });
    chart.update();
}