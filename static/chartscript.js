var chart = document.getElementById('myChart');
var ctx = chart.getContext('2d');
chart.height = 70;
$.ajax({
    url: "/getlifts" || window.location.pathname,
    type: "POST",
    success: function (data) {
        console.log("lifts fetched for chart")
        console.log(data)
        var chart = new Chart(ctx, {
            // The type of chart we want to create
            type: 'line',

            // The data for our dataset
            data: {
                labels: ['January', 'February', 'March', 'April', 'May', 'June', 'July'],
                datasets: [{
                    label: 'Weight',
                    backgroundColor: 'rgba(0, 0, 0, 0)',
                    borderColor: 'rgb(255, 99, 132)',
                    data: [0, 10, 5, 2, 20, 30, 45]
                },
                    {
                        label: 'Deadlift',
                        backgroundColor: 'rgba(0, 0, 0, 0)',
                        borderColor: 'rgb(99, 206, 255)',
                        data: [45, 5, 10, 2, 30, 20, 10]
                    }]
            },

            // Configuration options go here
            options: {}
        });
    },
    error: function (jXHR, textStatus, errorThrown) {
        alert(errorThrown);
    }
});