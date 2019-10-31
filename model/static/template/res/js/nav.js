$(function () {

    var left = $('.left');
    var right = $('.right');
    var down = $('.down');
    var up = $('.up');
    var bg = $('.bgDiv');
    var leftNav = $('.leftNav');
    var rightNav = $('.rightNav');
    var downNav = $('.downNav');
    var upNav = $('.upNav');

    showNav(left, leftNav, "left");
    showNav(right, rightNav, "right");
    showNav(up, upNav, "up");
    showNav(down, downNav, "down");


    function showNav(btn, navDiv, direction) {
        btn.on('click', function () {
            bg.css({
                display: "block",
                transition: "opacity .5s"
            });
            if (direction == "right") {
                navDiv.css({
                    right: "0px",
                    transition: "right 1s"
                });
            } else if (direction == "left") {
                navDiv.css({
                    left: "0px",
                    transition: "left 1s"
                });
            } else if (direction == "up") {
                navDiv.css({
                    top: "0px",
                    transition: "top 1s"
                });
            } else if (direction == "down") {
                navDiv.css({
                    bottom: "0px",
                    transition: "bottom 1s"
                });
            }


        });
    }


    bg.on('click', function () {
        hideNav();
    });

    function hideNav() {
        leftNav.css({
            left: "-50%",
            transition: "left .5s"
        });
        rightNav.css({
            right: "-50%",
            transition: "right .5s"
        });
        upNav.css({
            top: "-40%",
            transition: "top .5s"
        });
        downNav.css({
            bottom: "-50%",
            webkitTransition:"bottom .5s",
            oTransition:"bottom .5s",
            mozTransition:"bottom .5s",
            transition: "bottom .5s"
        });
        bg.css({
            display: "none",
            transition: "display 1s"
        });
    }
});