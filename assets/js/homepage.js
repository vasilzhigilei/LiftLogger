
var form = document.getElementById("analyzeform")
form.addEventListener('submit', submitForm);

function submitForm(event){
    // ugh this code is UGLY! Will have to DEFINITELY redo all of this!
    event.preventDefault()
    maximums = getMaximums()

    var insides = ''
    if(!isNaN(maximums["dl_max"])){
        insides += '<div class=\"row mb-3\"><div class="col-sm-4 col-md-2"><p>Deadlift: ' + maximums["dl_max"] + '</p></div></div>'
    }
    if(!isNaN(maximums["s_max"])){
        insides += '<div class=\"row mb-3\"><div class="col-sm-4 col-md-2"><p>Squat: ' + maximums["s_max"] + '</p></div></div>'
    }
    if(!isNaN(maximums["bp_max"])){
        insides += '<div class=\"row mb-3\"><div class="col-sm-4 col-md-2"><p>Bench Press: ' + maximums["bp_max"] + '</p></div></div>'
    }
    if(!isNaN(maximums["ohp_max"])){
        insides += '<div class=\"row mb-3\"><div class="col-sm-4 col-md-2"><p>Overhead Press: ' + maximums["ohp_max"] + '</p></div></div>'
    }

    $("#maincontent").append('<br><div class="row"><div class="col-md-12 bg-light padding rounded">' +
        insides +
        '</div></div>')
}

function getMaximums() {
    var dl_weight = parseFloat(document.getElementById("dl-weight").value)
    var dl_reps = parseFloat(document.getElementById("dl-reps").value)
    var s_weight = parseFloat(document.getElementById("s-weight").value)
    var s_reps = parseFloat(document.getElementById("s-reps").value)
    var bp_weight = parseFloat(document.getElementById("bp-weight").value)
    var bp_reps = parseFloat(document.getElementById("bp-reps").value)
    var ohp_weight = parseFloat(document.getElementById("ohp-weight").value)
    var ohp_reps = parseFloat(document.getElementById("ohp-reps").value)

    var dl_max = Math.round(generate1RM(dl_weight, dl_reps))
    var s_max = Math.round(generate1RM(s_weight, s_reps))
    var bp_max = Math.round(generate1RM(bp_weight, bp_reps))
    var ohp_max = Math.round(generate1RM(ohp_weight, ohp_reps))

    return {dl_max, s_max, bp_max, ohp_max}
}