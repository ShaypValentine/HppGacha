function roll() {
    if (isConnectedAndCookieExpired()) {
        window.location.replace(window.location.href);
    } else {
        fetch('/roll').then(function (response) {
            return response.text()
        }).then(function (html) {
            let rolled = document.getElementById("rolled")
            rolled.innerHTML = html + rolled.innerHTML;
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
