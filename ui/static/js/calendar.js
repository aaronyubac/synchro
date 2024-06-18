const daysContainer = document.querySelector(".days"),
todayBtn = document.querySelector(".today-btn")
nextBtn = document.querySelector(".next-btn"),
prevBtn = document.querySelector(".prev-btn"),
month = document.querySelector(".month");

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
    const prevLastDayDate  = prevLastDay.getDate();
    const nextDays = 7 - lastDayIndex - 1;

    // update current year and month in header
    month.innerHTML = `${months[currentMonth]} ${currentYear}`;

    // update days html
    let days = "";

    // prev days html
    for(let x = firstDay.getDay(); x > 0; x--) {
        days += `<div class="day prev">${prevLastDayDate - x + 1}</div>`;
    }

    // current month days
    for(let i = 1; i <= lastDayDate; i++) {
        // check if its today then add today class
        if (
            i === new Date().getDate() &&
            currentMonth === new Date().getMonth() &&
            currentYear === new Date().getFullYear()
        ) {
            days += `<div class="day current today" onclick="selectDate(this)">${i}</div>`;
        } else {
            days += `<div class="day current" onclick="selectDate(this)">${i}</div>`;
        }
    }


    // next month days
    for(let j = 1; j <= nextDays; j++) {
        days += `<div class="day next">${j}</div>`;
    }

    hideTodayBtn();
    daysContainer.innerHTML = days;

    // save selected day
    let dayElems = document.getElementsByClassName("current");
    for (let i = 0; i < dayElems.length; i++) {
        
        if (currentYear == selectedYear && currentMonth == selectedMonth && (i + 1) == selectedDay) {
            dayElems[i].classList.add("active");
            console.log(i);
        }
    }
}

renderCalendar();


nextBtn.addEventListener("click", () => {
    currentMonth++;
    if(currentMonth > 11) {
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
    if(
        currentMonth == new Date().getMonth() &&
        currentYear == new Date().getFullYear()
    ) {
        todayBtn.style.display = "none";
    } else {
        todayBtn.style.display = "flex";
    }
}

function selectDate(elem) {
    
let days = document.getElementsByClassName("day");

    for (let i = 0; i < days.length; i++) {
        days[i].classList.remove("active");
    }

    elem.classList.add("active");

    for (let i = 0; i < days.length; i++) {
        if (days[i].classList.contains("active")) {
            selectedYear = currentYear;
            selectedMonth = currentMonth;
            selectedDay = elem.innerHTML;

            let formMonth = selectedMonth + 1, formDay = selectedDay;

            if (formMonth < 10) formMonth = "0" + formMonth;
            if (formDay < 10) formDay = "0" + formDay;

            let formDate = currentYear + "-" + formMonth + "-" + formDay;

            document.getElementById("form-date").value = formDate;
      }
    }
}