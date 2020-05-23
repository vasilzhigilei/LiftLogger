
var form = document.getElementById("analyzeform")
form.addEventListener('submit', submitForm);



function submitForm(event){
    // ugh this code is UGLY! Will have to DEFINITELY redo all of this!
    event.preventDefault()
    maximums = getMaximums()


}

function getData() {
    inputblocks = document.querySelectorAll(".inputblock")
    results = []
    for(var i = 0; i < inputblocks.length; i++){
        fields = inputblocks[i].getElementsByClassName("form-control")
        fields[0].validate()
        for(field in fields){
            if(isNaN(field) || field != "")
        }
        results.append([inputblocks[i].id, ])
    }


    return {dl_max, s_max, bp_max, ohp_max}
}