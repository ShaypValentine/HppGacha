function roll() {
    var rollBtn = document.getElementById("roll")
    rollBtn.disabled = true
    if (isConnectedAndCookieExpired()) {
        window.location.replace(window.location.href);
    } else {
        let rolls = 0
        let counterAvailableRoll = ""
        isGuest = document.getElementById("connected").dataset.guest
        isConnected = document.getElementById("connected").dataset.connected;
        if (isGuest == undefined && isConnected) {
            counterAvailableRoll = document.getElementById("availableRolls")
            rolls = counterAvailableRoll.dataset.rolls
        }
        if (rolls > 0 || isGuest ) {
            fetch('/roll').then(function (response) {
                return response.text()
            }).then(function (html) {
                let rolled = document.getElementById("rolled")
                let cards =  document.getElementsByClassName("cardcount")
                if(cards.length >= 8){
                    var lastEle = cards[ cards.length-1 ];
                    if(lastEle !== undefined){
                        lastEle.remove()

                    }
                }
                rolled.innerHTML = html + rolled.innerHTML;
                if (isGuest == undefined && isConnected) {

                    if (rolls > 0) {
                        rolls = rolls - 1
                        counterAvailableRoll.innerHTML = rolls;
                        counterAvailableRoll.dataset.rolls = rolls
                    }
                }
            }).catch(function (err) {
                console.warn(err)
            })
        }
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
        var id = recycleTarget.dataset.id;
        if (quantity !== undefined && quantity >= 4 && id !== undefined) {
            data = { id: id }
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
                    var quantityText = document.getElementById(id + "-quantity")
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


document.addEventListener('click', function (event) {

    // If the clicked element doesn't have the right selector, bail
    if (!event.target.matches('.sacrifice')) return;
    event.preventDefault();
    sacrificeTarget = event.target
    if (sacrificeTarget !== null) {
        var quantity = sacrificeTarget.dataset.quantity;
        var rarity = sacrificeTarget.dataset.rarity;
        var id = sacrificeTarget.dataset.id;
        if (quantity !== undefined && quantity >= 2 && rarity !== undefined && id !== undefined) {
            data = { id: id,rarity: rarity}
            fetch("/sacrifice", {
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
                    sacrificeTarget.dataset.quantity = obj.new_quantity
                    var quantityText = document.getElementById(id + "-quantity")
                    quantityText.innerHTML = obj.new_quantity
                    Flashy('flash-messages', {
                        type: 'info',
                        title: 'Card sacrificed',
                        message: `The portal grows`,
                        globalClose: true,
                        expiry: 5000,
                        styles: {
                            icon: {
                                type: 'unicode',
                                val: '🌀'
                            }
                        }
                    });
                    if (obj.new_quantity < 2) {
                        sacrificeTarget.remove()
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

        document.getElementById("timerToRoll").innerHTML = hours + "h "
            + minutes + "m " + seconds + "s ";

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


document.addEventListener('click', function (event) {

    // If the clicked element doesn't have the right selector, bail
    if (!event.target.matches('.flip-card-back')) return;
    event.preventDefault();
    let flipInner = event.target.parentNode
    flipInner.style.transform = 'rotateY(180deg)'
})


const debouncedRoll = debounce(() => roll())

function debounce(func, wait = 200, immediate = true) {
    var timeout;
    return function() {
        var context = this, args = arguments;
        var later = function() {
            timeout = null;
            if (!immediate) func.apply(context, args);
        };
        var callNow = immediate && !timeout;
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
        if (callNow) func.apply(context, args);
    };
};