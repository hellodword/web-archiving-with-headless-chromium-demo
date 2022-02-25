(() => {
// https://stackoverflow.com/questions/50589137/scrollto-speed-duration-setting

    /*--------------------------------------------
     Functions to make scroll with speed control
    ---------------------------------------------*/

// Element or Position to move + Time in ms (milliseconds)
    function scrollTo(element, duration) {
        var e = document.documentElement;
        if (e.scrollTop === 0) {
            var t = e.scrollTop;
            ++e.scrollTop;
            e = t + 1 === e.scrollTop-- ? e : document.body;
        }
        scrollToC(e, e.scrollTop, element, duration);
    }

// Element to move, element or px from, element or px to, time in ms to animate
    function scrollToC(element, from, to, duration) {
        if (duration <= 0) return;
        if (typeof from === "object") from = from.offsetTop;
        if (typeof to === "object") to = to.offsetTop;
        // Choose one effect like easeInQuart
        scrollToX(element, from, to, 0, 1 / duration, 20, easeOutCuaic);
    }

    function scrollToX(element, xFrom, xTo, t01, speed, step, motion) {
        if (t01 < 0 || t01 > 1 || speed <= 0) {
            element.scrollTop = xTo;
            return;
        }
        element.scrollTop = xFrom - (xFrom - xTo) * motion(t01);
        t01 += speed * step;

        setTimeout(function () {
            scrollToX(element, xFrom, xTo, t01, speed, step, motion);
        }, step);
    }

    /* Effects List */
    function easeOutCuaic(t) {
        t--;
        return t * t * t + 1;
    }


    let duration = 20000;
    scrollTo(document.body.scrollHeight, duration);

    return duration+5000;
})()