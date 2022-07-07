function roll(){
    fetch('/roll').then(function (response){
       return response.text()
    }).then(function(html){
    let rolled = document.getElementById("rolled")
        rolled.innerHTML = html + rolled.innerHTML;
    }).catch(function (err){
        console.warn(err)
    })
}