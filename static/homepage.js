
var form = document.getElementById("analyzeform")
form.addEventListener('submit', submitForm);



function submitForm(event){
    // ugh this code is UGLY! Will have to DEFINITELY redo all of this!
    event.preventDefault()
    maximums = getMaximums()


}

function getData() {
    inputblocks = document.querySelectorAll(".inputblock")
    data = []
    for(var i = 0; i < inputblocks.length; i++){
        result = getFields(inputblocks[i])
        if(!isNaN(result)) {
            data.append([inputblocks[i].id, result[0], result[1]])
        }
    }


    return {dl_max, s_max, bp_max, ohp_max}
}

function getFields(inputblock){
    fields = inputblock.getElementsByClassName("form-control")
    result = []
    for(field in fields){
        if(isNaN(field.value) || field.value == ""){
            return NaN
        }else {
            result.append(parseFloat(field.value))
        }
    }
    return result
}