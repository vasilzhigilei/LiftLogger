
var form = document.getElementById("analyzeform")
form.addEventListener('submit', submitForm);
var clicked;

function submitForm(event){
    event.preventDefault()
    fetcheddata = getData()
    if(fetcheddata["lifts"].length < 1)
        return // if no data, don't do anything
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
    for(i = 0; i < data["lifts"].length; ++i){
        datum = data["lifts"][i]
        if(datum["name"] == "Personal")
            continue
        max = Math.round(generate1RM(datum["weight"], datum["reps"]))
        rowhtml = `<div class=\"row mb-3\">` +
            `<div class=\"col-sm-4 col-md-2\"><p class=\"h5 textright\">${datum["name"]}</p></div>` +
            `<div class=\"col-sm-4 col-md-3\"><p class="h5">${max}lbs</p></div>` +
            `</div>`
        htmlstring += rowhtml
    }
    htmlstring += "</div></div>"
    return htmlstring
}

function getData() {
    // returns array of sets of lifts and weights/reps
    inputblocks = document.querySelectorAll(".inputblock")
    data = {"lifts":[]}
    for(var i = 0; i < inputblocks.length; i++){
        result = getFields(inputblocks[i])
        if(result.length > 0) {
            data["lifts"].push({"name":inputblocks[i].id, "weight":result[0], "reps":result[1]})
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
            return []
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