const timebarGridAM = document.querySelectorAll(".timebar .am .timebar-grid");
const timebarGridPM = document.querySelectorAll(".timebar .pm .timebar-grid");

let timeslotsAM = "";
let timeslotsPM = "";

function renderTimeslots() {
  // AM times
  for (let i = 1; i < 49; i++) {
    if (i === 1) {
      timeslotsAM += `<div class="time-slot" data-time="15">
                            <div class="time-marker time-marker-start">12:00AM</div>
                        </div>`;
    } else if (i === 48) {
      timeslotsAM += `<div class="time-slot" data-time="720">
                                <div class="time-marker">12:00PM</div>
                            </div>`;
    } else if ((i * 15) % 60 === 0) {
      timeslotsAM += `<div class="time-slot" data-time="${i * 15}">
                                <div class="time-marker">${
                                  (i * 15) / 60
                                }:00AM</div>
                            </div>`;
    } else {
      timeslotsAM += `<div class="time-slot" data-time="${i * 15}">
                            </div>`;
    }
  }

  // PM times

  for (let i = 49; i < 97; i++) {
    if (i === 49) {
      timeslotsPM += `<div class="time-slot" data-time="735">
                                <div class="time-marker time-marker-start">12:00PM</div>
                            </div>`;
    } else if (i === 96) {
      timeslotsPM += `<div class="time-slot" data-time="1440">
                                <div class="time-marker">12:00AM</div>
                            </div>`;
    } else if ((i * 15) % 60 === 0) {
      timeslotsPM += `<div class="time-slot" data-time="${i * 15}">
                                <div class="time-marker">${
                                  ((i * 15) / 60) % 12
                                }:00PM</div>
                            </div>`;
    } else {
      timeslotsPM += `<div class="time-slot" data-time="${i * 15}">
                            </div>`;
    }
  }

  timebarGridAM[0].innerHTML = timeslotsAM;
  timebarGridPM[0].innerHTML = timeslotsPM;
  timebarGridAM[1].innerHTML = timeslotsAM;
  timebarGridPM[1].innerHTML = timeslotsPM;
}

function unavailabilitiesToBar() {
  let current = document.querySelector(".days .day.active");
  if (current != null) {
    let eventUnavailabilities = JSON.parse(current.dataset.eventUnavailabilities);
    let userUnavailabilities = JSON.parse(current.dataset.userUnavailabilities);

    const unavailabilityData = [userUnavailabilities, eventUnavailabilities];

    let timeboxes = [];
    let userTab = Boolean;
    
    for (let i = 0; i < 2; i++) {
      
      if (i === 0) {
        userTab = true;
        timeboxes = document.querySelectorAll(".user .time-slot[data-time]");
      } else if (i === 1) {
        userTab = false;
        timeboxes = document.querySelectorAll(".event .time-slot[data-time]");
      }

      timeboxes.forEach((timebox) => {
        timebox.classList = "time-slot";
      });

      unavailabilityData[i].forEach((unavailability) => {
        if (unavailability.allDay === "true") {

          timeboxes.forEach((timebox) => {
            timebox.classList.add("allDay");
          });

        } else {

          startTime = new Date(unavailability.start);
          endTime = new Date(unavailability.end);

          startAsMinutes = startTime.getHours() * 60 + startTime.getMinutes();
          endAsMinutes = endTime.getHours() * 60 + endTime.getMinutes();

          for (let i = startAsMinutes + 15; i < endAsMinutes + 15; i += 15) {
            let currentTimebox;
            if (userTab) {
              currentTimebox = document.querySelector(`.user .time-slot[data-time="${i}"]`);
            } else {
              currentTimebox = document.querySelector(`.event .time-slot[data-time="${i}"]`);
            }

            if (currentTimebox != null) {
              currentTimebox.classList.add("partial");
            }
          }
        }
      });
    }
  }
}
