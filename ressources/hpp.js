function roll() {
    if (isConnectedAndCookieExpired()) {
        window.location.replace(window.location.href);
    } else {
        fetch('/roll').then(function (response) {
            return response.text()
        }).then(function (html) {
            let rolled = document.getElementById("rolled")
            rolled.innerHTML = html + rolled.innerHTML;
            let counterAvailableRoll = document.getElementById("availableRolls")
            let rolls = counterAvailableRoll.dataset.rolls
            if (rolls > 0) {
                rolls = rolls - 1
                counterAvailableRoll.innerHTML = rolls;
                counterAvailableRoll.dataset.rolls = rolls
            }
        }).catch(function (err) {
            console.warn(err)
        })
    }
}

function isConnectedAndCookieExpired() {
    var elemConnect = document.getElementById("connected")
    if (elemConnect !== null) {
        if (elemConnect.dataset.connected) {
            if (document.cookie.indexOf('session_token') == -1) {
                return true
            }
        }
        return false
    }
    return false
}

const timeToNextRun = (start) => {
    const sixHoursInMs = 2 * 3600 * 1000;
    let remainingTime = sixHoursInMs - (start.getTime() % sixHoursInMs);
    return remainingTime;
};

if (document.getElementById("timerToRoll") != null) {
    var x = setInterval(function () {
        let now = new Date();
        let countdown = timeToNextRun(now);

        // Time calculations for days, hours, minutes and seconds
        var hours = Math.floor((countdown % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
        var minutes = Math.floor((countdown % (1000 * 60 * 60)) / (1000 * 60));
        var seconds = Math.floor((countdown % (1000 * 60)) / 1000);

        // Display the result in the element with id="demo"
        document.getElementById("timerToRoll").innerHTML = hours + "h "
            + minutes + "m " + seconds + "s ";

        // If the count down is finished, write some text
        if (countdown < 0) {
            clearInterval(x);
            let counterAvailableRoll = document.getElementById("availableRolls")
            let rolls = counterAvailableRoll.dataset.rolls
            if (rolls < 4) {
                rolls = rolls + 1
                counterAvailableRoll.innerHTML = rolls;
                counterAvailableRoll.dataset.rolls = rolls
            }
        }
    }, 1000);
}