
var form = document.getElementById("analyzeform")
form.addEventListener('submit', submitForm);

function submitForm(event){
    // ugh this code is UGLY! Will have to DEFINITELY redo all of this!
    event.preventDefault()
    data = getData()
    if(data.length < 0)
        return // if no data, don't do anything
    for(datum in data){

    }

}

function getData() {
    // returns array of sets of lifts and weights/reps
    inputblocks = document.querySelectorAll(".inputblock")
    data = []
    for(var i = 0; i < inputblocks.length; i++){
        result = getFields(inputblocks[i])
        if(!isNaN(result)) {
            data.append([inputblocks[i].id, result[0], result[1]])
        }
    }
    return data
}

function getFields(inputblock){
    // helper function for getData
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

// Restricts input for the given textbox to the given inputFilter.
function setInputFilter(textbox, inputFilter) {
    ["input", "keydown", "keyup", "mousedown", "mouseup", "select", "contextmenu", "drop"].forEach(function(event) {
        textbox.addEventListener(event, function() {
            if (inputFilter(this.value)) {
                this.oldValue = this.value;
                this.oldSelectionStart = this.selectionStart;
                this.oldSelectionEnd = this.selectionEnd;
            } else if (this.hasOwnProperty("oldValue")) {
                this.value = this.oldValue;
                this.setSelectionRange(this.oldSelectionStart, this.oldSelectionEnd);
            } else {
                this.value = "";
            }
        });
    });
}

idsInt = ["dl-weight", "dl-reps", "s-weight", "s-reps", "bp-weight", "bp-reps", "ohp-weight", "ohp-reps", "age"]

// limit lift weights/reps from 0 to 1200
for(id in idsInt) {
    setInputFilter(document.getElementById(id), function (value) {
        return /^\d*$/.test(value) && (value === "" || parseInt(value) <= 1200);
    });
}
// limit input of person weight to two decimal places and no more than 1000 (big upper limit :))
setInputFilter(document.getElementById("currencyTextBox"), function(value) {
    return /^-?\d*[.,]?\d{0,2}$/.test(value) && (value === "" || parseFloat(value) <= 1000); });