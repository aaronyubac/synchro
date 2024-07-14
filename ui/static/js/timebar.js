const timebarGridAM = document.querySelector(".timebar.am .timebar-grid");
const timebarGridPM = document.querySelector(".timebar.pm .timebar-grid");

let timeslotsAM = "";
let timeslotsPM = "";

function renderTimeslots() {
  // AM times
  for (let i = 1; i < 49; i++) {
    if (i === 1) {
      timeslotsAM += `<div class="time-slot" data-time="15">
                            <div class="time-box"></div>
                            <div class="time-marker time-marker-start">12:00AM</div>
                        </div>`;
    } else if (i === 48) {
      timeslotsAM += `<div class="time-slot" data-time="720">
                                <div class="time-box"></div>
                                <div class="time-marker">12:00PM</div>
                            </div>`;
    } else if ((i * 15) % 60 === 0) {
      timeslotsAM += `<div class="time-slot" data-time="${i * 15}">
                                <div class="time-box"></div>
                                <div class="time-marker">${
                                  (i * 15) / 60
                                }:00AM</div>
                            </div>`;
    } else {
      timeslotsAM += `<div class="time-slot" data-time="${i * 15}">
                            <div class="time-box"></div>
                            </div>`;
    }
  }

  // PM times

  for (let i = 49; i < 97; i++) {
    if (i === 49) {
      timeslotsPM += `<div class="time-slot" data-time="735">
                                <div class="time-box"></div>
                                <div class="time-marker time-marker-start">12:00PM</div>
                            </div>`;
    } else if (i === 96) {
      timeslotsPM += `<div class="time-slot" data-time="1440">
                                <div class="time-box"></div>
                                <div class="time-marker">12:00AM</div>
                            </div>`;
    } else if ((i * 15) % 60 === 0) {
      timeslotsPM += `<div class="time-slot" data-time="${i * 15}">
                                <div class="time-box"></div>
                                <div class="time-marker">${
                                  (i * 15) / 60
                                }:00PM</div>
                            </div>`;
    } else {
      timeslotsPM += `<div class="time-slot" data-time="${i * 15}">
                                <div class="time-box"></div>
                            </div>`;
    }
  }

  timebarGridAM.innerHTML = timeslotsAM;
  timebarGridPM.innerHTML = timeslotsPM;
}

function unavailabilitiesToBar() {
  
  let current = document.querySelector(".days .day.active");
  if (current != null) {
    let unavailabilityData = JSON.parse(current.dataset.eventUnavailabilities);
    
    let timeboxes = document.querySelectorAll(".time-slot[data-time] .time-box");

    timeboxes.forEach(timebox => {
      timebox.classList = "time-box";
    });

    unavailabilityData.forEach((unavailability) => {
      if (unavailability.allDay === "true") {


        timeboxes.forEach(timebox => {
            timebox.classList.add("allDay");
        });


      } else {
        // convert start and end to minutes
        const options = {
          hour: "numeric",
          minute: "numeric",
        };

        startTime = new Date(unavailability.start);
        endTime = new Date(unavailability.end);

        startAsMinutes = (startTime.getHours() * 60) + startTime.getMinutes();
        endAsMinutes = (endTime.getHours() * 60) + endTime.getMinutes();


        for (let i = (startAsMinutes + 15); i < (endAsMinutes + 15); i += 15) {
            let currentTimebox = document.querySelector(`.time-slot[data-time="${i}"] .time-box`);

            currentTimebox.classList.add("partial");
        }

    }
    

    });
  }

  //     // keep changing until end time(converted) matches data-time
}
