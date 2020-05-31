
var form = document.getElementById("analyzeform")
form.addEventListener('submit', submitForm);
var clicked;

function submitForm(event){
    event.preventDefault()
    fetcheddata = getData()
    document.getElementById("results").innerHTML = buildHTML(fetcheddata)
    if(clicked == "log") {
        $.ajax({
            url: $(form).attr('action') || window.location.pathname,
            type: "POST",
            data: fetcheddata,
            success: function (data) {
                console.log("lifts logged")
            },
            error: function (jXHR, textStatus, errorThrown) {
                alert(errorThrown);
            }
        });
    }
}

function buildHTML(data){
    htmlstring = "<br><div class=\"row\"><div class=\"col-md-12 bg-light padding rounded\">" +
            "<p class=\"h3\">Estimated 1 Rep Maxes</p><hr/>"
    for (var key in data){
        if(key == "Weight" || key == "Age")
            continue
        rowhtml = `<div class=\"row mb-3\">` +
            `<div class=\"col-sm-4 col-md-2\"><p class=\"h5 textright\">${key}</p></div>` +
            `<div class=\"col-sm-4 col-md-3\"><p class="h5">${data[key]}lbs</p></div>` +
            `</div>`
        htmlstring += rowhtml
    }
    htmlstring += "</div></div>"
    return htmlstring
}

function getData() {
    // returns array of sets of lifts and 1Rep Max
    inputblocks = document.querySelectorAll(".inputblock")
    data = {}
    for(var i = 0; i < inputblocks.length; i++){
        result = getFields(inputblocks[i])
        if(inputblocks[i].id === "Personal") {
            console.log(result[1])
            data["Weight"] = result[0]
            data["Age"] = result[1]
        }else {
            if (result.length > 1) {
                data[inputblocks[i].id] = Math.round(generate1RM(result[0], result[1]))
            } else {
                data[inputblocks[i].id] = 0
            }
        }
    }
    return data
}

function getFields(inputblock){
    // helper function for getData
    fields = inputblock.getElementsByClassName("form-control")
    result = []
    for(i = 0; i < fields.length; ++i){
        field = fields[i]
        if(isNaN(field.value) || field.value == ""){
            result.push(0.0)
        } else {
            result.push(parseFloat(field.value))
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

idsInt = ["dl-weight", "dl-reps", "s-weight", "s-reps", "bp-weight", "bp-reps", "ohp-weight", "ohp-reps", "u-age"]

// limit lift weights/reps from 0 to 1200
for(i = 0; i < idsInt.length; ++i) {
    setInputFilter(document.getElementById(idsInt[i]), function (value) {
        return /^\d*$/.test(value) && (value === "" || parseInt(value) <= 1200);
    });
}
// limit input of person weight to two decimal places and no more than 1000 (big upper limit :))
setInputFilter(document.getElementById("u-weight"), function(value) {
    return /^-?\d*[.,]?\d{0,2}$/.test(value) && (value === "" || parseFloat(value) <= 1000); });