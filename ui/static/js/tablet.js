
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
});

function loadUnavailabilityTableData() {
    let current = document.querySelector(".days .day.active");
    if (current != null) {
        let unavailabilityData = JSON.parse(current.dataset.userUnavailabilities);
    

        
    const table = document.querySelector(".unavailability-list")
    table.innerHTML = "";

    unavailabilityData.forEach( unavailability => {

        let row = table.insertRow();
        let start = row.insertCell(0);
        let end = row.insertCell(1);
        let allDay = row.insertCell(2);
        let remove = row.insertCell(3);


        if (unavailability.allDay === 'true') {
            allDay.innerHTML = '<svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" fill="#000000" viewBox="0 0 256 256"><path d="M173.66,98.34a8,8,0,0,1,0,11.32l-56,56a8,8,0,0,1-11.32,0l-24-24a8,8,0,0,1,11.32-11.32L112,148.69l50.34-50.35A8,8,0,0,1,173.66,98.34ZM232,128A104,104,0,1,1,128,24,104.11,104.11,0,0,1,232,128Zm-16,0a88,88,0,1,0-88,88A88.1,88.1,0,0,0,216,128Z"></path></svg>';
        } else {

            const options = {
                hour: "numeric",
                minute: "numeric",
            };

            startTime = new Date(unavailability.start).toLocaleTimeString('en-us', options);
            endTime = new Date(unavailability.end).toLocaleTimeString('en-us', options);

            start.innerHTML = startTime;
            end.innerHTML = endTime;
            allDay.innerHTML = '<svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" fill="#000000" viewBox="0 0 256 256"><path d="M165.66,101.66,139.31,128l26.35,26.34a8,8,0,0,1-11.32,11.32L128,139.31l-26.34,26.35a8,8,0,0,1-11.32-11.32L116.69,128,90.34,101.66a8,8,0,0,1,11.32-11.32L128,116.69l26.34-26.35a8,8,0,0,1,11.32,11.32ZM232,128A104,104,0,1,1,128,24,104.11,104.11,0,0,1,232,128Zm-16,0a88,88,0,1,0-88,88A88.1,88.1,0,0,0,216,128Z"></path></svg>';
        
        }

            remove.innerHTML = `
            <form action="/unavailability/remove" method="post">
                <input type="hidden" name="unavailabilityId" value="${unavailability.unavailabilityId}" />
                <input type="hidden" name="eventId" value="${unavailability.eventId}" />
                <input type="hidden" name="userId" value="${unavailability.userId}" />
                <button>
                    <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" fill="#000000" viewBox="0 0 256 256"><path d="M216,48H176V40a24,24,0,0,0-24-24H104A24,24,0,0,0,80,40v8H40a8,8,0,0,0,0,16h8V208a16,16,0,0,0,16,16H192a16,16,0,0,0,16-16V64h8a8,8,0,0,0,0-16ZM96,40a8,8,0,0,1,8-8h48a8,8,0,0,1,8,8v8H96Zm96,168H64V64H192ZM112,104v64a8,8,0,0,1-16,0V104a8,8,0,0,1,16,0Zm48,0v64a8,8,0,0,1-16,0V104a8,8,0,0,1,16,0Z"></path></svg>
                </button>
            </form>`;
            

    });
    }
}
