const daysContainer = document.querySelector(".days"),
  todayBtn = document.querySelector(".today-btn");
(nextBtn = document.querySelector(".next-btn")),
  (prevBtn = document.querySelector(".prev-btn")),
  (month = document.querySelector(".month"));

const months = [
  "January",
  "February",
  "March",
  "April",
  "May",
  "June",
  "July",
  "August",
  "September",
  "October",
  "Novemeber",
  "December",
];

const days = ["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"];

// get current date
const date = new Date();

// get current month
let currentMonth = date.getMonth();

// get current year
let currentYear = date.getFullYear();

let selectedDay, selectedMonth, selectedYear;

// function to render days
function renderCalendar() {




  // get prev month, current month, and next month days
  date.setDate(1);
  const firstDay = new Date(currentYear, currentMonth, 1);
  const lastDay = new Date(currentYear, currentMonth + 1, 0);
  const lastDayIndex = lastDay.getDay();
  const lastDayDate = lastDay.getDate();
  const prevLastDay = new Date(currentYear, currentMonth, 0);
  const prevLastDayDate = prevLastDay.getDate();
  const nextDays = 7 - lastDayIndex - 1;
  
  // update current year and month in header
  
    month.innerHTML = `${months[currentMonth]} ${currentYear}`;



  // update days html
  let days = "";

  // prev days html
  for (let x = firstDay.getDay(); x > 0; x--) {
    days += `<div class="day prev">${prevLastDayDate - x + 1}</div>`;
  }

  // current month days
  for (let i = 1; i <= lastDayDate; i++) {
    // check if its today then add today class
    if (
      i === new Date().getDate() &&
      currentMonth === new Date().getMonth() &&
      currentYear === new Date().getFullYear()
    ) {
      days += `<div class="day current today">${i}</div>`;
    } else {
      days += `<div class="day current">${i}</div>`;
    }
  }



  // next month days
  for (let j = 1; j <= nextDays; j++) {
    days += `<div class="day next">${j}</div>`;
  }

  hideTodayBtn();
  daysContainer.innerHTML = days;
  selectDate();


  // selected day
  for (let i = 0; i < currentDayElems.length; i++) {
    if (
      currentYear == selectedDate[0] &&
      currentMonth + 1 == selectedDate[1] &&
      i + 1 == selectedDate[2]
    ) {
      currentDayElems[i].classList.add("active");
    }

    // unavailability data
    let unavailabilitiesCurrent = [];
    let userUnavailabilitiesCurrent = [];
    
    for (let j = 0; j < unavailabilities.length; j++) {
      let start = new Date(unavailabilities[j].start);

    if (unavailabilities[j].allDay === 'true') {

      // due to way full day is stored in database
      start = new Date(start.getFullYear(), start.getMonth(), start.getDate() + 1, start.getUTCHours(), start.getUTCMinutes());

    } 
      if (i + 1 == start.getDate() && currentMonth == start.getMonth() && currentYear == start.getFullYear()) {
        if (unavailabilities[j].allDay === 'true') {
          currentDayElems[i].classList.add("allDay")
        } else {
          if (!currentDayElems[i].classList.contains("allDay")) {
             currentDayElems[i].classList.add("partial")
          }
        }
        unavailabilitiesCurrent.push(unavailabilities[j]);

        if (unavailabilities[j].userId === user.id) {
          userUnavailabilitiesCurrent.push(unavailabilities[j]);
        }
        }
      }
    
        currentDayElems[i].setAttribute("data-event-unavailabilities", JSON.stringify(unavailabilitiesCurrent));
        currentDayElems[i].setAttribute("data-user-unavailabilities", JSON.stringify(userUnavailabilitiesCurrent));
    }
  }
  



nextBtn.addEventListener("click", () => {
  currentMonth++;
  if (currentMonth > 11) {
    currentMonth = 0;
    currentYear++;
  }

  renderCalendar();
});

prevBtn.addEventListener("click", () => {
  currentMonth--;
  if (currentMonth < 0) {
    currentMonth = 11;
    currentYear--;
  }
  renderCalendar();
});

todayBtn.addEventListener("click", () => {
  currentMonth = date.getMonth();
  currentYear = date.getFullYear();

  renderCalendar();
});

function hideTodayBtn() {
  if (
    currentMonth == new Date().getMonth() &&
    currentYear == new Date().getFullYear()
  ) {
    todayBtn.style.display = "none";
  } else {
    todayBtn.style.display = "flex";
  }
}

function selectDate() {
currentDayElems = document.querySelectorAll(".current");
  currentDayElems.forEach(function (e) {
    e.addEventListener("click", function () {
      for (let i = 0; i < currentDayElems.length; i++) {
        currentDayElems[i].classList.remove("active");
      }

      e.classList.add("active");

      loadUnavailabilityTableData();
      unavailabilitiesToBar();

      for (let i = 0; i < currentDayElems.length; i++) {
        if (currentDayElems[i].classList.contains("active")) {
          selectedYear = currentYear;
          selectedMonth = currentMonth;
          selectedDay = e.innerHTML;

          let formMonth = selectedMonth + 1,
            formDay = selectedDay;

          if (formMonth < 10) formMonth = "0" + formMonth;
          if (formDay < 10) formDay = "0" + formDay;

          let formDate = selectedYear + "-" + formMonth + "-" + formDay;
          document.querySelector(".unavailability-date-input").value = formDate;

          let selectedDate =  document.querySelector(".unavailability-date-input").value.split("-");
          let displayDate = new Date(selectedDate[0], selectedDate[1] - 1, selectedDate[2]).toLocaleDateString("en-us", {year:"numeric", month:"short", day:"numeric"});
          document.querySelector(".unavailability-date-display").innerHTML = displayDate;
        }
      }
    });
  });
}
