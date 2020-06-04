var chart = document.getElementById('myChart');
var ctx = chart.getContext('2d');
chart.height = 70;
$.ajax({
    url: "/getlifts" || window.location.pathname,
    type: "POST",
    success: function (data) {
        jsondata = JSON.parse(data)
        var chart = new Chart(ctx, {
            // The type of chart we want to create
            type: 'line',

            // The data for our dataset
            data: {
                labels: jsondata["Date"],
                datasets: [{
                        label: 'Weight',
                        backgroundColor: 'rgba(255,96,23,0.3)',
                        borderColor: 'rgb(255,96,23)',
                        data: jsondata["Weight"]
                    },
                    {
                        label: 'Deadlift',
                        backgroundColor: 'rgba(0, 0, 0, 0)',
                        borderColor: 'rgb(255,23,58)',
                        data: jsondata["Deadlift"]
                    },
                    {
                        label: 'Squat',
                        backgroundColor: 'rgba(0, 0, 0, 0)',
                        borderColor: 'rgb(255,23,243)',
                        data: jsondata["Squat"]
                    },
                    {
                        label: 'Bench',
                        backgroundColor: 'rgba(0, 0, 0, 0)',
                        borderColor: 'rgb(23,62,255)',
                        data: jsondata["Bench"]
                    },
                    {
                        label: 'Overhead',
                        backgroundColor: 'rgba(0, 0, 0, 0)',
                        borderColor: 'rgb(23,255,247)',
                        data: jsondata["Overhead"]
                    }]
            },

            // Configuration options go here
            options: {}
        });
        console.log("lifts fetched for chart")
    },
    error: function (jXHR, textStatus, errorThrown) {
        alert(errorThrown);
    }
});