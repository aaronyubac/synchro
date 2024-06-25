function openTab(e, tabName) {

    var tabcontent, tablinks;

    tabcontent = document.getElementsByClassName("tabcontent");
    for (i = 0; i < tabcontent.length; i++) {
        tabcontent[i].style.display = "none";
    }

    tablinks = document.getElementsByClassName("tablinks");
    for (i = 0; i < tablinks.length; i++) {
        tablinks[i].className = tablinks[i].className.replace(" active", "");
    }

    document.getElementById(tabName).style.display = "block";
    e.currentTarget.className += " active";
}

var checkbox = document.querySelector(".unavailability-all-day");

checkbox.addEventListener("change", function() {
    // let unavailabilityTimeInputs = document.querySelector(".unavailability-time");
    let unavailabilityStart = document.querySelector(".unavailability-time-start");
    let unavailabilityEnd = document.querySelector(".unavailability-time-end");

    if (this.checked) {
        // unavailabilityTimeInput.style.display = "none";
        unavailabilityStart.disabled = true;
        unavailabilityEnd.disabled = true;
        unavailabilityStart.value = "";
        unavailabilityEnd.value = "";
    } else {
        // unavailabilityTimeInput.style.display = "inline-block";
        unavailabilityStart.disabled = false;
        unavailabilityEnd.disabled = false;
    }
})