# tracker
 
a simple server to monitor the time users spend on a website
build on TsunamiDB

# todo
1. add web interface for reading stats

# example js
```
document.addEventListener('DOMContentLoaded', function () {
    let totalTimeSpent = 0;
    let lastActiveTime = Date.now();
    let isTabActive = true;
    let userIsActive = false;
    let lastMousePosition = { x: null, y: null };
    let inactiveTime = 0;
    let id = null;
    const tracker_url = "http://127.0.0.1:5850"

    function trackTime() {
        if (isTabActive && userIsActive) {
            totalTimeSpent += Date.now() - lastActiveTime;
        }
        lastActiveTime = Date.now();
    }

    document.addEventListener('visibilitychange', () => {
        if (document.hidden) {
            isTabActive = false;
            trackTime();
        } else {
            isTabActive = true;
            lastActiveTime = Date.now();
        }
    });

    document.addEventListener('mousemove', (event) => {
        if (lastMousePosition.x !== event.clientX || lastMousePosition.y !== event.clientY) {
            userIsActive = true;
            inactiveTime = 0;
        }
        lastMousePosition = { x: event.clientX, y: event.clientY };
    });

    setInterval(() => {
        trackTime();
        inactiveTime++;
        if (inactiveTime >= 5) {
            userIsActive = false;
        }
        const totalTime = Math.round(totalTimeSpent / 1000)
        Raport(id, totalTime)
    }, 10000);

    window.addEventListener('beforeunload', () => {
        trackTime();
        console.log(`Total time spent on page: ${totalTimeSpent / 1000} seconds`);
    });

    function Register(){
        fetch(tracker_url + "/register", {
            method: "GET"
        })
        .then(response => {
            if (!response.ok) {
              throw new Error('Network response was not ok');
            }
            return response.text();
          })
        .then(response => {
            id = response;
            if(id) {
                return true
            } else {
                return false
            }
        })
    }

    function Raport(id, totalTimeSpent){
        fetch(tracker_url + "/raport", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                id: id,
                time: totalTimeSpent
            })
        })
    }

    Register();
})
```