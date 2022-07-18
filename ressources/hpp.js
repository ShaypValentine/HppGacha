function roll() {
    var rollBtn = document.getElementById("roll")
    rollBtn.disabled = true
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
    rollBtn.disabled = false;
}


document.addEventListener('click', function (event) {

    // If the clicked element doesn't have the right selector, bail
    if (!event.target.matches('.recycle')) return;
    event.preventDefault();
    recycleTarget = event.target
    if (recycleTarget !== null) {
        var quantity = recycleTarget.dataset.quantity;
        var name = recycleTarget.dataset.name;
        if (quantity !== undefined && quantity >= 4 && name !== undefined) {
            data = { quantity: quantity, name: name }
            fetch("/recycle", {
                method: "POST",
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(data)
            }).then(function (response) {
                return response.text()
            }).then(function (data) {
                obj = JSON.parse(data)
                if (obj.error_string !== "") {
                    alert(obj.error_string)
                } else {
                    recycleTarget.dataset.quantity = obj.new_quantity
                    var quantityText = document.getElementById(name + "-quantity")
                    quantityText.innerHTML = obj.new_quantity
                    Flashy('flash-messages', {
                        type: 'success',
                        title: 'Card recycled',
                        message: `One roll was added`,
                        globalClose: true,
                        expiry: 5000,
                        styles: {
                            icon: {
                                type: 'unicode',
                                val: '💬'
                            }
                        }
                    });
                    if (obj.new_quantity < 4) {
                        recycleTarget.remove()
                    }
                }
            }).catch(function (err) {
                console.warn(err)
            })
        }
    }
}, false);

  
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